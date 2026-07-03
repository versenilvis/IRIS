package integration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/logger"
	"golang.org/x/term"
)

const (
	boxWidth = 76 // total visual width, corners included
	maxItems = 6
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
						col = getParam(0, 1) - 1
						if col < 0 {
							col = 0
						}
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

type Theme struct {
	Border     lipgloss.Color
	Accent     lipgloss.Color
	Muted      lipgloss.Color
	Text       lipgloss.Color
	TextSel    lipgloss.Color
	Match      lipgloss.Color
	Desc       lipgloss.Color
	DescSel    lipgloss.Color
	SelBg      lipgloss.Color
	ScrollInfo lipgloss.Color
}

var currentTheme = Theme{
	Border:     lipgloss.Color("#a277ff"),
	Accent:     lipgloss.Color("#61ffca"),
	Muted:      lipgloss.Color("#6d6a7f"),
	Text:       lipgloss.Color("#edecee"),
	TextSel:    lipgloss.Color("#ffffff"),
	Match:      lipgloss.Color("#61ffca"),
	Desc:       lipgloss.Color("#9692a8"),
	DescSel:    lipgloss.Color("#edecee"),
	SelBg:      lipgloss.Color("#3d375e"),
	ScrollInfo: lipgloss.Color("#a277ff"),
}

func SetTheme(t Theme) {
	currentTheme = t
}

type Overlay struct {
	mu            sync.Mutex
	Visible       bool
	Items         []core.Suggestion
	Cursor        int
	StartIdx      int
	LastGhostLen  int
	TypedQuery    string
	UserNavigated bool
	PromptLen     int
}

func (o *Overlay) SetPromptLen(l int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.PromptLen != l {
		logger.Debugf("SetPromptLen: %d -> %d", o.PromptLen, l)
		o.PromptLen = l
	}
}

func NewOverlay() *Overlay {
	return &Overlay{Visible: false, Cursor: 0, StartIdx: 0}
}

func (o *Overlay) UpdateItems(items []core.Suggestion) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Items = items
	o.Visible = len(o.Items) > 0
	o.Cursor = 0
	o.StartIdx = 0
}

func fixedWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) > width {
		if width == 1 {
			return "…"
		}
		return string(runes[:width-1]) + "…"
	}
	if len(runes) < width {
		return s + strings.Repeat(" ", width-len(runes))
	}
	return s
}

func (o *Overlay) RenderGhostText(buffer string, userNavigated bool) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 {
		return ""
	}

	var s strings.Builder
	ghostText := ""
	if !userNavigated && buffer != "" {
		topCmd := o.Items[0].Cmd
		if strings.HasPrefix(strings.ToLower(topCmd), strings.ToLower(buffer)) {
			ghostText = topCmd[len(buffer):]
		}
	}

	padLen := o.LastGhostLen - len(ghostText)
	if padLen < 0 {
		padLen = 0
	}
	padLen += 10

	if ghostText != "" || padLen > 0 {
		s.WriteString("\0337")
		if ghostText != "" {
			s.WriteString("\033[90m")
			s.WriteString(ghostText)
			s.WriteString("\033[0m")
		}
		if padLen > 0 {
			s.WriteString(strings.Repeat(" ", padLen))
		}
		s.WriteString("\0338")
		o.LastGhostLen = len(ghostText)
	}

	return s.String()
}

func renderMatchedTitle(title, typed string, selected bool, w int) string {
	t := currentTheme
	textColor := t.Text
	if selected {
		textColor = t.TextSel
	}

	base := lipgloss.NewStyle().Foreground(textColor)
	match := lipgloss.NewStyle().Foreground(t.Match).Bold(true)
	if selected {
		base = base.Background(t.SelBg)
		match = match.Background(t.SelBg)
	}

	display := fixedWidth(title, w)

	if typed == "" || !strings.HasPrefix(strings.ToLower(display), strings.ToLower(typed)) {
		return base.Render(display)
	}

	matchLen := len(typed)
	if matchLen > len(display) {
		matchLen = len(display)
	}
	return match.Render(display[:matchLen]) + base.Render(display[matchLen:])
}

func (o *Overlay) Render() string {
	return o.draw()
}

func (o *Overlay) draw() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 {
		return ""
	}

	t := currentTheme
	border := lipgloss.NewStyle().Foreground(t.Border)
	scrollStyle := lipgloss.NewStyle().Foreground(t.ScrollInfo)

	var s strings.Builder
	s.WriteString("\033[?7l")

	typedLen := len([]rune(o.TypedQuery))
	targetCol := o.PromptLen + typedLen

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 120
	}
	if targetCol+boxWidth > width {
		targetCol = width - boxWidth
	}
	if targetCol < 0 {
		targetCol = 0
	}
	logger.Debugf("Overlay draw: pLen=%d, typedLen=%d, targetCol=%d, width=%d", o.PromptLen, typedLen, targetCol, width)

	s.WriteString("\0337")

	windowSize := maxItems
	if len(o.Items) < windowSize {
		windowSize = len(o.Items)
	}

	scrolloffUp := 1
	if windowSize <= 3 {
		scrolloffUp = 0
	}

	if o.Cursor < o.StartIdx+scrolloffUp {
		o.StartIdx = o.Cursor - scrolloffUp
	}
	if o.Cursor >= o.StartIdx+windowSize {
		o.StartIdx = o.Cursor - windowSize + 1
	}
	if o.StartIdx < 0 {
		o.StartIdx = 0
	}
	if o.StartIdx > len(o.Items)-windowSize {
		o.StartIdx = len(o.Items) - windowSize
	}
	if o.StartIdx < 0 {
		o.StartIdx = 0
	}

	start := o.StartIdx
	end := start + windowSize
	totalLines := windowSize + 2

	for range totalLines {
		s.WriteByte('\n')
	}
	fmt.Fprintf(&s, "\033[%dA", totalLines)

	s.WriteString("\0337")

	moveToTarget := func() {
		s.WriteString("\r")
		if targetCol > 0 {
			fmt.Fprintf(&s, "\033[%dC", targetCol)
		}
	}

	inner := boxWidth - 2 // width between the two border pipes/corners

	// top side border with scroll counter
	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", 1)
	s.WriteString("\033[2K")
	moveToTarget()

	scrollInfo := ""
	if len(o.Items) > windowSize {
		scrollInfo = fmt.Sprintf(" %d/%d ", o.Cursor+1, len(o.Items))
	}
	leftDash := 3
	rightDash := inner - leftDash - len(scrollInfo)
	if scrollInfo == "" {
		leftDash = 0
		rightDash = inner
	}
	fmt.Fprintf(&s, "%s%s%s%s%s",
		border.Render("╭"),
		border.Render(strings.Repeat("─", leftDash)),
		scrollStyle.Render(scrollInfo),
		border.Render(strings.Repeat("─", rightDash)),
		border.Render("╮"),
	)

	// left and right side border with item rows
	descW := 22
	padGap := 2
	markerW := 1
	sidePad := 1
	titleW := inner - sidePad*2 - markerW - 1 - padGap - descW

	for i := start; i < end; i++ {
		s.WriteString("\0338")
		fmt.Fprintf(&s, "\033[%dB", (i-start)+2)
		s.WriteString("\033[2K")
		moveToTarget()

		it := o.Items[i]
		selected := i == o.Cursor

		left := border.Render("│")
		right := border.Render("│")

		bg := lipgloss.NewStyle()
		if selected {
			bg = bg.Background(t.SelBg)
		}

		marker := " "
		markerStyle := bg.Foreground(t.Muted)
		if selected {
			marker = "▶"
			markerStyle = bg.Foreground(t.Accent).Bold(true)
		}

		title := renderMatchedTitle(it.Cmd, o.TypedQuery, selected, titleW)

		descColor := t.Desc
		if selected {
			descColor = t.DescSel
		}
		desc := bg.Foreground(descColor).Render(fixedWidth(it.Desc, descW))

		fmt.Fprintf(&s, "%s%s%s%s%s%s%s%s%s",
			left,
			bg.Render(" "),
			markerStyle.Render(marker),
			bg.Render(" "),
			title,
			bg.Render(strings.Repeat(" ", padGap)),
			desc,
			bg.Render(" "),
			right,
		)
	}

	// bottom side border with footer shortcut hints
	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", windowSize+2)
	s.WriteString("\033[2K")
	moveToTarget()
	footerInfo := " <Tab> Accept • <Ctrl+R> Mode "
	footerRunes := len([]rune(footerInfo))
	rightDash = 2
	leftDash = inner - footerRunes - rightDash
	if leftDash < 0 {
		leftDash = 0
	}
	fmt.Fprintf(&s, "%s%s%s%s%s",
		border.Render("╰"),
		border.Render(strings.Repeat("─", leftDash)),
		lipgloss.NewStyle().Foreground(t.ScrollInfo).Render(footerInfo),
		border.Render(strings.Repeat("─", rightDash)),
		border.Render("╯"),
	)

	s.WriteString("\0338")
	s.WriteString("\033[?7h")
	return s.String()
}

func (o *Overlay) Clear() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	var s strings.Builder
	s.WriteString("\033[?7l")
	s.WriteString("\0337")

	for i := 0; i < maxItems+2; i++ {
		s.WriteString("\0338")
		fmt.Fprintf(&s, "\033[%dB", i+1)
		s.WriteString("\r\033[2K")
	}

	s.WriteString("\0338")
	s.WriteString("\033[?7h")
	return s.String()
}

func (o *Overlay) ClearAndDisable() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible && len(o.Items) == 0 && o.LastGhostLen == 0 {
		return ""
	}

	o.Visible = false
	o.Items = nil
	o.TypedQuery = ""
	o.UserNavigated = false

	var s strings.Builder
	s.WriteString("\033[?7l")

	if o.LastGhostLen > 0 {
		s.WriteString("\0337")
		s.WriteString(strings.Repeat(" ", o.LastGhostLen+10))
		s.WriteString("\0338")
		o.LastGhostLen = 0
	}

	s.WriteString("\0337")

	for i := 0; i < maxItems+2; i++ {
		s.WriteString("\0338")
		fmt.Fprintf(&s, "\033[%dB", i+1)
		s.WriteString("\r\033[2K")
	}

	s.WriteString("\0338")
	s.WriteString("\033[?7h")
	return s.String()
}
