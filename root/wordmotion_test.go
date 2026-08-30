package root

import "testing"

func TestParseWordMotion(t *testing.T) {
	cases := []struct {
		input    string
		want     cursorMotion
		consumed int
	}{
		{"\x1b[1;5D", motionWordLeft, 6},
		{"\x1b[1;5C", motionWordRight, 6},
		{"\x1b[1;3D", motionWordLeft, 6},
		{"\x1b[5D", motionWordLeft, 4},
		{"\x1bOd", motionWordLeft, 3},
		{"\x1bOc", motionWordRight, 3},
		{"\x1b[D", motionNone, 0},
		{"\x1bOD", motionNone, 0},
		{"\x1b[1;2D", motionNone, 0},
		{"\x1b[3~", motionNone, 0},
		{"\x1b[1;5A", motionNone, 0},
		{"\x1b[", motionNone, 0},
	}
	for _, c := range cases {
		got, consumed := parseWordMotion([]byte(c.input))
		if got != c.want || consumed != c.consumed {
			t.Errorf("parseWordMotion(%q) = %v, %d; want %v, %d", c.input, got, consumed, c.want, c.consumed)
		}
	}
}

// The expected offsets were measured against real bash and zsh on a pty, with
// ctrl+left/right bound to backward-word/forward-word.
func TestWordMotionOffsetMatchesTheShells(t *testing.T) {
	const line = "echo alpha-beta gamma"

	cases := []struct {
		shell  string
		offset int
		motion cursorMotion
		want   int
	}{
		{"bash", 0, motionWordLeft, 5},
		{"bash", 5, motionWordLeft, 10},
		{"bash", 10, motionWordRight, 6},
		{"zsh", 0, motionWordLeft, 5},
		{"zsh", 5, motionWordLeft, 16},
		{"zsh", 16, motionWordRight, 5},
		{"bash", 21, motionWordLeft, 21},
		{"bash", 0, motionWordRight, 0},
		{"fish", 0, motionWordLeft, 5},
	}
	for _, c := range cases {
		if got := wordMotionOffset(c.shell, line, c.offset, c.motion); got != c.want {
			t.Errorf("wordMotionOffset(%s, offset=%d, %v) = %d; want %d", c.shell, c.offset, c.motion, got, c.want)
		}
	}
}

func TestWordMotionOffsetHandlesEmptyAndUnicode(t *testing.T) {
	if got := wordMotionOffset("bash", "", 0, motionWordLeft); got != 0 {
		t.Errorf("empty buffer = %d; want 0", got)
	}
	if got := wordMotionOffset("bash", "echo héllo wörld", 0, motionWordLeft); got != 5 {
		t.Errorf("unicode word = %d; want 5", got)
	}
}
