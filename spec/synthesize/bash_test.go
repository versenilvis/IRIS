package synthesize

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBashFile_Basic(t *testing.T) {
	absDir, err := filepath.Abs("testdata/completions")
	if err != nil {
		t.Fatal(err)
	}

	opts, err := parseBashFile(absDir, "toycmd", defaultMaxFileSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(opts) == 0 {
		t.Fatal("expected parsed options from bash completion script")
	}

	names := make(map[string]bool)
	for _, o := range opts {
		names[o.Name] = true
	}

	if !names["--file"] {
		t.Error("expected --file option")
	}
	if !names["--verbose"] {
		t.Error("expected --verbose option")
	}
	if !names["--dry-run"] {
		t.Error("expected --dry-run option")
	}
	if !names["--output"] {
		t.Error("expected --output option")
	}
}

func TestParseBashReader_CaseAliasesAndCompgen(t *testing.T) {
	script := `
_tool() {
    case $prev in
        '-s'|'--silent'|'-q'|'--quiet')
            return 0
            ;;
    esac
    COMPREPLY=( $(compgen -W "--color --version -v -h --help" -- "$cur") )
}
`
	opts, err := parseBashReader(strings.NewReader(script))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var foundSilent bool
	for _, o := range opts {
		if o.Name == "--silent" {
			foundSilent = true
			if len(o.Aliases) < 2 {
				t.Errorf("expected aliases for --silent, got: %v", o.Aliases)
			}
		}
	}
	if !foundSilent {
		t.Error("expected --silent with paired aliases from case branch")
	}
}
