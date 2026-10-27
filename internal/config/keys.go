package config

import (
	"strconv"
	"strings"
)

// MatchKey matches a sequence of terminal bytes against a configured key string.
// Key strings can be e.g. "ctrl+r", "tab", "shift+tab", "right".
// Returns matched=true if the byte sequence matches the configured key, and consumed=number of bytes to consume.
func MatchKey(input []byte, expected string) (matched bool, consumed int) {
	if len(input) == 0 {
		return false, 0
	}

	expected = strings.ToLower(strings.TrimSpace(expected))
	expected = strings.TrimPrefix(expected, "<")
	expected = strings.TrimSuffix(expected, ">")
	expected = strings.ReplaceAll(expected, "-", "+")
	expected = strings.ReplaceAll(expected, "ctrk", "ctrl")
	expected = strings.ReplaceAll(expected, "ctl", "ctrl")

	if len(expected) == 1 {
		if input[0] == expected[0] {
			return true, 1
		}
	}

	if strings.HasPrefix(expected, "ctrl+") && len(expected) == 6 {
		char := expected[5]
		if char >= 'a' && char <= 'z' {
			targetByte := char - 'a' + 1
			// 0x0d ('m') is the Carriage Return byte. In a raw terminal the
			// Enter/Return key and Ctrl+M both arrive as 0x0d, so they are
			// indistinguishable. Matching a "ctrl+m" keybinding here would let
			// it shadow the Enter key and break line submission (the wrapper
			// checks keybindings before its enter handler). Reserve this byte
			// for Enter so a "ctrl+m" keybinding can never hijack it.
			if targetByte == 0x0d {
				return false, 0
			}
			if input[0] == targetByte {
				return true, 1
			}
			if matched, consumed := matchKittyCtrl(input, int(char)); matched {
				return matched, consumed
			}
		}
	}

	switch expected {
	case "ctrl+space":
		if input[0] == 0x00 {
			return true, 1
		}
		if k, ok := parseKittySequence(input); ok && (k.codepoint == 32 || k.codepoint == 0) && (k.modifiers&4 != 0 || (k.modifiers-1)&4 != 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "tab":
		if input[0] == 0x09 {
			return true, 1
		}
		if k, ok := parseKittySequence(input); ok && k.codepoint == 9 && (k.modifiers == 1 || (k.modifiers-1)&1 == 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "shift+tab":
		if k, ok := parseModifiedCSI(input); ok && k.final == 'Z' && k.eventType != 3 {
			return true, k.consumed
		}
		if k, ok := parseKittySequence(input); ok && k.codepoint == 9 && (k.modifiers == 2 || (k.modifiers-1)&1 != 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "up":
		if k, ok := parseModifiedCSI(input); ok && k.final == 'A' && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
		if k, ok := parseKittySequence(input); ok && k.codepoint == 57352 && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "down":
		if k, ok := parseModifiedCSI(input); ok && k.final == 'B' && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
		if k, ok := parseKittySequence(input); ok && (k.codepoint == 57353 || k.codepoint == 57319) && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "right":
		if k, ok := parseModifiedCSI(input); ok && k.final == 'C' && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
		if k, ok := parseKittySequence(input); ok && (k.codepoint == 57351 || k.codepoint == 57320) && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "left":
		if k, ok := parseModifiedCSI(input); ok && k.final == 'D' && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
		if k, ok := parseKittySequence(input); ok && (k.codepoint == 57350 || k.codepoint == 57321) && (k.modifiers == 1 || k.modifiers == 0) && k.eventType != 3 {
			return true, k.consumed
		}
	case "enter", "cr", "return":
		if input[0] == 0x0d || input[0] == 0x0a {
			return true, 1
		}
		if k, ok := parseKittySequence(input); ok && (k.codepoint == 13 || k.codepoint == 10) && k.eventType != 3 {
			return true, k.consumed
		}
	}

	return false, 0
}

type kittyKey struct {
	codepoint int
	modifiers int
	eventType int
	consumed  int
}

func parseKittySequence(input []byte) (kittyKey, bool) {
	if len(input) < 4 || input[0] != 0x1b || input[1] != '[' {
		return kittyKey{}, false
	}
	uIdx := -1
	for i := 2; i < len(input); i++ {
		if input[i] == 'u' {
			uIdx = i
			break
		}
		if input[i] < 0x20 || input[i] > 0x3f {
			return kittyKey{}, false
		}
	}
	if uIdx == -1 {
		return kittyKey{}, false
	}

	body := string(input[2:uIdx])
	parts := strings.Split(body, ";")
	if len(parts) == 0 {
		return kittyKey{}, false
	}

	codeSub := strings.Split(parts[0], ":")
	codepoint, err := strconv.Atoi(codeSub[0])
	if err != nil {
		return kittyKey{}, false
	}

	modifiers := 1
	eventType := 1

	if len(codeSub) >= 2 {
		if ev, err := strconv.Atoi(codeSub[1]); err == nil {
			eventType = ev
		}
	}

	if len(parts) >= 2 {
		modSub := strings.Split(parts[1], ":")
		if m, err := strconv.Atoi(modSub[0]); err == nil {
			modifiers = m
		}
		if len(modSub) >= 2 {
			if ev, err := strconv.Atoi(modSub[1]); err == nil {
				eventType = ev
			}
		}
	}

	return kittyKey{
		codepoint: codepoint,
		modifiers: modifiers,
		eventType: eventType,
		consumed:  uIdx + 1,
	}, true
}

type csiKey struct {
	final     byte
	modifiers int
	eventType int
	consumed  int
}

func parseModifiedCSI(input []byte) (csiKey, bool) {
	if len(input) < 3 || input[0] != 0x1b {
		return csiKey{}, false
	}
	if input[1] == 'O' {
		switch input[2] {
		case 'A', 'B', 'C', 'D':
			return csiKey{final: input[2], modifiers: 1, eventType: 1, consumed: 3}, true
		}
		return csiKey{}, false
	}
	if input[1] != '[' {
		return csiKey{}, false
	}

	if len(input) >= 3 {
		switch input[2] {
		case 'A', 'B', 'C', 'D', 'Z':
			return csiKey{final: input[2], modifiers: 1, eventType: 1, consumed: 3}, true
		}
	}

	finalIdx := -1
	for i := 2; i < len(input); i++ {
		b := input[i]
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '~' {
			finalIdx = i
			break
		}
		if b < 0x20 || b > 0x3f {
			return csiKey{}, false
		}
	}
	if finalIdx == -1 {
		return csiKey{}, false
	}

	finalByte := input[finalIdx]
	body := string(input[2:finalIdx])
	parts := strings.Split(body, ";")

	modifiers := 1
	eventType := 1

	if len(parts) > 0 && parts[0] != "" {
		sub := strings.Split(parts[0], ":")
		if len(sub) >= 2 {
			if ev, err := strconv.Atoi(sub[1]); err == nil {
				eventType = ev
			}
		}
	}

	if len(parts) >= 2 {
		sub := strings.Split(parts[1], ":")
		if m, err := strconv.Atoi(sub[0]); err == nil {
			modifiers = m
		}
		if len(sub) >= 2 {
			if ev, err := strconv.Atoi(sub[1]); err == nil {
				eventType = ev
			}
		}
	}

	return csiKey{
		final:     finalByte,
		modifiers: modifiers,
		eventType: eventType,
		consumed:  finalIdx + 1,
	}, true
}

func matchKittyCtrl(input []byte, expectedASCII int) (bool, int) {
	k, ok := parseKittySequence(input)
	if !ok || k.eventType == 3 {
		return false, 0
	}
	if (k.modifiers&4 != 0 || (k.modifiers-1)&4 != 0) && k.codepoint == expectedASCII {
		return true, k.consumed
	}
	return false, 0
}

func IsKeyRelease(input []byte) (bool, int) {
	if len(input) == 0 || input[0] != 0x1b {
		return false, 0
	}
	if k, ok := parseKittySequence(input); ok && k.eventType == 3 {
		return true, k.consumed
	}
	if k, ok := parseModifiedCSI(input); ok && k.eventType == 3 {
		return true, k.consumed
	}
	return false, 0
}

// FormatKeyName takes a config key string like "ctrl+r" and formats it for UI display, e.g. "<Ctrl+R>".
func FormatKeyName(key string) string {
	parts := strings.Split(key, "+")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		}
	}
	return "<" + strings.Join(parts, "+") + ">"
}
