package config

import (
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

	if len(expected) == 1 {
		if input[0] == expected[0] {
			return true, 1
		}
	}

	if strings.HasPrefix(expected, "ctrl+") && len(expected) == 6 {
		char := expected[5]
		if char >= 'a' && char <= 'z' {
			targetByte := char - 'a' + 1
			if input[0] == targetByte {
				return true, 1
			}
		}
	}

	switch expected {
	case "ctrl+space":
		if input[0] == 0x00 {
			return true, 1
		}
	case "tab":
		if input[0] == 0x09 {
			return true, 1
		}
	case "shift+tab":
		// typically \033[Z
		if len(input) >= 3 && input[0] == 0x1b && input[1] == '[' && input[2] == 'Z' {
			return true, 3
		}
	case "up":
		if len(input) >= 3 && input[0] == 0x1b && (input[1] == '[' || input[1] == 'O') && input[2] == 'A' {
			return true, 3
		}
	case "down":
		if len(input) >= 3 && input[0] == 0x1b && (input[1] == '[' || input[1] == 'O') && input[2] == 'B' {
			return true, 3
		}
	case "right":
		// typically \033[C or \033OC
		if len(input) >= 3 && input[0] == 0x1b && (input[1] == '[' || input[1] == 'O') && input[2] == 'C' {
			return true, 3
		}
	case "left":
		if len(input) >= 3 && input[0] == 0x1b && (input[1] == '[' || input[1] == 'O') && input[2] == 'D' {
			return true, 3
		}
	case "enter", "cr", "return":
		if input[0] == 0x0d || input[0] == 0x0a {
			return true, 1
		}
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
