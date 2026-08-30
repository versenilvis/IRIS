package tui

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

const (
	cols = 197
	rows = 24
)

var (
	buildOnce sync.Once
	irisBin   string
	errBuild  error
)

// binary builds iris once per run. The tests drive the real wrapper, so there
// is no substitute for the actual binary.
func binary(t *testing.T) string {
	t.Helper()
	// IRIS_TUI_BIN points the tests at a binary built elsewhere, so the same
	// scenarios can be replayed against an older commit to confirm they fail
	if custom := os.Getenv("IRIS_TUI_BIN"); custom != "" {
		return custom
	}
	buildOnce.Do(func() {
		// go test caches results per package, and it cannot see that these
		// tests build and run a binary. Reading the sources registers them as
		// test inputs so editing any of them invalidates the cached result.
		registerSourcesAsInputs(t)

		dir, err := os.MkdirTemp("", "iris-tui-*")
		if err != nil {
			errBuild = err
			return
		}
		irisBin = filepath.Join(dir, "iris")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "go", "build", "-o", irisBin, "github.com/versenilvis/iris/cmd/iris")
		cmd.Dir = repoRoot(t)
		if out, err := cmd.CombinedOutput(); err != nil {
			errBuild = err
			t.Logf("go build: %s", out)
		}
	})
	if errBuild != nil {
		t.Fatalf("building iris: %v", errBuild)
	}
	return irisBin
}

func registerSourcesAsInputs(t *testing.T) {
	t.Helper()
	root := repoRoot(t)
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "testdata", "docs", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			_, _ = os.ReadFile(path)
		}
		return nil
	})
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// start brings up iris wrapping zsh on its own pty, with every path pointed at
// a temp dir so the test never reads or writes the developer's real config,
// state or history.
func start(t *testing.T, extraEnv ...string) *tuitest.Terminal {
	t.Helper()

	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return startIn(t, home, extraEnv...)
}

func startIn(t *testing.T, home string, extraEnv ...string) *tuitest.Terminal {
	t.Helper()

	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not installed")
	}

	bin := binary(t)

	// a bare prompt keeps the geometry assertions readable, and the iris
	// integration has to be sourced the way a real .zshrc sources it
	prompt := os.Getenv("IRIS_TUI_PROMPT")
	if prompt == "" {
		prompt = "> "
	}
	// sourced last so a test can add its own bindkeys on top of the integration
	zshrc := "PROMPT='" + prompt + "'\nRPROMPT=''\nunsetopt PROMPT_SP\neval \"$(" + bin + " init zsh)\"\n" +
		"[[ -f $ZDOTDIR/.zshrc.extra ]] && source $ZDOTDIR/.zshrc.extra\n"
	if err := os.WriteFile(filepath.Join(home, ".zshrc"), []byte(zshrc), 0o644); err != nil {
		t.Fatal(err)
	}

	env := []string{
		"HOME=" + home,
		"ZDOTDIR=" + home,
		"XDG_CONFIG_HOME=" + filepath.Join(home, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(home, ".local/share"),
		"XDG_CACHE_HOME=" + filepath.Join(home, ".cache"),
		"SHELL=/bin/zsh",
		"IRIS_ACTIVE_SHELL=zsh",
		"PATH=" + os.Getenv("PATH"),
		"TERM=xterm-256color",
	}
	env = append(env, extraEnv...)

	term := tuitest.StartT(t, []string{bin},
		tuitest.WithSize(cols, rows),
		tuitest.WithEnv(env...),
		tuitest.WithDir(home),
	)

	if err := term.WaitForText(">", 20*time.Second); err != nil {
		t.Fatalf("iris never reached a prompt: %v\n%s", err, term.Snapshot())
	}
	return term
}

// screen returns the visible screen with trailing blank lines and trailing
// spaces removed, which is what the assertions care about.
func screen(term *tuitest.Terminal) string {
	lines := strings.Split(term.Snapshot(), "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

func TestIrisStartsAndShowsAPrompt(t *testing.T) {
	term := start(t)
	defer func() { _ = term.Close() }()

	if got := screen(term); !strings.Contains(got, ">") {
		t.Fatalf("no prompt on screen:\n%s", got)
	}
}
