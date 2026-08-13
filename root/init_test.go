package root

import (
	"io"
	"os"
	"strings"
	"testing"
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
