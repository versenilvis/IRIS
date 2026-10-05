package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

func TestFnmMultishellNotDuplicatedInPath(t *testing.T) {
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not installed")
	}

	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	bin := binary(t)

	outerFnmPath := filepath.Join(home, ".local/state/fnm_multishells/5487_outer/bin")
	innerFnmPath := filepath.Join(home, ".local/state/fnm_multishells/5450_inner/bin")
	_ = os.MkdirAll(outerFnmPath, 0o755)
	_ = os.MkdirAll(innerFnmPath, 0o755)

	// simulate .zshrc where fnm runs before iris init
	zshrc := `
PROMPT='> '
RPROMPT=''
unsetopt PROMPT_SP
if [ -z "$IRIS_PID" ]; then
    export FNM_MULTISHELL_PATH="` + outerFnmPath + `"
    export PATH="$FNM_MULTISHELL_PATH:$PATH"
fi
eval "$(` + bin + ` init zsh)"
if [ -n "$IRIS_PID" ]; then
    export FNM_MULTISHELL_PATH="` + innerFnmPath + `"
    export PATH="$FNM_MULTISHELL_PATH:$PATH"
fi
`
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

	term := tuitest.StartT(t, []string{bin},
		tuitest.WithSize(cols, rows),
		tuitest.WithEnv(env...),
		tuitest.WithDir(home),
	)
	defer func() { _ = term.Close() }()

	if err := term.WaitForText(">", 20*time.Second); err != nil {
		t.Fatalf("iris never reached a prompt: %v\n%s", err, term.Snapshot())
	}

	outLog := filepath.Join(home, "path_out.txt")
	cmd := "echo $PATH > " + outLog + "\n"
	if err := term.SendKeys(cmd); err != nil {
		t.Fatal(err)
	}

	// wait for path dump file to be written by the shell
	var writtenPath string
	for range 20 {
		time.Sleep(100 * time.Millisecond)
		data, err := os.ReadFile(outLog)
		if err == nil && len(data) > 0 {
			writtenPath = string(data)
			break
		}
	}

	if writtenPath == "" {
		t.Fatalf("path file was not written, screen:\n%s", screen(term))
	}

	if strings.Contains(writtenPath, "5487_outer") {
		t.Fatalf("PATH should not contain dead outer fnm multishell (5487_outer), got:\n%s", writtenPath)
	}
	if !strings.Contains(writtenPath, "5450_inner") {
		t.Fatalf("PATH must contain active inner fnm multishell (5450_inner), got:\n%s", writtenPath)
	}
}

func TestRealFnmWithIris(t *testing.T) {
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not installed")
	}

	fnmBin := "/tmp/fnm-bin/fnm"
	if _, err := os.Stat(fnmBin); err != nil {
		if path, errLook := exec.LookPath("fnm"); errLook == nil {
			fnmBin = path
		} else {
			t.Skip("fnm binary not available")
		}
	}

	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	bin := binary(t)

	// exact reproduction steps from issue #170
	zshrc := `
PROMPT='> '
RPROMPT=''
unsetopt PROMPT_SP
eval "$(` + fnmBin + ` env --use-on-cd --shell zsh)"
eval "$(` + bin + ` init zsh)"
`
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

	term := tuitest.StartT(t, []string{bin},
		tuitest.WithSize(cols, rows),
		tuitest.WithEnv(env...),
		tuitest.WithDir(home),
	)
	defer func() { _ = term.Close() }()

	if err := term.WaitForText(">", 20*time.Second); err != nil {
		t.Fatalf("iris never reached a prompt: %v\n%s", err, term.Snapshot())
	}

	outLog := filepath.Join(home, "path_out.txt")
	cmd := "echo $PATH > " + outLog + "\n"
	if err := term.SendKeys(cmd); err != nil {
		t.Fatal(err)
	}

	var writtenPath string
	for range 20 {
		time.Sleep(100 * time.Millisecond)
		data, err := os.ReadFile(outLog)
		if err == nil && len(data) > 0 {
			writtenPath = string(data)
			break
		}
	}

	if writtenPath == "" {
		t.Fatalf("path file was not written, screen:\n%s", screen(term))
	}

	parts := strings.Split(strings.TrimSpace(writtenPath), ":")
	fnmCount := 0
	for _, p := range parts {
		if strings.Contains(p, "fnm_multishells") {
			fnmCount++
		}
	}

	if fnmCount != 1 {
		t.Fatalf("expected exactly 1 fnm_multishells entry in PATH, got %d:\n%s", fnmCount, writtenPath)
	}
}

