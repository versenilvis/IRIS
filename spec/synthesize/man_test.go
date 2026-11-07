package synthesize

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseManFile_Uncompressed(t *testing.T) {
	absDir, err := filepath.Abs("testdata/man1")
	if err != nil {
		t.Fatal(err)
	}

	desc, opts, err := parseManFile(absDir, "toycmd", defaultMaxFileSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if desc != "a demonstration command for iris" {
		t.Errorf("got desc %q, want 'a demonstration command for iris'", desc)
	}

	if len(opts) != 3 {
		t.Fatalf("expected 3 options, got %d", len(opts))
	}

	var fileOpt, verbOpt, quietOpt bool
	for _, o := range opts {
		switch o.Name {
		case "--file":
			fileOpt = true
			if len(o.Aliases) != 1 || o.Aliases[0] != "-f" {
				t.Errorf("expected -f alias for --file, got: %v", o.Aliases)
			}
			if !strings.Contains(o.Description, "Specify the input configuration file path") {
				t.Errorf("unexpected desc: %q", o.Description)
			}
		case "--verbose":
			verbOpt = true
			if len(o.Aliases) != 1 || o.Aliases[0] != "-v" {
				t.Errorf("expected -v alias for --verbose, got: %v", o.Aliases)
			}
		case "-q":
			quietOpt = true
			if !strings.Contains(o.Description, "Suppress standard informational output") {
				t.Errorf("unexpected desc: %q", o.Description)
			}
		}
	}

	if !fileOpt || !verbOpt || !quietOpt {
		t.Errorf("missing parsed options: file=%v, verb=%v, quiet=%v", fileOpt, verbOpt, quietOpt)
	}
}

func TestParseManFile_CompressedGzip(t *testing.T) {
	tmpDir := t.TempDir()
	manContent := []byte(`.TH GZCMD "1" "Nov 2026" "Gzcmd 1.0" "User Commands"
.SH NAME
gzcmd \- compressed manual test tool
.SH OPTIONS
.TP
\fB\-c\fP, \fB\-\-count\fP=\fIN\fP
Set iteration count.
.TP
\fB\-\-force\fP
Force operation without prompting.
`)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(manContent); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}

	gzPath := filepath.Join(tmpDir, "gzcmd.1.gz")
	if err := os.WriteFile(gzPath, buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}

	desc, opts, err := parseManFile(tmpDir, "gzcmd", defaultMaxFileSize)
	if err != nil {
		t.Fatalf("unexpected error parsing gzipped manpage: %v", err)
	}

	if desc != "compressed manual test tool" {
		t.Errorf("got desc %q, want 'compressed manual test tool'", desc)
	}
	if len(opts) != 2 {
		t.Fatalf("expected 2 options, got %d", len(opts))
	}
}

func TestParseManReader_GzipBombRejected(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	zeros := make([]byte, 64*1024)
	for range 32 {
		if _, err := gw.Write(zeros); err != nil {
			t.Fatal(err)
		}
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}

	_, _, err := parseManReader(&buf, "bomb", 16*1024)
	if err == nil {
		t.Fatal("expected gzip bomb to be rejected by decompressed size limit")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size limit") {
		t.Errorf("expected size limit error, got: %v", err)
	}
}

func TestParseManReader_MacrosAndMdoc(t *testing.T) {
	roff := `.Dd November 7, 2026
.Dt TESTTOOL 1
.Sh NAME
.Nm testtool
.Nd utility for testing roff parsing
.Sh OPTIONS
.It Fl a , Fl Fl all
Process all available items
.It Fl Fl dry-run
Simulate execution without modifying state
`
	desc, opts, err := parseManReader(strings.NewReader(roff), "testtool", defaultMaxFileSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if desc != "utility for testing roff parsing" {
		t.Errorf("got desc %q", desc)
	}

	if len(opts) != 2 {
		t.Fatalf("expected 2 options, got %d", len(opts))
	}
	if opts[0].Name != "--all" || len(opts[0].Aliases) != 1 || opts[0].Aliases[0] != "-a" {
		t.Errorf("unexpected option 0: %+v", opts[0])
	}
}

func TestParseManReader_ControlCharactersAndAnsiStripped(t *testing.T) {
	roff := ".TH BAD \"1\" \"Nov 2026\" \"Bad 1.0\" \"User Commands\"\n" +
		".SH NAME\n" +
		"bad \\- tool with \x1b[31mcolors\x1b[0m and \x1b]52;c;Y29weQ==\x07injection\n" +
		".SH OPTIONS\n" +
		".TP\n" +
		"\\fB\\-\\-clean\\fR, \\fB\\-c\\fR\n" +
		"Clean flag description with \x1b[1mbold\x1b[0m text.\n" +
		".TP\n" +
		"\\fB\\-\\-bad\x1b[31mflag\\fR\n" +
		"Flag with embedded escape code.\n"
	desc, opts, err := parseManReader(strings.NewReader(roff), "bad", defaultMaxFileSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(desc, "\x1b") || strings.Contains(desc, "52;c;") || strings.Contains(desc, "[31m") {
		t.Errorf("desc still contains escape codes: %q", desc)
	}
	if desc != "tool with colors and injection" {
		t.Errorf("unexpected cleaned desc: %q", desc)
	}

	for _, o := range opts {
		if strings.Contains(o.Name, "\x1b") || strings.Contains(o.Description, "\x1b") {
			t.Errorf("option contains escape: %+v", o)
		}
		if strings.Contains(o.Name, "31m") || strings.Contains(o.Description, "31m") {
			t.Errorf("option contains CSI code: %+v", o)
		}
		if strings.Contains(o.Name, "52;c;") || strings.Contains(o.Description, "52;c;") {
			t.Errorf("option contains OSC sequence: %+v", o)
		}
		if !reValidFlag.MatchString(o.Name) {
			t.Errorf("option name %q does not match valid flag regex", o.Name)
		}
		for _, a := range o.Aliases {
			if !reValidFlag.MatchString(a) {
				t.Errorf("alias %q does not match valid flag regex", a)
			}
		}
	}
}

func TestCleanDescription_RuneTruncation(t *testing.T) {
	// create a long unicode string with Vietnamese multibyte characters
	multibyte := strings.Repeat("Tiếng Việt rất phong phú và đa dạng ", 10)
	truncated := truncateRunes(multibyte, 160)
	runes := []rune(truncated)
	if len(runes) > 160 {
		t.Errorf("rune count %d exceeds 160", len(runes))
	}

	lines := []string{multibyte}
	cleaned := cleanDescription(lines)
	cleanedRunes := []rune(cleaned)
	if len(cleanedRunes) > 160 {
		t.Errorf("cleaned rune count %d exceeds 160", len(cleanedRunes))
	}
}
