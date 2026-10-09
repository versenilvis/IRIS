package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

func TestAutoExecuteExactMatch(t *testing.T) {
	for _, tt := range []struct {
		name       string
		coreConfig string
		burst      bool
		navigate   bool
		accept     bool
		want       string
	}{
		{"default-on", "auto-execute = true\n", false, false, false, "IRIS_EXACT"},
		{"auto-on", "auto-execute = true\nfilter-exact-match = \"auto\"\n", false, false, false, "IRIS_EXACT"},
		{"always-filter", "auto-execute = true\nfilter-exact-match = true\n", false, false, false, "IRIS_EXACT_WRONG"},
		{"never-filter", "auto-execute = false\nfilter-exact-match = false\n", false, false, true, "IRIS_EXACT"},
		{"auto-off", "auto-execute = false\nfilter-exact-match = \"auto\"\n", false, false, true, "IRIS_EXACT_WRONG"},
		{"final-byte-enter", "auto-execute = true\nfilter-exact-match = \"auto\"\n", true, false, false, "IRIS_EXACT"},
		{"navigate-auto-on", "auto-execute = true\nfilter-exact-match = \"auto\"\n", false, true, false, "IRIS_EXACT_WRONG"},
		{"navigate-auto-off", "auto-execute = false\nfilter-exact-match = false\n", false, true, false, "IRIS_EXACT_WRONG"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			term := startExactHistory(t, tt.coreConfig)
			defer func() { _ = term.Close() }()

			query := "echo IRIS_EXACT"
			if tt.burst {
				query = "echo IRIS_EXAC"
			}
			if err := term.Type(query); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitForText("Accept", 10*time.Second); err != nil {
				t.Fatalf("no menu: %v\n%s", err, screen(term))
			}
			if err := term.WaitStable(time.Second); err != nil {
				t.Fatal(err)
			}
			if tt.navigate {
				if err := term.SendKeys(tuitest.Down); err != nil {
					t.Fatal(err)
				}
				if err := term.WaitStable(time.Second); err != nil {
					t.Fatal(err)
				}
			}
			if tt.accept {
				if err := term.SendKeys(tuitest.Tab); err != nil {
					t.Fatal(err)
				}
				if err := term.WaitStable(time.Second); err != nil {
					t.Fatal(err)
				}
			}

			if tt.burst {
				if err := term.SendKeys("T", tuitest.Enter); err != nil {
					t.Fatal(err)
				}
			} else if err := term.SendKeys(tuitest.Enter); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitForText("\n"+tt.want+"\n", 10*time.Second); err != nil {
				t.Fatalf("expected output %q: %v\n%s", tt.want, err, screen(term))
			}
			for line := range strings.SplitSeq(screen(term), "\n") {
				if (line == "IRIS_EXACT" || line == "IRIS_EXACT_WRONG") && line != tt.want {
					t.Fatalf("executed %q; want %q\n%s", line, tt.want, screen(term))
				}
			}
		})
	}
}

func TestAutoExecuteEnterAcceptsFreshMidLineSuggestion(t *testing.T) {
	term := startExactHistory(t, "auto-execute = true\n")
	defer func() { _ = term.Close() }()

	if err := term.Type("echo IRIS_EXAC"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("no menu: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(time.Second); err != nil {
		t.Fatal(err)
	}
	if err := term.SendKeys(tuitest.Left); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(time.Second); err != nil {
		t.Fatal(err)
	}
	if err := term.SendKeys(tuitest.Enter); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("\nIRIS_EXACT_WRONG\n", 10*time.Second); err != nil {
		t.Fatalf("selected suggestion was not executed: %v\n%s", err, screen(term))
	}
}

func TestAutoExecuteEnterDoesNotRerunGenerators(t *testing.T) {
	home := ghostHome(t)
	binDir := filepath.Join(home, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	generatorLog := filepath.Join(home, "generator.log")
	script := `#!/bin/sh
if [ "$1" = ps ]; then
    printf 'query\n' >> "$IRIS_TEST_GENERATOR_LOG"
    printf 'irisfoo\nirisfoo_WRONG\n'
else
    printf 'IRIS_SUBMITTED %s\n' "$2"
fi
`
	if err := os.WriteFile(filepath.Join(binDir, "docker"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	histFile := filepath.Join(home, ".zsh_history")
	history := ": 1700000000:0;docker stop irisfoo\n: 1700000001:0;docker stop irisfoo_WRONG\n"
	if err := os.WriteFile(histFile, []byte(history), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "iris", "config.toml"), []byte("[core]\nauto-execute = true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"HISTFILE="+histFile,
		"PATH="+binDir+":"+os.Getenv("PATH"),
		"IRIS_TEST_GENERATOR_LOG="+generatorLog,
	)
	defer func() { _ = term.Close() }()

	if err := term.Type("docker stop irisfoo"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("no menu: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(time.Second); err != nil {
		t.Fatal(err)
	}
	before, readErr := os.ReadFile(generatorLog)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(before) == 0 {
		t.Fatal("generator was not called")
	}
	if err := term.SendKeys(tuitest.Enter); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("\nIRIS_SUBMITTED irisfoo\n", 10*time.Second); err != nil {
		t.Fatalf("selected command was not executed: %v\n%s", err, screen(term))
	}
	after, readErr := os.ReadFile(generatorLog)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(after) != string(before) {
		t.Fatalf("Enter reran the generator: %d calls before, %d after", len(strings.Fields(string(before))), len(strings.Fields(string(after))))
	}
}

func startExactHistory(t *testing.T, coreConfig string) *tuitest.Terminal {
	t.Helper()
	home := ghostHome(t)
	histFile := filepath.Join(home, ".zsh_history")
	history := ": 1700000000:0;echo IRIS_EXACT\n: 1700000001:0;echo IRIS_EXACT_WRONG\n"
	if err := os.WriteFile(histFile, []byte(history), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(home, ".config", "iris", "config.toml")
	if err := os.WriteFile(configPath, []byte("[core]\n"+coreConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	return startIn(t, home, "IRIS_CORE_MODE=history", "HISTFILE="+histFile)
}
