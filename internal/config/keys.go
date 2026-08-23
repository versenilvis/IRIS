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
		// Some terminals (foot, kitty, etc.) send kitty keyboard protocol
		// escape sequences for Ctrl+letter instead of raw control bytes:
		// CSI <keycode> ; <modifiers> <action>
		// where modifiers include Ctrl=4 and action='u' means press.
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

// matchKittyCtrl matches the kitty keyboard protocol CSI sequence for Ctrl+<letter>.
// Format: ESC [ <keycode> ; <modifiers> <action>
// For Ctrl+<letter>: keycode is the ASCII code of the letter, modifiers must have bit 2 (value 4) set, action is 'u'.
func matchKittyCtrl(input []byte, expectedASCII int) (matched bool, consumed int) {
	if len(input) < 6 || input[0] != 0x1b || input[1] != '[' {
		return false, 0
	}

	uIdx := -1
	for i := 2; i < len(input); i++ {
		if input[i] == 'u' {
			uIdx = i
			break
		}
	}
	if uIdx == -1 {
		return false, 0
	}

	body := string(input[2:uIdx])
	parts := strings.Split(body, ";")
	if len(parts) != 2 {
		return false, 0
	}

	keycode, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, 0
	}
	modifiers, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, 0
	}

	if modifiers&4 == 0 {
		return false, 0
	}

	if keycode != expectedASCII {
		return false, 0
	}

	return true, uIdx + 1
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
