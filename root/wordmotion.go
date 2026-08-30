package root

import (
	"strconv"
	"strings"
	"unicode"
)

type cursorMotion int

const (
	motionNone cursorMotion = iota
	motionWordLeft
	motionWordRight
)

// zsh counts these as word characters on top of the alphanumerics; readline
// counts only alphanumerics, which is also the closest fit for anything else.
const zshWordChars = "*?_-.[]~=/&;!#$%^(){}<>"

func parseWordMotion(input []byte) (cursorMotion, int) {
	if len(input) < 3 || input[0] != 0x1b {
		return motionNone, 0
	}

	// rxvt spells ctrl+left and ctrl+right as ESC O d and ESC O c
	if input[1] == 'O' {
		switch input[2] {
		case 'd':
			return motionWordLeft, 3
		case 'c':
			return motionWordRight, 3
		}
		return motionNone, 0
	}
	if input[1] != '[' {
		return motionNone, 0
	}

	end := 2
	for end < len(input) && (input[end] == ';' || (input[end] >= '0' && input[end] <= '9')) {
		end++
	}
	if end >= len(input) || !hasWordModifier(string(input[2:end])) {
		return motionNone, 0
	}

	switch input[end] {
	case 'D':
		return motionWordLeft, end + 1
	case 'C':
		return motionWordRight, end + 1
	}
	return motionNone, 0
}

// xterm encodes modifiers as 1 plus a bitmask, 2 for alt and 4 for ctrl. A bare
// arrow carries no parameter at all and is tracked one step at a time instead.
func hasWordModifier(params string) bool {
	if i := strings.LastIndex(params, ";"); i >= 0 {
		params = params[i+1:]
	}
	mod, err := strconv.Atoi(params)
	if err != nil {
		return false
	}
	return (mod-1)&(2|4) != 0
}

// wordMotionOffset returns where the cursor lands, still counted back from the
// end of the line, once the shell has done the move itself.
func wordMotionOffset(shellName, buf string, offset int, motion cursorMotion) int {
	runes := []rune(buf)
	pos := min(max(len(runes)-offset, 0), len(runes))
	word := func(r rune) bool { return isWordRune(shellName, r) }

	if motion == motionWordLeft {
		for pos > 0 && !word(runes[pos-1]) {
			pos--
		}
		for pos > 0 && word(runes[pos-1]) {
			pos--
		}
		return len(runes) - pos
	}

	// zsh stops at the start of the next word, readline at the end of this one
	if shellName == "zsh" {
		for pos < len(runes) && word(runes[pos]) {
			pos++
		}
		for pos < len(runes) && !word(runes[pos]) {
			pos++
		}
	} else {
		for pos < len(runes) && !word(runes[pos]) {
			pos++
		}
		for pos < len(runes) && word(runes[pos]) {
			pos++
		}
	}
	return len(runes) - pos
}

func isWordRune(shellName string, r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	return shellName == "zsh" && strings.ContainsRune(zshWordChars, r)
}

// parseLineReport reads what the shell says its line editor holds. zsh reports
// the whole buffer with the cursor's place in it; a shell without that hook
// only ever sends the text itself, and is taken to be at the end of it.
func parseLineReport(query string) (line string, cursorFromEnd int, ok bool) {
	rest, found := strings.CutPrefix(query, "IRIS_LINE:")
	if !found {
		return query, 0, true
	}

	head, buf, found := strings.Cut(rest, ":")
	if !found {
		return "", 0, false
	}
	left, err := strconv.Atoi(head)
	if err != nil {
		return "", 0, false
	}
	return buf, max(len([]rune(buf))-left, 0), true
}
