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
	return startInShell(t, home, "zsh", extraEnv...)
}

func startInShell(t *testing.T, home, shellName string, extraEnv ...string) *tuitest.Terminal {
	return startInDirShell(t, home, home, shellName, extraEnv...)
}

func startInDirShell(t *testing.T, home, workDir, shellName string, extraEnv ...string) *tuitest.Terminal {
	t.Helper()

	shellBin, err := exec.LookPath(shellName)
	if err != nil {
		if os.Getenv("IRIS_REQUIRE_SHELLS") == "1" {
			t.Fatalf("required shell %s not installed", shellName)
		}
		t.Skipf("%s not installed", shellName)
	}

	bin := binary(t)

	// a bare prompt keeps the geometry assertions readable
	prompt := os.Getenv("IRIS_TUI_PROMPT")
	if prompt == "" {
		prompt = "> "
	}

	switch shellName {
	case "zsh":
		// prevent system zshrc from prompting compinit in test pty
		zshenv := "unsetopt GLOBAL_RCS\nskip_global_compinit=1\n"
		if err := os.WriteFile(filepath.Join(home, ".zshenv"), []byte(zshenv), 0o644); err != nil {
			t.Fatal(err)
		}
		zshrc := "PROMPT='" + prompt + "'\nRPROMPT=''\nunsetopt PROMPT_SP\neval \"$(" + bin + " init zsh)\"\n" +
			"[[ -f $ZDOTDIR/.zshrc.extra ]] && source $ZDOTDIR/.zshrc.extra\n"
		if err := os.WriteFile(filepath.Join(home, ".zshrc"), []byte(zshrc), 0o644); err != nil {
			t.Fatal(err)
		}
	case "bash":
		bashrc := "PS1='" + prompt + "'\neval \"$(" + bin + " init bash)\"\n" +
			"[[ -f $HOME/.bashrc.extra ]] && source $HOME/.bashrc.extra\n"
		if err := os.WriteFile(filepath.Join(home, ".bashrc"), []byte(bashrc), 0o644); err != nil {
			t.Fatal(err)
		}
	case "fish":
		fishConfDir := filepath.Join(home, ".config/fish")
		_ = os.MkdirAll(fishConfDir, 0o755)
		configFish := "function fish_prompt\n    echo -n '" + prompt + "'\nend\n" +
			"function fish_update_completions\n    return 0\nend\n" +
			bin + " init fish | source\n"
		if err := os.WriteFile(filepath.Join(fishConfDir, "config.fish"), []byte(configFish), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	cacheDir := filepath.Join(os.TempDir(), "iris-tui-cache")
	_ = os.MkdirAll(cacheDir, 0o755)

	binDir := filepath.Dir(bin)
	env := []string{
		"HOME=" + home,
		"ZDOTDIR=" + home,
		"XDG_CONFIG_HOME=" + filepath.Join(home, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(home, ".local/share"),
		"XDG_CACHE_HOME=" + cacheDir,
		"SHELL=" + shellBin,
		"IRIS_ACTIVE_SHELL=" + shellName,
		"PATH=" + binDir + ":" + os.Getenv("PATH"),
		"TERM=xterm-256color",
		"skip_global_compinit=1",
	}
	env = append(env, extraEnv...)

	term := tuitest.StartT(t, []string{bin},
		tuitest.WithSize(cols, rows),
		tuitest.WithEnv(env...),
		tuitest.WithDir(workDir),
	)

	if err := term.WaitForText(">", 20*time.Second); err != nil {
		t.Fatalf("iris never reached a prompt in %s: %v\n%s", shellName, err, term.Snapshot())
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
