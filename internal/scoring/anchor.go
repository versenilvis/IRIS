package scoring

import (
	"strings"
)

type Anchor struct {
	buffer string
	prefix string
	head   string
}

func commandWord(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, ' '); i >= 0 {
		return s[:i]
	}
	return s
}

func NewAnchor(candidates []string, buffer string) Anchor {
	buffer = strings.TrimLeft(buffer, " ")
	if buffer == "" {
		return Anchor{}
	}

	lowerBuf := strings.ToLower(buffer)
	head := strings.ToLower(commandWord(buffer))
	headKnown := false

	for _, c := range candidates {
		cTrim := strings.TrimSpace(c)
		lowerC := strings.ToLower(cTrim)
		// tier 1: candidate starts with buffer and continues past it
		if len(lowerC) > len(lowerBuf) && strings.HasPrefix(lowerC, lowerBuf) {
			return Anchor{buffer: buffer, prefix: buffer}
		}
		// tier 2: first word matches an existing command
		if !headKnown && head != "" && lowerC != lowerBuf && strings.ToLower(commandWord(cTrim)) == head {
			headKnown = true
		}
	}

	if headKnown {
		return Anchor{buffer: buffer, head: head}
	}
	return Anchor{}
}

func (a Anchor) Allows(cmd string) bool {
	cmdTrim := strings.TrimSpace(cmd)
	lowerCmd := strings.ToLower(cmdTrim)
	switch {
	case a.prefix != "":
		return len(lowerCmd) > len(a.prefix) && strings.HasPrefix(lowerCmd, strings.ToLower(a.prefix))
	case a.head != "":
		return lowerCmd != strings.ToLower(a.buffer) && strings.ToLower(commandWord(cmdTrim)) == a.head
	default:
		return true
	}
}

func isSubsequenceWithGap(sub, full string, maxGap int) bool {
	subRunes := []rune(sub)
	fullRunes := []rune(full)
	if len(subRunes) == 0 {
		return true
	}
	if len(fullRunes) < len(subRunes) {
		return false
	}

	i := 0
	prevIdx := -1
	for j := 0; j < len(fullRunes) && i < len(subRunes); j++ {
		if subRunes[i] == fullRunes[j] {
			if prevIdx >= 0 && maxGap > 0 {
				if j-prevIdx-1 > maxGap {
					return false
				}
			}
			prevIdx = j
			i++
		}
	}
	return i == len(subRunes)
}
