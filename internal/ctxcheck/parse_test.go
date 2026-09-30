package ctxcheck

import (
	"reflect"
	"testing"
)

func TestParse_PosixBasicAndSeparators(t *testing.T) {
	cmd := "git add . && git commit -m 'initial commit' ; echo done"
	parsed := Parse(cmd, Posix)
	if parsed.IsUnknown {
		t.Fatalf("expected command not to be unknown")
	}
	if len(parsed.Segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(parsed.Segments))
	}
	expected := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "initial commit"},
		{"echo", "done"},
	}
	for i, seg := range parsed.Segments {
		if !reflect.DeepEqual(seg.Tokens, expected[i]) {
			t.Errorf("seg %d: expected %v, got %v", i, expected[i], seg.Tokens)
		}
	}
}

func TestParse_SingleQuoteBackslash(t *testing.T) {
	// in posix, backslash in '...' is literal
	posixCmd := `echo 'foo\bar'`
	pParsed := Parse(posixCmd, Posix)
	if len(pParsed.Segments) != 1 || pParsed.Segments[0].Tokens[1] != `foo\bar` {
		t.Fatalf("posix expected 'foo\\bar', got %v", pParsed.Segments[0].Tokens)
	}

	// in fish, \' escapes single quote
	fishCmd := `echo 'foo\'bar'`
	fParsed := Parse(fishCmd, Fish)
	if len(fParsed.Segments) != 1 || fParsed.Segments[0].Tokens[1] != `foo'bar` {
		t.Fatalf("fish expected 'foo'bar', got %v", fParsed.Segments[0].Tokens)
	}
}

func TestParse_FishAndOr(t *testing.T) {
	cmd := "git add . and git commit -m test or echo failed"
	parsed := Parse(cmd, Fish)
	if len(parsed.Segments) != 3 {
		t.Fatalf("fish expected 3 segments, got %d", len(parsed.Segments))
	}
	if parsed.Segments[0].Tokens[0] != "git" || parsed.Segments[1].Tokens[0] != "git" || parsed.Segments[2].Tokens[0] != "echo" {
		t.Fatalf("unexpected segments: %v", parsed.Segments)
	}

	// in posix, 'and' and 'or' are normal arguments
	pParsed := Parse("echo a and echo b", Posix)
	if len(pParsed.Segments) != 1 {
		t.Fatalf("posix should treat 'and' as argument, got %d segments", len(pParsed.Segments))
	}
}

func TestParse_Substitutions(t *testing.T) {
	cases := []struct {
		cmd     string
		dialect Dialect
		unknown bool
	}{
		{"echo $(whoami)", Posix, true},
		{"echo `whoami`", Posix, true},
		{"echo \"$(whoami)\"", Posix, true},
		{"echo (whoami)", Fish, true},
		{"echo (whoami)", Posix, false}, // in posix outside quote '(' is subshell/syntax, handled as token
		{"echo '(whoami)'", Fish, false},
		{"echo '$(whoami)'", Posix, false},
	}

	for _, tc := range cases {
		p := Parse(tc.cmd, tc.dialect)
		isUnk := p.IsUnknown || (len(p.Segments) > 0 && p.Segments[0].IsUnknown)
		if isUnk != tc.unknown {
			t.Errorf("cmd %q (dialect %v): expected unknown %v, got %v", tc.cmd, tc.dialect, tc.unknown, isUnk)
		}
	}
}

func TestParse_Redirects(t *testing.T) {
	// 2>&1 is allowed
	pOk := Parse("cmd 2>&1", Posix)
	if pOk.IsUnknown || (len(pOk.Segments) > 0 && pOk.Segments[0].IsUnknown) {
		t.Errorf("2>&1 should be allowed, got unknown")
	}

	// other redirects -> unknown
	badRedirects := []string{
		"cmd > out.txt",
		"cmd >> out.txt",
		"cmd < in.txt",
		"cmd 2> err.log",
	}
	for _, br := range badRedirects {
		p := Parse(br, Posix)
		isUnk := p.IsUnknown || (len(p.Segments) > 0 && p.Segments[0].IsUnknown)
		if !isUnk {
			t.Errorf("expected redirect %q to be unknown", br)
		}
	}
}

func TestParse_CwdChanges(t *testing.T) {
	compoundCases := []string{
		"cd foo && just test",
		"just test && cd foo",
		"pushd /tmp ; make",
		"popd && npm test",
	}
	for _, c := range compoundCases {
		p := Parse(c, Posix)
		if !p.IsUnknown {
			t.Errorf("expected compound command with cwd change to be unknown: %q", c)
		}
	}

	singleCdCases := []string{
		"cd /tmp",
		"cd ~",
		"cd -",
		"cd nonexistent",
	}
	for _, c := range singleCdCases {
		p := Parse(c, Posix)
		if p.IsUnknown {
			t.Errorf("expected single cd command not to be marked unknown by parser: %q", c)
		}
		if len(p.Segments) != 1 || p.Segments[0].Tokens[0] != "cd" {
			t.Errorf("expected 1 segment with 'cd' for %q, got %v", c, p.Segments)
		}
	}
}

func TestParse_EnvAndWrappers(t *testing.T) {
	// env variables stripped
	p1 := Parse("FOO=1 BAR=2 just test", Posix)
	if len(p1.Segments) != 1 || p1.Segments[0].Tokens[0] != "just" {
		t.Fatalf("expected 'just', got %v", p1.Segments[0].Tokens)
	}

	// wrappers stripped
	p2 := Parse("sudo env time nohup just test", Posix)
	if len(p2.Segments) != 1 || p2.Segments[0].Tokens[0] != "just" {
		t.Fatalf("expected 'just', got %v", p2.Segments[0].Tokens)
	}

	// wrapper with flag -> unknown
	p3 := Parse("sudo -u admin just test", Posix)
	if len(p3.Segments) == 0 || !p3.Segments[0].IsUnknown {
		t.Fatalf("expected wrapper with flag to be unknown, got %v", p3.Segments)
	}
}
