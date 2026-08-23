package root

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/versenilvis/iris/internal/config"
)

func captureInitScript(t *testing.T, shell string) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	original := os.Stdout
	os.Stdout = w
	initCmd.Run(initCmd, []string{shell})
	_ = w.Close()
	os.Stdout = original

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// a tool runner that sources the rc file in a non-interactive shell must not
// exec a second iris, which would seize the tty from the one already running
func TestInitAutostartRequiresInteractiveShell(t *testing.T) {
	guards := map[string]string{
		"zsh":  "[[ -o interactive ]]",
		"bash": "[[ $- == *i* ]]",
		"fish": "status is-interactive",
	}

	for shell, guard := range guards {
		t.Run(shell, func(t *testing.T) {
			script := captureInitScript(t, shell)

			head, _, found := strings.Cut(script, "exec iris")
			if !found {
				t.Fatalf("%s init script has no autostart", shell)
			}

			cond := strings.LastIndex(head, "\nif ")
			if cond < 0 {
				t.Fatalf("%s autostart is not inside an if", shell)
			}

			if !strings.Contains(head[cond:], guard) {
				t.Fatalf("%s autostart is not guarded by %q:\n%s", shell, guard, head[cond:])
			}
		})
	}
}

func TestFishAutosuggestionsFollowGhostTextMode(t *testing.T) {
	const marker = "set -g fish_autosuggestion_enabled 0"

	cases := []struct {
		mode config.GhostTextMode
		want bool
	}{
		{config.GhostTextOff, false},
		{config.GhostTextOn, true},
		{config.GhostTextIndividual, true},
	}

	for _, tc := range cases {
		cfg := config.DefaultConfig()
		cfg.UI.GhostText = tc.mode
		config.Init(cfg)

		script := captureInitScript(t, "fish")
		if got := strings.Contains(script, marker); got != tc.want {
			t.Errorf("ghost-text=%d: fish autosuggestions disabled=%v, want %v", tc.mode, got, tc.want)
		}
	}
}
