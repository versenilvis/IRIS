package root

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/versenilvis/iris/internal/config"
)

func TestFetchReleasesFiltersPrereleasesOnStableChannel(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	cfg := config.DefaultConfig()
	cfg.Updater.Channel = "stable"
	config.Init(cfg)

	releases := []Release{
		{TagName: "v0.6.0-nightly.abc", Prerelease: true, Body: "nightly"},
		{TagName: "v0.5.2", Prerelease: false, Body: "stable"},
		{TagName: "v0.5.1", Prerelease: false, Body: "stable"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	got, err := FetchReleases(0)
	if err != nil {
		t.Fatalf("FetchReleases: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 stable releases, got %d: %+v", len(got), got)
	}
	for _, r := range got {
		if r.Prerelease {
			t.Errorf("prerelease %q leaked into stable channel results", r.TagName)
		}
	}
}

func TestFetchReleasesLimitsCount(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	config.Init(config.DefaultConfig())

	releases := []Release{
		{TagName: "v0.5.2"}, {TagName: "v0.5.1"}, {TagName: "v0.5.0"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	got, err := FetchReleases(2)
	if err != nil {
		t.Fatalf("FetchReleases: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected limit of 2, got %d", len(got))
	}
}

func seedChangelogCache(t *testing.T, fetchedAt time.Time, releases []Release) {
	t.Helper()
	path, err := changelogCachePath()
	if err != nil {
		t.Fatalf("changelogCachePath: %v", err)
	}
	data, err := json.Marshal(changelogCache{FetchedAt: fetchedAt, Channel: config.Get().Updater.Channel, Releases: releases})
	if err != nil {
		t.Fatalf("marshal cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir cache dir: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write cache: %v", err)
	}
}

func TestFetchReleasesCachedServesWithinTTL(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	seedChangelogCache(t, time.Now(), []Release{{TagName: "v1.0.0", Body: "cached"}})

	// point at an address nothing listens on: a network call here is a test failure
	t.Setenv("IRIS_CHANGELOG_URL", "http://127.0.0.1:1/unreachable")

	got, rateLimited, err := FetchReleasesCached(1, false)
	if err != nil {
		t.Fatalf("FetchReleasesCached: %v", err)
	}
	if rateLimited {
		t.Error("expected rateLimited=false for a fresh cache hit")
	}
	if len(got) != 1 || got[0].TagName != "v1.0.0" {
		t.Fatalf("expected cached release, got %+v", got)
	}
}

func TestFetchReleasesCachedChannelChangeBypassesStaleCache(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	cfg := config.DefaultConfig()
	cfg.Updater.Channel = "nightly"
	config.Init(cfg)

	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	seedChangelogCache(t, time.Now(), []Release{{TagName: "v1.0.0-nightly.abc", Body: "nightly cache", Prerelease: true}})

	cfg.Updater.Channel = "stable"
	config.Init(cfg)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Release{{TagName: "v1.0.0", Body: "stable release"}})
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	got, _, err := FetchReleasesCached(1, false)
	if err != nil {
		t.Fatalf("FetchReleasesCached: %v", err)
	}
	if len(got) != 1 || got[0].TagName != "v1.0.0" {
		t.Fatalf("expected the channel switch to bypass the nightly-channel cache and fetch fresh, got %+v", got)
	}
}

func TestFetchReleasesCachedRefreshBypassesTTL(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	config.Init(config.DefaultConfig())

	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	seedChangelogCache(t, time.Now(), []Release{{TagName: "v1.0.0", Body: "stale-but-fresh"}})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Release{{TagName: "v2.0.0", Body: "fresh"}})
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	got, _, err := FetchReleasesCached(1, true)
	if err != nil {
		t.Fatalf("FetchReleasesCached: %v", err)
	}
	if len(got) != 1 || got[0].TagName != "v2.0.0" {
		t.Fatalf("expected --refresh to bypass the cache, got %+v", got)
	}
}

func TestFetchReleasesCachedRateLimitFallsBackToStaleCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	seedChangelogCache(t, time.Now().Add(-2*changelogCacheTTL), []Release{{TagName: "v1.0.0", Body: "stale"}})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	got, rateLimited, err := FetchReleasesCached(1, false)
	if err != nil {
		t.Fatalf("FetchReleasesCached: %v", err)
	}
	if !rateLimited {
		t.Error("expected rateLimited=true when the API 403s and a stale cache exists")
	}
	if len(got) != 1 || got[0].TagName != "v1.0.0" {
		t.Fatalf("expected stale cached release, got %+v", got)
	}
}

func TestFetchReleasesCachedRateLimitNoCacheReturnsError(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	_, _, err := FetchReleasesCached(1, false)
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestPrintReleaseSuppressesUpdateCTAWhenCurrent(t *testing.T) {
	originalVersion := Version
	Version = "v1.0.0"
	t.Cleanup(func() { Version = originalVersion })

	var buf bytes.Buffer
	printRelease(&buf, Release{TagName: "v1.0.0", Body: "## Changelog\n### Bug fixes\n* abc1234  fix something\n"}, true)

	if strings.Contains(buf.String(), "run `iris update`") {
		t.Error("expected no update CTA when already on the latest version")
	}
	if !strings.Contains(buf.String(), "(current)") {
		t.Error("expected the version label to be marked (current)")
	}
}

func TestPrintReleaseShowsUpdateCTAWhenNewer(t *testing.T) {
	originalVersion := Version
	Version = "v0.9.0"
	t.Cleanup(func() { Version = originalVersion })

	var buf bytes.Buffer
	printRelease(&buf, Release{TagName: "v1.0.0", Body: "## Changelog\n### Bug fixes\n* abc1234  fix something\n"}, true)

	if !strings.Contains(buf.String(), "run `iris update`") {
		t.Error("expected an update CTA when a newer release is shown")
	}
}

func TestPrintChangelogBodyStopsAtUpdateFooter(t *testing.T) {
	var buf bytes.Buffer
	body := "## Changelog\n### Bug fixes\n* abc1234  fix something\n## Update\n```bash\niris update\n```\n"
	printChangelogBody(&buf, body)
	plain := ansi.Strip(buf.String())

	if strings.Contains(plain, "Update") {
		t.Error("expected the GoReleaser footer to be excluded from rendered body")
	}
	if !strings.Contains(plain, "fix something") {
		t.Errorf("expected the changelog entry to be rendered, got %q", plain)
	}
}

func TestPrintChangelogBodyRendersArbitraryMarkdown(t *testing.T) {
	var buf bytes.Buffer
	body := "# 🌙 **Nightly build** - `v0.6.0-nightly.abc123`\nNightly release from commit: abc123\n"
	printChangelogBody(&buf, body)
	plain := ansi.Strip(buf.String())

	if !strings.Contains(plain, "Nightly build") {
		t.Errorf("expected the body rendered via glamour, got %q", plain)
	}
}

func TestPrintChangelogBodyEmptyBodyPrintsNothing(t *testing.T) {
	var buf bytes.Buffer
	printChangelogBody(&buf, "")

	if buf.Len() != 0 {
		t.Errorf("expected no output for an empty body, got %q", buf.String())
	}
}

func TestPrintChangelogBodyStripsRedundantHeadingAndBlankRuns(t *testing.T) {
	var buf bytes.Buffer
	body := "## Changelog\n### Bug fixes\n* abc1234  fix something\n"
	printChangelogBody(&buf, body)
	plain := ansi.Strip(buf.String())

	if strings.Contains(plain, "Changelog") {
		t.Errorf("expected the redundant ## Changelog heading to be stripped, got %q", plain)
	}
	if strings.Contains(buf.String(), "\n\n\n") {
		t.Errorf("expected consecutive blank lines to be squeezed, got %q", buf.String())
	}
}

func TestStripRedundantHeadingRemovesExactGoReleaserLabelOnly(t *testing.T) {
	got := stripRedundantHeading("## Changelog\n### Bug fixes\n* item\n")
	want := "### Bug fixes\n* item\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripRedundantHeadingPreservesHandWrittenHeadings(t *testing.T) {
	body := "## Theme is here\n<img src=\"...\" />\n"
	got := stripRedundantHeading(body)
	if got != body {
		t.Errorf("expected a hand-written ## heading to be preserved, got %q", got)
	}
}

func TestPrintChangelogBodyPreservesHandWrittenAnnouncement(t *testing.T) {
	var buf bytes.Buffer
	body := "## Theme is here\r\n<img width=\"1818\" src=\"...\" />\r\n\r\n## Update\r\n```bash\r\niris update\r\n```\r\n"
	printChangelogBody(&buf, body)
	plain := ansi.Strip(buf.String())

	if !strings.Contains(plain, "Theme is here") {
		t.Errorf("expected the hand-written release heading to be rendered, got %q", plain)
	}
}

func TestSqueezeBlankLinesCollapsesRuns(t *testing.T) {
	got := squeezeBlankLines("a\n\n\n\nb\n")
	want := "a\n\nb\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChangelogCmdShowsSpecificVersion(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	config.Init(config.DefaultConfig())

	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	releases := []Release{
		{TagName: "v0.5.2", Body: "## Changelog\n### Bug fixes\n* xyz9999  unrelated fix\n"},
		{TagName: "v0.5.1", Body: "## Changelog\n### Bug fixes\n* abc1234  fix mode\n"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	var out, errOut bytes.Buffer
	ChangelogCmd.SetOut(&out)
	ChangelogCmd.SetErr(&errOut)
	t.Cleanup(func() {
		ChangelogCmd.SetOut(nil)
		ChangelogCmd.SetErr(nil)
	})

	ChangelogCmd.Run(ChangelogCmd, []string{"v0.5.1"})
	got := ansi.Strip(out.String())

	if !strings.Contains(got, "v0.5.1") {
		t.Fatalf("expected the requested version in output, got %q", got)
	}
	if strings.Contains(got, "v0.5.2") || strings.Contains(got, "unrelated fix") {
		t.Errorf("expected only the requested version, got %q", got)
	}
	if !strings.Contains(got, "fix mode") {
		t.Errorf("expected the matched release's body, got %q", got)
	}
}

func TestChangelogCmdVersionNotFound(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() { config.Init(originalConfig) })
	config.Init(config.DefaultConfig())

	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	releases := []Release{{TagName: "v0.5.2", Body: "latest"}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()
	t.Setenv("IRIS_CHANGELOG_URL", srv.URL)

	var out, errOut bytes.Buffer
	ChangelogCmd.SetOut(&out)
	ChangelogCmd.SetErr(&errOut)
	t.Cleanup(func() {
		ChangelogCmd.SetOut(nil)
		ChangelogCmd.SetErr(nil)
	})

	ChangelogCmd.Run(ChangelogCmd, []string{"v9.9.9"})

	if !strings.Contains(errOut.String(), "not found") {
		t.Errorf("expected a not-found error, got stdout=%q stderr=%q", out.String(), errOut.String())
	}
}
