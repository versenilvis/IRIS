package integration

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/logger"
	"golang.org/x/term"
)

func renderMatchedTitle(title, typed string, selected bool, w int) string {
	t := theme()
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

	t := theme()
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
