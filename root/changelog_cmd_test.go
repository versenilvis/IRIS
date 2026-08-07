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
	data, err := json.Marshal(changelogCache{FetchedAt: fetchedAt, Releases: releases})
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

	if strings.Contains(buf.String(), "Update") {
		t.Error("expected the GoReleaser footer to be excluded from rendered body")
	}
	if !strings.Contains(buf.String(), "fix something") {
		t.Error("expected the changelog entry to be rendered")
	}
}

func TestPrintChangelogBodyFallsBackToRawTextForUnstructuredBody(t *testing.T) {
	var buf bytes.Buffer
	// nightly releases don't go through GoReleaser's changelog: block (see
	// .github/workflows/nightly.yml) - their body has none of the "### "/
	// "* " markers, so it should be shown as-is rather than silently dropped
	body := "# 🌙 **Nightly build** - `v0.6.0-nightly.abc123`\nNightly release from commit: abc123\n\n# Update: \n`iris update`\n"
	printChangelogBody(&buf, body)

	if !strings.Contains(buf.String(), "Nightly build") {
		t.Errorf("expected the raw nightly body as a fallback, got %q", buf.String())
	}
}

func TestPrintChangelogBodyEmptyBodyPrintsNothing(t *testing.T) {
	var buf bytes.Buffer
	printChangelogBody(&buf, "")

	if buf.Len() != 0 {
		t.Errorf("expected no output for an empty body, got %q", buf.String())
	}
}
