package root

import "testing"

func TestLastAltScreenTransition(t *testing.T) {
	// the exact startup probe fish 4.8 writes, captured from a bare PTY: it
	// enters and leaves the alternate screen inside one write
	fishProbe := []byte("\x1b[?u\x1b[>0q\x1b]11;?\x1b\\\x1b[?1049h\x1bP+q696e646e\x1b\\\x1bP+q71756572792d6f732d6e616d65\x1b\\\x1b[?1049l\x1b[0c")

	tests := []struct {
		name      string
		chunk     []byte
		wantEnter bool
		wantFound bool
	}{
		{"fish startup probe leaves the main screen", fishProbe, false, true},
		{"tui takes over", []byte("\x1b[?1049h\x1b[2J"), true, true},
		{"tui gives the screen back", []byte("\x1b[?1049l\x1b[K"), false, true},
		{"legacy 47h", []byte("\x1b[?47h"), true, true},
		{"legacy 1047l", []byte("\x1b[?1047l"), false, true},
		{"plain output", []byte("hello\x1b[0m"), false, false},
		{"enter after a probe stays active", append(append([]byte{}, fishProbe...), []byte("\x1b[?1049h")...), true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enter, found := lastAltScreenTransition(tt.chunk)
			if found != tt.wantFound || enter != tt.wantEnter {
				t.Fatalf("got (enter=%v, found=%v), want (enter=%v, found=%v)", enter, found, tt.wantEnter, tt.wantFound)
			}
		})
	}
}

func TestLastAltScreenTransitionAcrossChunks(t *testing.T) {
	first := []byte("some output\x1b[?10")
	second := []byte("49h\x1b[2J")

	if _, found := lastAltScreenTransition(first); found {
		t.Fatal("half a sequence should not count as a transition")
	}

	carry := keepAltScreenCarry(first)
	enter, found := scanAltScreen(carry, second)
	if !found || !enter {
		t.Fatalf("split sequence not detected: enter=%v found=%v", enter, found)
	}
}

func TestScanAltScreenPrefersTheLaterSwitch(t *testing.T) {
	// a TUI restoring the main screen right where the previous read ended, then
	// a second one taking over inside the same read
	carry := keepAltScreenCarry([]byte("frame\x1b[?104"))
	enter, found := scanAltScreen(carry, []byte("9l\x1b[K\x1b[?1049h"))
	if !found || !enter {
		t.Fatalf("later switch ignored: enter=%v found=%v", enter, found)
	}

	enter, found = scanAltScreen(carry, []byte("9l\x1b[K"))
	if !found || enter {
		t.Fatalf("boundary exit not seen: enter=%v found=%v", enter, found)
	}
}

func TestKeepAltScreenCarryDoesNotAliasTheReadBuffer(t *testing.T) {
	buf := []byte("abcdefghij")
	carry := keepAltScreenCarry(buf)
	copy(buf, "0000000000")
	if string(carry) != "defghij" {
		t.Fatalf("carry changed with the read buffer: %q", carry)
	}
}
