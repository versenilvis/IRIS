package integration

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func ComputeCursorCol(data []byte) int {
	col := 0
	i := 0
	n := len(data)
	for i < n {
		b := data[i]
		if b == '\r' {
			col = 0
			i++
			continue
		}
		if b == '\b' || b == 0x7f {
			col--
			if col < 0 {
				col = 0
			}
			i++
			continue
		}
		if b == '\t' {
			col = (col + 8) &^ 7
			i++
			continue
		}
		if b == '\033' {
			if i+1 < n && data[i+1] == '[' {
				j := i + 2
				for j < n && data[j] >= 0x20 && data[j] <= 0x3F {
					j++
				}
				if j < n {
					cmd := data[j]
					paramsStr := string(data[i+2 : j])
					paramsStr = strings.TrimLeft(paramsStr, "?>=")
					parts := strings.Split(paramsStr, ";")
					getParam := func(idx, def int) int {
						if idx < len(parts) && parts[idx] != "" {
							if v, err := strconv.Atoi(parts[idx]); err == nil && v > 0 {
								return v
							}
						}
						return def
					}
					switch cmd {
					case 'C':
						col += getParam(0, 1)
					case 'D':
						col -= getParam(0, 1)
						if col < 0 {
							col = 0
						}
					case 'G':
						col = max(getParam(0, 1)-1, 0)
					}
					i = j + 1
					continue
				}
				break
			} else if i+1 < n && data[i+1] == ']' {
				j := i + 2
				for j < n {
					if data[j] == '\007' {
						j++
						break
					}
					if data[j] == '\033' && j+1 < n && data[j+1] == '\\' {
						j += 2
						break
					}
					j++
				}
				i = j
				continue
			} else if i+1 < n && (data[i+1] == 'P' || data[i+1] == 'X' || data[i+1] == '^' || data[i+1] == '_') {
				j := i + 2
				for j < n {
					if data[j] == '\033' && j+1 < n && data[j+1] == '\\' {
						j += 2
						break
					}
					j++
				}
				i = j
				continue
			} else if i+1 < n {
				i += 2
				continue
			} else {
				break
			}
		}
		if b < 0x20 {
			i++
			continue
		}
		if b < 0x7f {
			col++
			i++
			continue
		}
		r, size := utf8.DecodeRune(data[i:])
		w := lipgloss.Width(string(r))
		col += w
		i += size
	}
	return col
}
