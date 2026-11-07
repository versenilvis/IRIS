package root

import "unicode"

type AutocloseAction int

const (
	AutocloseNone AutocloseAction = iota
	AutoclosePair
	AutocloseSkip
	AutocloseDeletePair
)

type AutocloseDecision struct {
	Action      AutocloseAction
	NewBuffer   string
	NewOffset   int
	PtySequence []byte
}

func matchingClosing(b byte) (byte, bool) {
	switch b {
	case '(':
		return ')', true
	case '[':
		return ']', true
	case '{':
		return '}', true
	case '"':
		return '"', true
	case '\'':
		return '\'', true
	case '`':
		return '`', true
	default:
		return 0, false
	}
}

func isCloseChar(b byte) bool {
	return b == ')' || b == ']' || b == '}' || b == '"' || b == '\'' || b == '`'
}

func isQuote(b byte) bool {
	return b == '"' || b == '\'' || b == '`'
}

func canAutocloseQuote(runes []rune, pos int, b byte) bool {
	if pos > 0 {
		prev := runes[pos-1]
		if unicode.IsLetter(prev) || unicode.IsDigit(prev) || prev == rune(b) {
			return false
		}
	}
	if pos < len(runes) {
		next := runes[pos]
		if unicode.IsLetter(next) || unicode.IsDigit(next) {
			return false
		}
	}
	return true
}

func isEmptyPair(left, right rune) bool {
	switch left {
	case '(':
		return right == ')'
	case '[':
		return right == ']'
	case '{':
		return right == '}'
	case '"':
		return right == '"'
	case '\'':
		return right == '\''
	case '`':
		return right == '`'
	default:
		return false
	}
}

func isEscaped(runes []rune, pos int) bool {
	count := 0
	for i := pos - 1; i >= 0 && runes[i] == '\\'; i-- {
		count++
	}
	return count%2 != 0
}

func countUnescaped(runes []rune, target rune, upTo int) int {
	count := 0
	for i := range upTo {
		if runes[i] == target && !isEscaped(runes, i) {
			count++
		}
	}
	return count
}

func CanAutoclose(enabled bool, inPaste bool, isAltScreen bool, isSynced bool) bool {
	return enabled && !inPaste && !isAltScreen && isSynced
}

func DecideAutoclose(b byte, buffer string, cursorOffset int) AutocloseDecision {
	runes := []rune(buffer)
	if cursorOffset < 0 || cursorOffset > len(runes) {
		return AutocloseDecision{Action: AutocloseNone}
	}

	pos := len(runes) - cursorOffset

	if b == 127 || b == 0x08 {
		if pos > 0 && pos < len(runes) && isEmptyPair(runes[pos-1], runes[pos]) {
			newRunes := make([]rune, 0, len(runes)-2)
			newRunes = append(newRunes, runes[:pos-1]...)
			newRunes = append(newRunes, runes[pos+1:]...)
			// delete both characters in shell line editor without relying on terminal delete binding
			return AutocloseDecision{
				Action:      AutocloseDeletePair,
				NewBuffer:   string(newRunes),
				NewOffset:   cursorOffset - 1,
				PtySequence: []byte("\x1b[C\x7f\x7f"),
			}
		}
		return AutocloseDecision{Action: AutocloseNone}
	}

	if isEscaped(runes, pos) {
		return AutocloseDecision{Action: AutocloseNone}
	}

	if isCloseChar(b) && pos < len(runes) && runes[pos] == rune(b) {
		shouldSkip := true
		if isQuote(b) {
			shouldSkip = countUnescaped(runes, rune(b), pos)%2 != 0
		}
		if shouldSkip {
			// step over closing character when typed directly before it
			return AutocloseDecision{
				Action:      AutocloseSkip,
				NewBuffer:   buffer,
				NewOffset:   cursorOffset - 1,
				PtySequence: []byte("\x1b[C"),
			}
		}
	}

	if isQuote(b) && !canAutocloseQuote(runes, pos, b) {
		// avoid autoclosing inside contractions or when touching alphanumeric text
		return AutocloseDecision{Action: AutocloseNone}
	}

	if closing, ok := matchingClosing(b); ok {
		newRunes := make([]rune, 0, len(runes)+2)
		newRunes = append(newRunes, runes[:pos]...)
		newRunes = append(newRunes, rune(b), rune(closing))
		newRunes = append(newRunes, runes[pos:]...)
		return AutocloseDecision{
			Action:      AutoclosePair,
			NewBuffer:   string(newRunes),
			NewOffset:   cursorOffset + 1,
			PtySequence: []byte{b, closing, 0x1b, '[', 'D'},
		}
	}

	return AutocloseDecision{Action: AutocloseNone}
}

func HandleAutoclose(b byte, buffer string, cursorOffset int, enabled bool, inPaste bool, isAltScreen bool, isSynced bool) AutocloseDecision {
	if !CanAutoclose(enabled, inPaste, isAltScreen, isSynced) {
		return AutocloseDecision{Action: AutocloseNone}
	}
	return DecideAutoclose(b, buffer, cursorOffset)
}
