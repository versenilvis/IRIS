package root

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
)

// Release mirrors the fields we need from the GitHub releases API.
type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	Prerelease  bool      `json:"prerelease"`
}

type changelogCache struct {
	FetchedAt time.Time `json:"fetched_at"`
	Releases  []Release `json:"releases"`
}

const changelogCacheTTL = time.Hour

func changelogCachePath() (string, error) {
	dir, err := config.CachePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "changelog-cache.json"), nil
}

func loadChangelogCache() (*changelogCache, error) {
	path, err := changelogCachePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cache changelogCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func saveChangelogCache(releases []Release) {
	path, err := changelogCachePath()
	if err != nil {
		return
	}
	data, err := json.Marshal(changelogCache{FetchedAt: time.Now(), Releases: releases})
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, data, 0o644)
}

// FetchReleases hits the GitHub releases API and returns up to limit
// releases, filtering out prereleases (nightlies) unless the nightly
// channel is active. A limit of 0 returns everything fetched.
func FetchReleases(limit int) ([]Release, error) {
	ctx, cancel := newGitHubRequestContext()
	defer cancel()

	endpoint := os.Getenv("IRIS_CHANGELOG_URL")
	if endpoint == "" {
		endpoint = "https://api.github.com/repos/versenilvis/iris/releases?per_page=100"
	}

	body, err := fetchGitHubBody(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	if config.Get().Updater.Channel != "nightly" {
		filtered := releases[:0]
		for _, r := range releases {
			if !r.Prerelease {
				filtered = append(filtered, r)
			}
		}
		releases = filtered
	}

	return truncateReleases(releases, limit), nil
}

func truncateReleases(releases []Release, limit int) []Release {
	if limit > 0 && len(releases) > limit {
		return releases[:limit]
	}
	return releases
}

// FetchReleasesCached serves cached releases when they're within TTL, and
// falls back to a stale cache when the API is rate limited rather than
// failing outright. rateLimited reports whether the result came from that
// fallback, so the caller can say so.
func FetchReleasesCached(limit int, refresh bool) (releases []Release, rateLimited bool, err error) {
	if !refresh {
		if cache, cacheErr := loadChangelogCache(); cacheErr == nil && time.Since(cache.FetchedAt) < changelogCacheTTL {
			return truncateReleases(cache.Releases, limit), false, nil
		}
	}

	fresh, fetchErr := FetchReleases(100)
	if fetchErr != nil {
		if errors.Is(fetchErr, ErrRateLimited) {
			if cache, cacheErr := loadChangelogCache(); cacheErr == nil && len(cache.Releases) > 0 {
				return truncateReleases(cache.Releases, limit), true, nil
			}
		}
		return nil, false, fetchErr
	}

	saveChangelogCache(fresh)
	return truncateReleases(fresh, limit), false, nil
}

func releaseURL(tag string) string {
	return "https://github.com/versenilvis/iris/releases/tag/" + tag
}

// printChangelogBody re-renders a release body as styled terminal output.
// Stable releases go through GoReleaser (see the `changelog:` block in
// .goreleaser.yaml): "### Group" headings, "* {sha}  {message}" entries,
// then a "## Update" footer we stop at since this command prints its own
// call-to-action. Nightly releases don't go through GoReleaser at all
// (built directly by .github/workflows/nightly.yml) and have a fixed
// plain-text body that matches none of those patterns - so if nothing
// matched by the end, the raw body is printed as a fallback rather than
// silently showing nothing.
func printChangelogBody(out io.Writer, body string) {
	firstGroup := true
	matched := false
	for line := range strings.SplitSeq(body, "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case strings.HasPrefix(line, "## Update"):
			return
		case strings.HasPrefix(line, "## "):
			continue
		case strings.HasPrefix(line, "### "):
			matched = true
			if !firstGroup {
				fmt.Fprintln(out)
			}
			firstGroup = false
			fmt.Fprintf(out, "\033[1m%s\033[0m\n", strings.TrimPrefix(line, "### "))
		case strings.HasPrefix(line, "* "):
			matched = true
			entry := strings.TrimPrefix(line, "* ")
			sha, msg, ok := strings.Cut(entry, "  ")
			if ok {
				fmt.Fprintf(out, "  \033[2m%s\033[0m  %s\n", sha, msg)
			} else {
				fmt.Fprintf(out, "  %s\n", entry)
			}
		}
	}

	if !matched {
		if trimmed := strings.TrimSpace(body); trimmed != "" {
			fmt.Fprintln(out, trimmed)
		}
	}
}

func printRelease(out io.Writer, release Release, showUpdateCTA bool) {
	label := release.TagName
	if strings.TrimPrefix(release.TagName, "v") == strings.TrimPrefix(Version, "v") {
		label += " (current)"
	}

	date := release.PublishedAt.Format("2006-01-02")
	fmt.Fprintf(out, "\033[1;36mIRIS %s - %s\033[0m\n", label, date)
	fmt.Fprintf(out, "\033[2m%s\033[0m\n\n", releaseURL(release.TagName))

	printChangelogBody(out, release.Body)

	if showUpdateCTA && IsNewer(Version, release.TagName) {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "run `iris update` to install")
	}
}

var (
	changelogCount   int
	changelogRefresh bool
)

// ChangelogCmd shows grouped release notes fetched from GitHub releases.
var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "show what changed in recent iris releases",
	Run: func(cmd *cobra.Command, args []string) {
		releases, rateLimited, err := FetchReleasesCached(changelogCount, changelogRefresh)
		if err != nil {
			if errors.Is(err, ErrRateLimited) {
				fmt.Fprintln(cmd.ErrOrStderr(), "\033[31m[IRIS] rate limited by GitHub API, try again later\033[0m")
				return
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "\033[31m[IRIS] could not fetch changelog: %v\033[0m\n", err)
			return
		}
		if len(releases) == 0 {
			cmd.Println("no releases found")
			return
		}

		out := cmd.OutOrStdout()
		for i, release := range releases {
			if i > 0 {
				fmt.Fprintln(out)
			}
			printRelease(out, release, i == 0)
		}
		if rateLimited {
			fmt.Fprintln(out)
			fmt.Fprintln(out, "\033[33m[IRIS] rate limited by GitHub API, showing cached data\033[0m")
		}
	},
}

func init() {
	ChangelogCmd.Flags().IntVarP(&changelogCount, "count", "n", 1, "number of releases to show")
	ChangelogCmd.Flags().BoolVar(&changelogRefresh, "refresh", false, "bypass the cache and fetch fresh data")
	rootCmd.AddCommand(ChangelogCmd)
}
