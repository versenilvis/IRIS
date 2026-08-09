package integration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/logger"
	"github.com/versenilvis/iris/spec"
	"golang.org/x/term"
)

const (
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

func (o *Overlay) InjectAISuggestion(sugg spec.Suggestion) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.TypedQuery == "" {
		return false
	}
	if o.UserNavigated {
		return false
	}

	var currentConf int
	if len(o.Items) > 0 {
		currentConf = o.Items[0].Confidence
		if currentConf == 0 {
			if o.Items[0].Source == "history" {
				currentConf = 70
			} else {
				currentConf = 50
			}
		}
	}

	if !strings.HasPrefix(strings.ToLower(sugg.Cmd), strings.ToLower(o.TypedQuery)) {
		return false
	}
	if sugg.Confidence <= currentConf && len(o.Items) > 0 {
		return false
	}

	if len(o.Items) == 0 {
		o.Items = []spec.Suggestion{sugg}
	} else if strings.EqualFold(o.Items[0].Cmd, sugg.Cmd) {
		if o.Visible && o.Items[0].Confidence == sugg.Confidence {
			return false
		}
		o.Items[0] = sugg
	} else {
		o.Items = append([]spec.Suggestion{sugg}, o.Items...)
		if len(o.Items) > 100 {
			o.Items = o.Items[:100]
		}
	}
	o.Visible = true
	o.Cursor = 0
	o.StartIdx = 0
	return true
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
	if visualWidth > width {
		s = ansi.Truncate(s, width, "…")
	}
	if rem := width - lipgloss.Width(s); rem > 0 {
		s += strings.Repeat(" ", rem)
	}
	return s
}

func truncateToWidth(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	return ansi.Truncate(s, maxW, "…")
}

// titledEdge renders one horizontal box edge (top or bottom) with styled
// content embedded in the dash run. inner is the width between corners.
// leftPad positions the content that many dashes from the left corner;
// pass -1 to center the content (classic style).
func titledEdge(left, right string, inner int, content string, border lipgloss.Style, leftPad int) string {
	contentWidth := lipgloss.Width(content)
	if content == "" {
		return border.Render(left + strings.Repeat("─", inner) + right)
	}
	leftDash := leftPad
	if leftPad < 0 {
		leftDash = max((inner-contentWidth)/2, 0)
	}
	rightDash := inner - leftDash - contentWidth
	if rightDash < 0 {
		leftDash = max(leftDash+rightDash, 0)
		rightDash = 0
	}
	return border.Render(left+strings.Repeat("─", leftDash)) +
		content +
		border.Render(strings.Repeat("─", rightDash)+right)
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

func (o *Overlay) ClearGhostTextState() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.LastGhostLen = 0
}

func (o *Overlay) HideGhostTextSync() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.LastGhostLen > 0 {
		padLen := o.LastGhostLen + 4
		res := ansi.SaveCursor + strings.Repeat(" ", padLen) + ansi.RestoreCursor
		o.LastGhostLen = 0
		return res
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
			return ansi.SaveCursor + strings.Repeat(" ", padLen) + ansi.RestoreCursor
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
		width := termWidth()
		totalCol := o.PromptLen + lipgloss.Width(buffer)
		cursorCol := totalCol
		if width > 0 {
			cursorCol = totalCol % width
		}
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

	s.WriteString(ansi.SaveCursor)
	if ghostText != "" {
		styled := lipgloss.NewStyle().Foreground(lipgloss.Color(config.Theme().GhostText)).Render(ghostText)
		s.WriteString(styled)
	}
	if padLen > 0 {
		s.WriteString(strings.Repeat(" ", padLen))
	}
	s.WriteString(ansi.RestoreCursor)
	o.LastGhostLen = ghostWidth

	return s.String()
}

// termWidth returns the terminal width, falling back to 120 when the size
// can't be determined (e.g. stdout isn't a TTY in tests).
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 120
	}
	return w
}

func renderMatchedTitle(title, typed string, selected bool, w int) string {
	t := config.Theme()
	textColor := t.Text
	if selected {
		textColor = t.TextSel
	}

	base := lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))
	match := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Match)).Bold(true)
	if selected {
		base = base.Background(lipgloss.Color(t.SelBg))
		match = match.Background(lipgloss.Color(t.SelBg))
	}

	display := fixedWidth(title, w)

	if typed == "" || !strings.HasPrefix(strings.ToLower(display), strings.ToLower(typed)) {
		return base.Render(display)
	}

	// Split by display width, not rune count, so a typed prefix containing
	// wide runes (CJK, emoji) highlights exactly the matching cells.
	matchW := lipgloss.Width(typed)
	highlighted := ansi.Truncate(display, matchW, "")
	rest := strings.TrimPrefix(display, highlighted)
	return match.Render(highlighted) + base.Render(rest)
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

	t := config.Theme()
	border := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Border))
	scrollStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.ScrollInfo))

	var s strings.Builder
	s.WriteString(ansi.ResetModeAutoWrap)

	typedLen := lipgloss.Width(o.TypedQuery)
	width := termWidth()

	// ComputeCursorCol returns the total visual width, not the column on the
	// current line. When the prompt is wider than the terminal the cursor has
	// wrapped, so using PromptLen directly overflows the screen and the box
	// lands at the wrong horizontal position.
	boxWidth := config.Get().UI.MaxWidth
	if boxWidth <= 0 {
		boxWidth = 76 // Default if 0
	}
	if boxWidth > width {
		boxWidth = width // Responsive: don't overflow
	}
	if boxWidth < 40 {
		boxWidth = 40 // Minimum safe width
	}

	// ComputeCursorCol returns the total visual width, not the column on the
	// current line. When the prompt is wider than the terminal the cursor has
	// wrapped, so using PromptLen directly overflows the screen and the box
	// lands at the wrong horizontal position.
	totalCol := o.PromptLen + typedLen
	cursorCol := totalCol
	if width > 0 {
		cursorCol = totalCol % width
	}
	targetCol := cursorCol
	if targetCol+boxWidth > width {
		targetCol = width - boxWidth
	}
	if targetCol < 0 {
		targetCol = 0
	}
	logger.Debugf("Overlay draw: pLen=%d, typedLen=%d, totalCol=%d, cursorCol=%d, targetCol=%d, width=%d", o.PromptLen, typedLen, totalCol, cursorCol, targetCol, width)

	s.WriteString(ansi.SaveCursor)

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

	s.WriteString(strings.Repeat("\n", totalLines))
	s.WriteString(ansi.CursorUp(totalLines))

	s.WriteString(ansi.SaveCursor)

	moveToTarget := func() {
		s.WriteString("\r")
		if targetCol > 0 {
			s.WriteString(ansi.CursorForward(targetCol))
		}
	}

	inner := boxWidth - 2 // width between the two border pipes/corners

	style := strings.ToLower(config.Get().UI.Style)
	isClassic := style == "classic" || style == "minimal" || style == "minimalist"

	// top side border with scroll counter
	s.WriteString(ansi.RestoreCursor)
	s.WriteString(ansi.CursorDown(1))
	s.WriteString(ansi.EraseEntireLine)
	moveToTarget()

	scrollInfo := ""
	if len(o.Items) > windowSize {
		scrollInfo = scrollStyle.Render(fmt.Sprintf(" %d/%d ", o.Cursor+1, len(o.Items)))
	}
	scrollPad := 3
	if isClassic {
		scrollPad = -1 // centered
	}
	s.WriteString(titledEdge("╭", "╮", inner, scrollInfo, border, scrollPad))

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
		s.WriteString(ansi.RestoreCursor)
		s.WriteString(ansi.CursorDown((i - start) + 2))
		s.WriteString(ansi.EraseEntireLine)
		moveToTarget()

		it := o.Items[i]
		selected := i == o.Cursor

		left := border.Render("│")
		right := border.Render("│")

		bg := lipgloss.NewStyle()
		if selected {
			bg = bg.Background(lipgloss.Color(t.SelBg))
		}

		marker := " "
		markerStyle := bg.Foreground(lipgloss.Color(t.Muted))
		if selected {
			marker = "▶"
			markerStyle = bg.Foreground(lipgloss.Color(t.Accent)).Bold(true)
		}

		iconGlyph := lookupIcon(it.Icon)
		iconColor := t.Muted
		if selected {
			iconColor = t.Accent
		}
		iconStr := bg.Foreground(lipgloss.Color(iconColor)).Render(fixedWidth(iconGlyph, iconW))

		title := renderMatchedTitle(it.Cmd, o.TypedQuery, selected, titleW)

		descColor := t.Desc
		if selected {
			descColor = t.DescSel
		}

		var desc string
		if isClassic {
			if it.Icon == "alias" {
				desc = bg.Foreground(lipgloss.Color(descColor)).Render(fixedWidth("alias: "+it.Desc, descW))
			} else {
				desc = bg.Foreground(lipgloss.Color(descColor)).Render(fixedWidth(it.Desc, descW))
			}
		} else {
			switch it.Icon {
			case "alias":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color(t.Alias)).Foreground(lipgloss.Color(t.AliasSel))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color(t.AliasSel)).Foreground(lipgloss.Color(t.SelText)).Bold(true)
				}
				tag := boxStyle.Render(" alias ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw-1, 0)
				desc = tag + bg.Render(" ") + bg.Foreground(lipgloss.Color(descColor)).Render(fixedWidth(it.Desc, rem))
			case "atuin":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color(t.Alias)).Foreground(lipgloss.Color(t.AliasSel))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color(t.AliasSel)).Foreground(lipgloss.Color(t.SelText)).Bold(true)
				}
				tag := boxStyle.Render(" atuin ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw-1, 0)
				desc = tag + bg.Render(" ") + bg.Foreground(lipgloss.Color(descColor)).Render(fixedWidth("atuin history", rem))
			case "history":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color(t.History)).Foreground(lipgloss.Color(t.HistorySel))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color(t.HistorySel)).Foreground(lipgloss.Color(t.SelText)).Bold(true)
				}
				tag := boxStyle.Render(" history ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw, 0)
				desc = tag + bg.Render(strings.Repeat(" ", rem))
			case "system":
				boxStyle := lipgloss.NewStyle().Background(lipgloss.Color(t.Sys)).Foreground(lipgloss.Color(t.SysSel))
				if selected {
					boxStyle = lipgloss.NewStyle().Background(lipgloss.Color(t.SysSel)).Foreground(lipgloss.Color(t.SelText)).Bold(true)
				}
				tag := boxStyle.Render(" system ")
				tw := lipgloss.Width(tag)
				rem := max(descW-tw, 0)
				desc = tag + bg.Render(strings.Repeat(" ", rem))
			default:
				desc = bg.Foreground(lipgloss.Color(descColor)).Render(fixedWidth(it.Desc, descW))
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
	s.WriteString(ansi.RestoreCursor)
	s.WriteString(ansi.CursorDown(windowSize + 2))
	s.WriteString(ansi.EraseEntireLine)
	moveToTarget()

	footerInfo := ""
	if !isClassic {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Key)).Bold(true)
		selectKey := keyStyle.Render(config.FormatKeyName(config.Get().Keybindings.SelectSuggestion))
		ctrlRKey := keyStyle.Render(config.FormatKeyName(config.Get().Keybindings.ToggleMode))
		acceptText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.ScrollInfo)).Render(" Accept")
		modeText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.ScrollInfo)).Render(" Mode")
		footerInfo = fmt.Sprintf(" %s%s • %s%s ", selectKey, acceptText, ctrlRKey, modeText)
	}

	s.WriteString(titledEdge("╰", "╯", inner, footerInfo, border, inner-lipgloss.Width(footerInfo)-2))

	s.WriteString(ansi.RestoreCursor)
	s.WriteString(ansi.SetModeAutoWrap)
	return s.String()
}

// clearLinesBelow erases n lines below the cursor, leaving the cursor where
// it started. Caller is responsible for wrapping/resetting auto-wrap mode.
func clearLinesBelow(s *strings.Builder, n int) {
	s.WriteString(ansi.SaveCursor)
	for i := range n {
		s.WriteString(ansi.RestoreCursor)
		s.WriteString(ansi.CursorDown(i + 1))
		s.WriteString("\r" + ansi.EraseEntireLine)
	}
	s.WriteString(ansi.RestoreCursor)
}

func (o *Overlay) Clear() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	var s strings.Builder
	s.WriteString(ansi.ResetModeAutoWrap)
	clearLinesBelow(&s, maxItems+2)
	s.WriteString(ansi.SetModeAutoWrap)
	return s.String()
}

func (o *Overlay) HideMenu(query string) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.TypedQuery = query
	if !o.Visible && len(o.Items) == 0 && o.LastGhostLen == 0 {
		return ""
	}

	o.Visible = false
	o.Items = nil
	o.UserNavigated = false
	o.Cursor = 0
	o.StartIdx = 0

	var s strings.Builder
	s.WriteString(ansi.ResetModeAutoWrap)

	if o.LastGhostLen > 0 {
		s.WriteString(ansi.SaveCursor)
		s.WriteString(strings.Repeat(" ", o.LastGhostLen+10))
		s.WriteString(ansi.RestoreCursor)
		o.LastGhostLen = 0
	}

	clearLinesBelow(&s, maxItems+2)
	s.WriteString(ansi.SetModeAutoWrap)
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
	s.WriteString(ansi.ResetModeAutoWrap)

	if o.LastGhostLen > 0 {
		s.WriteString(ansi.SaveCursor)
		s.WriteString(strings.Repeat(" ", o.LastGhostLen+10))
		s.WriteString(ansi.RestoreCursor)
		o.LastGhostLen = 0
	}

	clearLinesBelow(&s, maxItems+2)
	s.WriteString(ansi.SetModeAutoWrap)
	return s.String()
}
