package config

import (
	"strconv"
	"strings"
)

// MatchKey matches a sequence of terminal bytes against a configured key string.
// Key strings can be e.g. "ctrl+r", "tab", "shift+tab", "right".
// Returns matched=true if the byte sequence matches the configured key, and consumed=number of bytes to consume.
// splitModArrow splits a binding like "ctrl+up" into its modifier and arrow
// direction. Returns ok=false unless the string is "<modifier>+<arrow>".
func splitModArrow(s string) (mod, dir string, ok bool) {
	i := strings.IndexByte(s, '+')
	if i <= 0 {
		return "", "", false
	}
	dir = s[i+1:]
	switch dir {
	case "up", "down", "left", "right":
		return s[:i], dir, true
	}
	return "", "", false
}

// arrowModifier extracts the modifier parameter from a parameterized CSI arrow
// sequence (ESC [ <col> ; <mod> <A-D>) or a kitty CSI-u sequence (ESC [ <key> ;
// <mod> u). It returns 0 for plain arrows with no modifier parameter.
func arrowModifier(input []byte) int {
	if len(input) < 4 || input[1] != '[' {
		return 0
	}
	lastSemi := -1
	end := len(input)
	for j := 2; j < len(input); j++ {
		c := input[j]
		if c == 'A' || c == 'B' || c == 'C' || c == 'D' || c == 'u' {
			end = j
			break
		}
		if c == ';' {
			lastSemi = j
		}
	}
	if lastSemi == -1 {
		return 0
	}
	mod, err := strconv.Atoi(strings.Trim(string(input[lastSemi+1:end]), " "))
	if err != nil {
		return 0
	}
	return mod
}

// arrowIsCtrl reports whether an arrow sequence carries the Ctrl modifier,
// using the common arrow-modifier encoding where the Ctrl bit has value 4
// (so mod&4 != 0). Covers xterm parameterized arrows and kitty CSI-u.
func arrowIsCtrl(input []byte) bool {
	return arrowModifier(input)&4 != 0
}

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

	// Modifier + arrow-key bindings, e.g. "ctrl+up"/"ctrl+down"/"ctrl+left"/"ctrl+right".
	// Plain "up"/"down"/"left"/"right" bindings (in the switch below) match an arrow
	// of that direction regardless of modifier; a modifier-prefixed binding only
	// matches when the arrow carries that modifier.
	if mod, dir, ok := splitModArrow(expected); ok {
		m, c, adir := MatchArrowKey(input)
		if !m || adir != dir {
			return false, 0
		}
		switch mod {
		case "ctrl":
			if arrowIsCtrl(input) {
				return true, c
			}
		}
		return false, 0
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
		if m, c, d := MatchArrowKey(input); m && d == "up" {
			return m, c
		}
	case "down":
		if m, c, d := MatchArrowKey(input); m && d == "down" {
			return m, c
		}
	case "right":
		if m, c, d := MatchArrowKey(input); m && d == "right" {
			return m, c
		}
	case "left":
		if m, c, d := MatchArrowKey(input); m && d == "left" {
			return m, c
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

// MatchArrowKey matches terminal escape sequences for arrow keys.
// Supports three formats:
//   - Standard CSI: ESC [ A/B/C/D or ESC O A/B/C/D
//   - Parameterized CSI: ESC [ <params> A/B/C/D (e.g., ESC [ 1;3A)
//   - CSI u protocol: ESC [ <keycode> ; <modifiers> u (e.g., ESC [ 107 ; 133 u)
//
// Returns matched=true if the byte sequence is an arrow key, and consumed=number of bytes to consume.
func MatchArrowKey(input []byte) (matched bool, consumed int, direction string) {
	if len(input) < 3 || input[0] != 0x1b {
		return false, 0, ""
	}

	// Standard and parameterized CSI / SS3 arrow keys: ESC [ ... A/B/C/D or ESC O ... A/B/C/D
	if input[1] == '[' || input[1] == 'O' {
		switch input[1] {
		case '[':
			// Find the final terminating byte (A/B/C/D).
			// CSI parameter bytes are 0x30-0x3f (digits, ;, etc.),
			// intermediate bytes are 0x20-0x2f, final byte is 0x40-0x7e.
			for i := 2; i < len(input); i++ {
				c := input[i]
				if c == 'A' || c == 'B' || c == 'C' || c == 'D' {
					dir := ""
					switch c {
					case 'A':
						dir = "up"
					case 'B':
						dir = "down"
					case 'C':
						dir = "right"
					case 'D':
						dir = "left"
					}
					return true, i + 1, dir
				}
				if c < 0x20 || c > 0x7e {
					break
				}
				if c >= 0x40 && c <= 0x7e && c != 'A' && c != 'B' && c != 'C' && c != 'D' {
					break
				}
			}
		case 'O':
			// SS3 form: ESC O A/B/C/D
			if len(input) >= 3 {
				switch input[2] {
				case 'A':
					return true, 3, "up"
				case 'B':
					return true, 3, "down"
				case 'C':
					return true, 3, "right"
				case 'D':
					return true, 3, "left"
				}
			}
		}
	}

	// CSI u protocol (modify keyboard protocol): ESC [ <keycode> ; <modifiers> u
	// Arrow key keycodes: Up=107, Down=108, Right=106, Left=105 (when modifier bits are set)
	// Also: Up=4134?, Down=4133?, etc. — depends on terminal. Common mapping:
	// modifier=2 (Shift): Up=57357, etc. For simplicity, handle both 3-digit and 5-digit keycodes.
	if len(input) >= 6 && input[1] == '[' {
		// Find the 'u' terminator
		uIdx := -1
		for i := 2; i < len(input); i++ {
			if input[i] == 'u' {
				uIdx = i
				break
			}
			if (input[i] < '0' || input[i] > '9') && input[i] != ';' {
				break
			}
		}
		if uIdx > 0 {
			body := string(input[2:uIdx])
			parts := strings.Split(body, ";")
			if len(parts) == 2 {
				keycode, err := strconv.Atoi(parts[0])
				if err != nil {
					return false, 0, ""
				}
				mod, _ := strconv.Atoi(parts[1])
				_ = mod // modifiers don't change which arrow key
				switch keycode {
				case 105, 4133: // left
					return true, uIdx + 1, "left"
				case 106, 4134: // right (fish sends 1;129 but that's parameterized CSI, not CSI u)
					return true, uIdx + 1, "right"
				case 107, 4135: // up
					return true, uIdx + 1, "up"
				case 108, 4136: // down
					return true, uIdx + 1, "down"
				}
			}
		}
	}

	return false, 0, ""
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
