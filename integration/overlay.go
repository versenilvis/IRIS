package integration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/spec"
	"github.com/versenilvis/iris/config"
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
	Items         []spec.Suggestion
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

func (o *Overlay) UpdateItems(items []spec.Suggestion) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Items = items
	o.Visible = len(o.Items) > 0
	o.Cursor = 0
	o.StartIdx = 0
}

func (o *Overlay) IsVisible() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.Visible
}

func (o *Overlay) GetUserNavigated() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.UserNavigated
}

func (o *Overlay) SetUserNavigated(v bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.UserNavigated = v
}

func (o *Overlay) GetTypedQuery() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.TypedQuery
}

func (o *Overlay) GetCurrentCmd() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.Items) > 0 && o.Cursor >= 0 && o.Cursor < len(o.Items) {
		return o.Items[o.Cursor].Cmd
	}
	return ""
}

func (o *Overlay) GetTopCmd() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.Items) > 0 {
		return o.Items[0].Cmd
	}
	return ""
}

func (o *Overlay) Show() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.UserNavigated = false
	o.Visible = true
}

func (o *Overlay) ResetCursor() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Cursor = 0
}

func (o *Overlay) SetQueryAndItems(query string, items []spec.Suggestion) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.TypedQuery = query
	o.UserNavigated = false
	o.Items = items
	o.Visible = len(o.Items) > 0
	o.Cursor = 0
	o.StartIdx = 0
}

func (o *Overlay) ClearGhostLen() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	l := o.LastGhostLen
	o.LastGhostLen = 0
	return l
}

func (o *Overlay) MoveCursor(dir string) (moved bool, selectedCmd string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.Visible || len(o.Items) == 0 {
		return false, ""
	}
	o.UserNavigated = true
	oldCursor := o.Cursor
	if dir == "up" {
		o.Cursor--
		if o.Cursor < 0 {
			o.Cursor = 0
		}
	} else {
		o.Cursor++
		if o.Cursor >= len(o.Items) {
			o.Cursor = len(o.Items) - 1
		}
	}
	if o.Cursor == oldCursor {
		return false, ""
	}
	return true, o.Items[o.Cursor].Cmd
}

func (o *Overlay) SetHistoryList(items []spec.Suggestion, startAtBottom bool) string {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.TypedQuery = ""
	o.UserNavigated = true
	o.Items = items
	o.Visible = len(o.Items) > 0
	if startAtBottom && len(o.Items) > 0 {
		o.Cursor = len(o.Items) - 1
	} else {
		o.Cursor = 0
	}
	o.StartIdx = 0
	if len(o.Items) > 0 && o.Cursor >= 0 && o.Cursor < len(o.Items) {
		return o.Items[o.Cursor].Cmd
	}
	return ""
}

func fixedWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	visualWidth := lipgloss.Width(s)
	if visualWidth == width {
		return s
	}
	if visualWidth < width {
		return s + strings.Repeat(" ", width-visualWidth)
	}
	var sb strings.Builder
	currentWidth := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if currentWidth+rw > width-1 {
			break
		}
		sb.WriteRune(r)
		currentWidth += rw
	}
	sb.WriteString("…")
	rem := width - lipgloss.Width(sb.String())
	if rem > 0 {
		sb.WriteString(strings.Repeat(" ", rem))
	}
	return sb.String()
}

func truncateToWidth(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxW {
		return s
	}
	if maxW == 1 {
		return "…"
	}
	runes := []rune(s)
	var sb strings.Builder
	w := 0
	for _, r := range runes {
		rw := lipgloss.Width(string(r))
		if w+rw > maxW-1 { // leave 1 column for '…'
			break
		}
		sb.WriteRune(r)
		w += rw
	}
	sb.WriteRune('…')
	return sb.String()
}

func (o *Overlay) GetGhostText(buffer string, cursorAtEnd bool) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 || !cursorAtEnd || buffer == "" {
		return ""
	}

	var topCmd string
	if o.Cursor >= 0 && o.Cursor < len(o.Items) {
		topCmd = o.Items[o.Cursor].Cmd
	} else {
		topCmd = o.Items[0].Cmd
	}

	if strings.HasPrefix(strings.ToLower(topCmd), strings.ToLower(buffer)) {
		return topCmd[len(buffer):]
	}
	return ""
}

func (o *Overlay) RenderGhostText(buffer string, userNavigated bool, cursorAtEnd bool) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 {
		if o.LastGhostLen > 0 {
			padLen := o.LastGhostLen + 4
			o.LastGhostLen = 0
			return "\0337" + strings.Repeat(" ", padLen) + "\0338"
		}
		return ""
	}

	var s strings.Builder
	ghostText := ""
	if cursorAtEnd && buffer != "" {
		var topCmd string
		if o.Cursor >= 0 && o.Cursor < len(o.Items) {
			topCmd = o.Items[o.Cursor].Cmd
		} else {
			topCmd = o.Items[0].Cmd
		}
		if strings.HasPrefix(strings.ToLower(topCmd), strings.ToLower(buffer)) {
			ghostText = topCmd[len(buffer):]
		}
	}

	if ghostText != "" {
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || width <= 0 {
			width = 120
		}
		cursorCol := o.PromptLen + lipgloss.Width(buffer)
		availableCols := width - cursorCol
		if availableCols <= 0 {
			ghostText = ""
		} else if lipgloss.Width(ghostText) > availableCols {
			ghostText = truncateToWidth(ghostText, availableCols)
		}
	}

	if ghostText == "" && o.LastGhostLen == 0 {
		return ""
	}

	ghostWidth := lipgloss.Width(ghostText)
	padLen := max(o.LastGhostLen-ghostWidth, 0)
	if o.LastGhostLen > 0 {
		padLen += 4
	}

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
	o.LastGhostLen = ghostWidth

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

	typedRunes := []rune(typed)
	displayRunes := []rune(display)
	matchLen := min(len(typedRunes), len(displayRunes))
	return match.Render(string(displayRunes[:matchLen])) + base.Render(string(displayRunes[matchLen:]))
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

	windowSize := min(len(o.Items), maxItems)

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

	style := strings.ToLower(config.Get().UI.Style)
	isClassic := style == "classic" || style == "minimal" || style == "minimalist"

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
	if isClassic && scrollInfo != "" {
		leftDash = (inner - len(scrollInfo)) / 2
	}
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
	descW := 24
	padGap := 2
	markerW := 1
	iconW := 2
	if isClassic || !config.Get().UI.NerdFonts {
		iconW = 0
	}
	sidePad := 1
	titleW := inner - sidePad*2 - markerW - 1 - iconW
	if iconW > 0 {
		titleW--
	}
	titleW = titleW - padGap - descW

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

		iconGlyph := lookupIcon(it.Icon)
		iconColor := t.Muted
		if selected {
			iconColor = t.Accent
		}
		iconStr := bg.Foreground(iconColor).Render(fixedWidth(iconGlyph, iconW))

		title := renderMatchedTitle(it.Cmd, o.TypedQuery, selected, titleW)

		descColor := t.Desc
		if selected {
			descColor = t.DescSel
		}

		var desc string
		if isClassic {
			if it.Icon == "alias" {
				desc = bg.Foreground(descColor).Render(fixedWidth("alias: "+it.Desc, descW))
			} else {
				desc = bg.Foreground(descColor).Render(fixedWidth(it.Desc, descW))
			}
		} else {
			switch it.Icon {
			case "alias":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color("#2a2342")).Foreground(lipgloss.Color("#a277ff"))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color("#a277ff")).Foreground(lipgloss.Color("#110f18")).Bold(true)
				}
				tag := boxStyle.Render(" alias ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw-1, 0)
				desc = tag + bg.Render(" ") + bg.Foreground(descColor).Render(fixedWidth(it.Desc, rem))
			case "history":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color("#1a2d36")).Foreground(lipgloss.Color("#61ffca"))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color("#61ffca")).Foreground(lipgloss.Color("#110f18")).Bold(true)
				}
				tag := boxStyle.Render(" history ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw, 0)
				desc = tag + bg.Render(strings.Repeat(" ", rem))
			case "system":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color("#1e1d28")).Foreground(lipgloss.Color("#a277ff"))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color("#a277ff")).Foreground(lipgloss.Color("#110f18")).Bold(true)
				}
				tag := boxStyle.Render(" system ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw, 0)
				desc = tag + bg.Render(strings.Repeat(" ", rem))
			default:
				desc = bg.Foreground(descColor).Render(fixedWidth(it.Desc, descW))
			}
		}

		iconSection := ""
		if iconW > 0 {
			iconSection = iconStr + bg.Render(" ")
		}

		fmt.Fprintf(&s, "%s%s%s%s%s%s%s%s%s%s",
			left,
			bg.Render(" "),
			markerStyle.Render(marker),
			bg.Render(" "),
			iconSection,
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

	footerInfo := ""
	if !isClassic {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a277ff")).Bold(true)
		tabKey := keyStyle.Render("<Tab>")
		ctrlRKey := keyStyle.Render("<Ctrl+R>")
		acceptText := lipgloss.NewStyle().Foreground(t.ScrollInfo).Render(" Accept")
		modeText := lipgloss.NewStyle().Foreground(t.ScrollInfo).Render(" Mode")
		footerInfo = fmt.Sprintf(" %s%s • %s%s ", tabKey, acceptText, ctrlRKey, modeText)
	}

	footerRunes := lipgloss.Width(footerInfo)
	rightDash = 2
	leftDash = inner - footerRunes - rightDash
	if footerInfo == "" {
		leftDash = 0
		rightDash = inner
	}
	if leftDash < 0 {
		leftDash = 0
	}
	fmt.Fprintf(&s, "%s%s%s%s%s",
		border.Render("╰"),
		border.Render(strings.Repeat("─", leftDash)),
		footerInfo,
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

	for i := range maxItems + 2 {
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
	o.Cursor = 0
	o.StartIdx = 0

	var s strings.Builder
	s.WriteString("\033[?7l")

	if o.LastGhostLen > 0 {
		s.WriteString("\0337")
		s.WriteString(strings.Repeat(" ", o.LastGhostLen+10))
		s.WriteString("\0338")
		o.LastGhostLen = 0
	}

	s.WriteString("\0337")

	for i := range maxItems + 2 {
		s.WriteString("\0338")
		fmt.Fprintf(&s, "\033[%dB", i+1)
		s.WriteString("\r\033[2K")
	}

	s.WriteString("\0338")
	s.WriteString("\033[?7h")
	return s.String()
}
