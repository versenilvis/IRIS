package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

// bashHook returns the part of the bash init script that runs inside iris.
func bashHook(t *testing.T) string {
	t.Helper()
	script := captureInitScript(t, "bash")
	_, hook, found := strings.Cut(script, "# Iris Autocomplete Hook")
	if !found {
		t.Fatal("bash init script has no autocomplete hook")
	}
	return hook
}

func TestBashHookKeepsBracketedPaste(t *testing.T) {
	bash, lookErr := exec.LookPath("bash")
	if lookErr != nil {
		t.Skip("bash not installed")
	}

	dir := t.TempDir()
	report := filepath.Join(dir, "report")
	ran := filepath.Join(dir, "ran")
	rc := filepath.Join(dir, "rc")
	rcBody := "PS1='$ '\nexec 9>" + report + "\nIRIS_PID=1\nIRIS_FD=9\n" + bashHook(t)
	if err := os.WriteFile(rc, []byte(rcBody), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext(t.Context(), bash, "--noprofile", "--rcfile", rc, "-i")
	cmd.Env = append(os.Environ(), "TERM=xterm", "INPUTRC=/dev/null")
	ptmx, startErr := pty.Start(cmd)
	if startErr != nil {
		t.Fatal(startErr)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		_ = ptmx.Close()
	}()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, err := ptmx.Read(buf); err != nil {
				return
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	// a two line paste: the first line must not run before the user hits enter
	paste := "\x1b[200~touch " + ran + "\necho second\x1b[201~"
	if _, err := ptmx.Write([]byte("true; " + paste)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(700 * time.Millisecond)

	if _, err := os.Stat(ran); err == nil {
		t.Fatal("a pasted line ran without enter: bracketed paste was taken over")
	}

	got, readErr := os.ReadFile(report)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Contains(got, []byte("IRIS_LINE:6:true; \x00")) {
		t.Fatalf("the line was not reported at paste start, got %q", got)
	}
}
