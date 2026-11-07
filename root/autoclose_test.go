package root

import (
	"bytes"
	"testing"
)

func TestAutoclose_OpenPairs(t *testing.T) {
	tests := []struct {
		name       string
		b          byte
		buf        string
		offset     int
		wantBuf    string
		wantOffset int
		wantPty    []byte
	}{
		{
			name:       "open parenthesis on empty buffer",
			b:          '(',
			buf:        "",
			offset:     0,
			wantBuf:    "()",
			wantOffset: 1,
			wantPty:    []byte("()\x1b[D"),
		},
		{
			name:       "open bracket on non-empty buffer",
			b:          '[',
			buf:        "echo ",
			offset:     0,
			wantBuf:    "echo []",
			wantOffset: 1,
			wantPty:    []byte("[]\x1b[D"),
		},
		{
			name:       "open brace in middle of buffer",
			b:          '{',
			buf:        "foobar",
			offset:     3,
			wantBuf:    "foo{}bar",
			wantOffset: 4,
			wantPty:    []byte("{}\x1b[D"),
		},
		{
			name:       "open double quote",
			b:          '"',
			buf:        "git commit -m ",
			offset:     0,
			wantBuf:    "git commit -m \"\"",
			wantOffset: 1,
			wantPty:    []byte("\"\"\x1b[D"),
		},
		{
			name:       "open single quote",
			b:          '\'',
			buf:        "",
			offset:     0,
			wantBuf:    "''",
			wantOffset: 1,
			wantPty:    []byte("''\x1b[D"),
		},
		{
			name:       "open backtick",
			b:          '`',
			buf:        "echo ",
			offset:     0,
			wantBuf:    "echo ``",
			wantOffset: 1,
			wantPty:    []byte("``\x1b[D"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := DecideAutoclose(tt.b, tt.buf, tt.offset)
			if dec.Action != AutoclosePair {
				t.Fatalf("expected AutoclosePair, got %v", dec.Action)
			}
			if dec.NewBuffer != tt.wantBuf {
				t.Errorf("NewBuffer = %q, want %q", dec.NewBuffer, tt.wantBuf)
			}
			if dec.NewOffset != tt.wantOffset {
				t.Errorf("NewOffset = %d, want %d", dec.NewOffset, tt.wantOffset)
			}
			if !bytes.Equal(dec.PtySequence, tt.wantPty) {
				t.Errorf("PtySequence = %q, want %q", dec.PtySequence, tt.wantPty)
			}
		})
	}
}

func TestAutoclose_SkipClosing(t *testing.T) {
	tests := []struct {
		name       string
		b          byte
		buf        string
		offset     int
		wantOffset int
	}{
		{
			name:       "skip parenthesis",
			b:          ')',
			buf:        "()",
			offset:     1,
			wantOffset: 0,
		},
		{
			name:       "skip bracket",
			b:          ']',
			buf:        "foo[]bar",
			offset:     4,
			wantOffset: 3,
		},
		{
			name:       "skip brace",
			b:          '}',
			buf:        "{}",
			offset:     1,
			wantOffset: 0,
		},
		{
			name:       "skip double quote when inside string",
			b:          '"',
			buf:        `echo "hello"`,
			offset:     1,
			wantOffset: 0,
		},
		{
			name:       "skip single quote when inside string",
			b:          '\'',
			buf:        `'test'`,
			offset:     1,
			wantOffset: 0,
		},
		{
			name:       "skip backtick when inside string",
			b:          '`',
			buf:        "``",
			offset:     1,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := DecideAutoclose(tt.b, tt.buf, tt.offset)
			if dec.Action != AutocloseSkip {
				t.Fatalf("expected AutocloseSkip, got %v", dec.Action)
			}
			if dec.NewBuffer != tt.buf {
				t.Errorf("NewBuffer = %q, want unchanged %q", dec.NewBuffer, tt.buf)
			}
			if dec.NewOffset != tt.wantOffset {
				t.Errorf("NewOffset = %d, want %d", dec.NewOffset, tt.wantOffset)
			}
			if !bytes.Equal(dec.PtySequence, []byte("\x1b[C")) {
				t.Errorf("PtySequence = %q, want %q", dec.PtySequence, "\x1b[C")
			}
		})
	}
}

func TestAutoclose_QuoteOpenRatherThanSkipOutside(t *testing.T) {
	// when cursor is before opening quote, count of quotes before is even so typing quote opens pair
	buf := `"hello"`
	offset := 7 // at pos 0, before first quote
	dec := DecideAutoclose('"', buf, offset)
	if dec.Action != AutoclosePair {
		t.Fatalf("expected AutoclosePair when typing quote before string, got %v", dec.Action)
	}
	if dec.NewBuffer != `"""hello"` {
		t.Errorf("NewBuffer = %q, want %q", dec.NewBuffer, `"""hello"`)
	}
	if dec.NewOffset != 8 {
		t.Errorf("NewOffset = %d, want 8", dec.NewOffset)
	}
}

func TestAutoclose_DeleteEmptyPair(t *testing.T) {
	tests := []struct {
		name       string
		b          byte
		buf        string
		offset     int
		wantBuf    string
		wantOffset int
	}{
		{
			name:       "delete parenthesis pair with del",
			b:          127,
			buf:        "()",
			offset:     1,
			wantBuf:    "",
			wantOffset: 0,
		},
		{
			name:       "delete parenthesis pair with bs",
			b:          0x08,
			buf:        "()",
			offset:     1,
			wantBuf:    "",
			wantOffset: 0,
		},
		{
			name:       "delete bracket pair",
			b:          127,
			buf:        "foo[]bar",
			offset:     4,
			wantBuf:    "foobar",
			wantOffset: 3,
		},
		{
			name:       "delete brace pair",
			b:          127,
			buf:        "{}",
			offset:     1,
			wantBuf:    "",
			wantOffset: 0,
		},
		{
			name:       "delete double quote pair",
			b:          127,
			buf:        `echo ""`,
			offset:     1,
			wantBuf:    "echo ",
			wantOffset: 0,
		},
		{
			name:       "delete single quote pair",
			b:          127,
			buf:        "''",
			offset:     1,
			wantBuf:    "",
			wantOffset: 0,
		},
		{
			name:       "delete backtick pair",
			b:          127,
			buf:        "``",
			offset:     1,
			wantBuf:    "",
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := DecideAutoclose(tt.b, tt.buf, tt.offset)
			if dec.Action != AutocloseDeletePair {
				t.Fatalf("expected AutocloseDeletePair, got %v", dec.Action)
			}
			if dec.NewBuffer != tt.wantBuf {
				t.Errorf("NewBuffer = %q, want %q", dec.NewBuffer, tt.wantBuf)
			}
			if dec.NewOffset != tt.wantOffset {
				t.Errorf("NewOffset = %d, want %d", dec.NewOffset, tt.wantOffset)
			}
			if !bytes.Equal(dec.PtySequence, []byte("\x1b[C\x7f\x7f")) {
				t.Errorf("PtySequence = %q, want %q", dec.PtySequence, "\x1b[C\x7f\x7f")
			}
		})
	}
}

func TestAutoclose_BackspaceNonEmpty(t *testing.T) {
	tests := []struct {
		name   string
		b      byte
		buf    string
		offset int
	}{
		{
			name:   "non-empty parenthesis",
			b:      127,
			buf:        "(a)",
			offset: 1, // cursor between 'a' and ')'
		},
		{
			name:   "mismatched brackets",
			b:      127,
			buf:        "(]",
			offset: 1,
		},
		{
			name:   "cursor at end",
			b:      127,
			buf:        "()",
			offset: 0,
		},
		{
			name:   "empty buffer",
			b:      127,
			buf:        "",
			offset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := DecideAutoclose(tt.b, tt.buf, tt.offset)
			if dec.Action != AutocloseNone {
				t.Errorf("expected AutocloseNone, got %v", dec.Action)
			}
		})
	}
}

func TestAutoclose_EscapedCharacters(t *testing.T) {
	// escaped quote should not be autoclosed into pair
	dec := DecideAutoclose('"', `echo \`, 0)
	if dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for escaped quote, got %v", dec.Action)
	}

	// escaped parenthesis should not be autoclosed
	dec = DecideAutoclose('(', `echo \`, 0)
	if dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for escaped paren, got %v", dec.Action)
	}

	// escaped backslash before quote means quote is not escaped
	dec = DecideAutoclose('"', `echo \\`, 0)
	if dec.Action != AutoclosePair {
		t.Errorf("expected AutoclosePair when backslash is escaped, got %v", dec.Action)
	}
}

func TestAutoclose_Conditions(t *testing.T) {
	// disabled config
	if dec := HandleAutoclose('(', "", 0, false, false, false, true); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone when disabled")
	}

	// during paste
	if dec := HandleAutoclose('(', "", 0, true, true, false, true); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone during paste")
	}

	// during altscreen
	if dec := HandleAutoclose('(', "", 0, true, false, true, true); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone during altscreen")
	}

	// buffer not synced
	if dec := HandleAutoclose('(', "", 0, true, false, false, false); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone when buffer not synced")
	}

	// out of bounds offset
	if dec := DecideAutoclose('(', "abc", -1); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for negative offset")
	}
	if dec := DecideAutoclose('(', "abc", 10); dec.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for offset > len")
	}
}

func TestAutoclose_UTF8Runes(t *testing.T) {
	buf := "tiếng Việt ()"
	offset := 1 // cursor between ()
	dec := DecideAutoclose(127, buf, offset)
	if dec.Action != AutocloseDeletePair {
		t.Fatalf("expected AutocloseDeletePair with UTF-8 runes, got %v", dec.Action)
	}
	if dec.NewBuffer != "tiếng Việt " {
		t.Errorf("NewBuffer = %q, want %q", dec.NewBuffer, "tiếng Việt ")
	}
	if dec.NewOffset != 0 {
		t.Errorf("NewOffset = %d, want 0", dec.NewOffset)
	}
}

func TestAutoclose_Contractions(t *testing.T) {
	// typing single quote in don should not autoclose
	decDont := DecideAutoclose('\'', "don", 0)
	if decDont.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for don', got %v", decDont.Action)
	}

	// typing single quote in it should not autoclose
	decIts := DecideAutoclose('\'', "it", 0)
	if decIts.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone for it', got %v", decIts.Action)
	}

	// typing quote right before a letter should not autoclose
	decBeforeWord := DecideAutoclose('"', "hello", 5) // cursor before 'h'
	if decBeforeWord.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone before letter, got %v", decBeforeWord.Action)
	}

	// typing quote after a digit should not autoclose
	decAfterDigit := DecideAutoclose('\'', "item5", 0)
	if decAfterDigit.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone after digit, got %v", decAfterDigit.Action)
	}

	// typing quote after identical quote should not autoclose
	decSameQuote := DecideAutoclose('"', `""`, 0)
	if decSameQuote.Action != AutocloseNone {
		t.Errorf("expected AutocloseNone after identical quote, got %v", decSameQuote.Action)
	}
}

func TestAutoclose_EchoQuotesFlow(t *testing.T) {
	// 1. type double quote after echo
	step1 := DecideAutoclose('"', "echo ", 0)
	if step1.Action != AutoclosePair || step1.NewBuffer != `echo ""` || step1.NewOffset != 1 {
		t.Fatalf("step 1 failed: %+v", step1)
	}

	// 2. cursor inside echo "a", type closing double quote to skip
	step2 := DecideAutoclose('"', `echo "a"`, 1)
	if step2.Action != AutocloseSkip || step2.NewOffset != 0 {
		t.Fatalf("step 2 failed: %+v", step2)
	}

	// 3. cursor after echo "a" , type single quote
	step3 := DecideAutoclose('\'', `echo "a" `, 0)
	if step3.Action != AutoclosePair || step3.NewBuffer != `echo "a" ''` || step3.NewOffset != 1 {
		t.Fatalf("step 3 failed: %+v", step3)
	}

	// 4. cursor inside echo "a" 'b', type closing single quote to skip
	step4 := DecideAutoclose('\'', `echo "a" 'b'`, 1)
	if step4.Action != AutocloseSkip || step4.NewOffset != 0 {
		t.Fatalf("step 4 failed: %+v", step4)
	}
}

func TestAutoclose_NestedQuotes(t *testing.T) {
	// single quote inside double quotes: echo "" with cursor at offset 1
	innerSingle := DecideAutoclose('\'', `echo ""`, 1)
	if innerSingle.Action != AutoclosePair || innerSingle.NewBuffer != `echo "''"` || innerSingle.NewOffset != 2 {
		t.Fatalf("nested single inside double failed: %+v", innerSingle)
	}

	// double quote inside single quotes: echo '' with cursor at offset 1
	innerDouble := DecideAutoclose('"', `echo ''`, 1)
	if innerDouble.Action != AutoclosePair || innerDouble.NewBuffer != `echo '""'` || innerDouble.NewOffset != 2 {
		t.Fatalf("nested double inside single failed: %+v", innerDouble)
	}

	// skip inner closing single quote in "''"
	skipInnerSingle := DecideAutoclose('\'', `"''"`, 2)
	if skipInnerSingle.Action != AutocloseSkip || skipInnerSingle.NewOffset != 1 {
		t.Fatalf("skip inner single failed: %+v", skipInnerSingle)
	}

	// skip outer closing double quote in "''"
	skipOuterDouble := DecideAutoclose('"', `"''"`, 1)
	if skipOuterDouble.Action != AutocloseSkip || skipOuterDouble.NewOffset != 0 {
		t.Fatalf("skip outer double failed: %+v", skipOuterDouble)
	}
}
