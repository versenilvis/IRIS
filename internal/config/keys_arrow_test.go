package config

import "testing"

func TestMatchKeyCtrlArrow(t *testing.T) {
	if m, _ := MatchKey([]byte("\x1b[B"), "down"); !m {
		t.Error("plain 'down' binding should match Down arrow")
	}
	if m, _ := MatchKey([]byte("\x1b[1;5B"), "down"); !m {
		t.Error("plain 'down' binding should match Ctrl+Down too")
	}
	if m, _ := MatchKey([]byte("\x1b[1;5B"), "ctrl+down"); !m {
		t.Error("'ctrl+down' should match Ctrl+Down")
	}
	if m, _ := MatchKey([]byte("\x1b[B"), "ctrl+down"); m {
		t.Error("'ctrl+down' must NOT match plain Down")
	}
	if m, _ := MatchKey([]byte("\x1b[1;5A"), "ctrl+up"); !m {
		t.Error("'ctrl+up' should match Ctrl+Up")
	}
	if m, _ := MatchKey([]byte("\x1b[1;5C"), "ctrl+right"); !m {
		t.Error("'ctrl+right' should match Ctrl+Right")
	}
	if m, _ := MatchKey([]byte("\x1b[1;5D"), "ctrl+left"); !m {
		t.Error("'ctrl+left' should match Ctrl+Left")
	}
}
