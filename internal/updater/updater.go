package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/versenilvis/iris/internal/config"
)

const InstallScriptURL = "https://raw.githubusercontent.com/versenilvis/iris/main/scripts/install.sh"

var ErrRateLimited = errors.New("rate limited by GitHub API")

type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	Prerelease  bool      `json:"prerelease"`
}

func ResolveInstallScriptURL() string {
	if url := os.Getenv("IRIS_INSTALL_URL"); url != "" {
		return url
	}
	return InstallScriptURL
}

func PerformUpdate(latest string, interactive bool) (string, error) {
	ctx := context.Background()
	if !interactive {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}

	cmdRun := exec.CommandContext(ctx, "sh", "-c", "curl -sSL "+ResolveInstallScriptURL()+" | sh")
	if config.Get().Updater.Channel == "nightly" {
		cmdRun.Env = append(os.Environ(), "IRIS_RELEASE_TAG="+latest)
	}

	if interactive {
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr
		cmdRun.Stdin = os.Stdin
		return "", cmdRun.Run()
	}

	out, runErr := cmdRun.CombinedOutput()
	return string(out), runErr
}

func NewGitHubRequestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func FetchGitHubBody(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func FetchLatestRelease() (Release, error) {
	ctx, cancel := NewGitHubRequestContext()
	defer cancel()

	endpoint := os.Getenv("IRIS_UPDATE_URL")
	if endpoint == "" {
		if config.Get().Updater.Channel == "nightly" {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases"
		} else {
			endpoint = "https://api.github.com/repos/versenilvis/iris/releases/latest"
		}
	}

	body, err := FetchGitHubBody(ctx, endpoint)
	if err != nil {
		return Release{}, err
	}

	if config.Get().Updater.Channel == "nightly" && os.Getenv("IRIS_UPDATE_URL") == "" {
		var releases []Release
		if err := json.Unmarshal(body, &releases); err != nil {
			return Release{}, err
		}
		if len(releases) == 0 {
			return Release{}, fmt.Errorf("no releases found")
		}
		return releases[0], nil
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return Release{}, err
	}
	if release.TagName == "" {
		return Release{}, fmt.Errorf("no tag_name in response")
	}
	return release, nil
}

func FetchLatestVersion() (string, error) {
	release, err := FetchLatestRelease()
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

func IsNewer(current, latest string) bool {
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")
	channel := config.Get().Updater.Channel
	// dev builds or empty versions never trigger an update
	if c == "" || c == "dev" || l == "" || l == "dev" {
		return false
	}

	// nightly builds are never shown as stable update targets
	if channel != "nightly" && strings.Contains(l, "-nightly.") {
		return false
	}

	if c == l {
		return false
	}

	cParts := strings.Split(c, ".")
	lParts := strings.Split(l, ".")

	// compare major.minor.patch
	for i := 0; i < len(cParts) && i < len(lParts); i++ {
		// strip pre-release tags like -beta or -rc for numeric comparison
		cClean := strings.Split(cParts[i], "-")[0]
		lClean := strings.Split(lParts[i], "-")[0]

		cv, _ := strconv.Atoi(cClean)
		lv, _ := strconv.Atoi(lClean)
		if lv > cv {
			return true
		}
		if lv < cv {
			return false
		}
	}

	if channel == "nightly" && strings.Contains(l, "-nightly.") && c != l {
		return true
	}

	// if all parts are equal, the one with more parts is newer (e.g. 1.0.1 > 1.0)
	return len(lParts) > len(cParts)
}
