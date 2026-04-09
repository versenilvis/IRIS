package integration

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/commands/core"
)

const (
	boxWidth = 72
	maxItems = 6 // max items showing in the menu preview
)

type Overlay struct {
	mu      sync.Mutex
	Visible bool
	Items   []core.Suggestion
	Cursor  int
}

var (
	selBgColor = lipgloss.Color("#44475A")

	iconStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	descStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))

	iconSel  = iconStyle.Background(selBgColor)
	titleSel = titleStyle.Background(selBgColor)
	descSel  = descStyle.Background(selBgColor)
	padSel   = lipgloss.NewStyle().Background(selBgColor)

	borderColor = lipgloss.Color("#6272A4")
	borderStyle = lipgloss.NewStyle().Foreground(borderColor)
)

func NewOverlay() *Overlay {
	return &Overlay{
		Visible: false,
		Cursor:  0,
	}
}

func (o *Overlay) UpdateItems(items []core.Suggestion) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.Items = items
	o.Visible = len(o.Items) > 0
	o.Cursor = 0 // reset to top result on update
}

// fixedWidth pads or truncates a string to exact rune width
func fixedWidth(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		return string(runes[:width-1]) + "…"
	}
	if len(runes) < width {
		return s + strings.Repeat(" ", width-len(runes))
	}
	return s
}

func (o *Overlay) Render() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString("\033[?7l")
	s.WriteString("\0337") // DEC save cursor

	windowSize := maxItems
	if len(o.Items) < windowSize {
		windowSize = len(o.Items)
	}

	start := 0
	if o.Cursor >= windowSize {
		start = o.Cursor - windowSize + 1
	}
	end := start + windowSize

	totalLines := windowSize + 2 // top border + items + bottom border

	// if we reach the last lines of terminal, it will auto expand space to have space for the menu
	for range totalLines {
		s.WriteByte('\n')
	}
	fmt.Fprintf(&s, "\033[%dA", totalLines)

	// re-save after scroll up
	// it means when you scroll up, then scroll down to the prompt, the menu is still be there
	s.WriteString("\0337")

	// top border with scroll indicator
	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", 1)
	s.WriteString("\033[K")

	scrollInfo := ""
	if len(o.Items) > windowSize {
		scrollInfo = fmt.Sprintf(" %d/%d ", o.Cursor+1, len(o.Items))
	}

	borderWidth := boxWidth - len(scrollInfo)
	topBorder := "╭" + strings.Repeat("─", borderWidth/2) + scrollInfo + strings.Repeat("─", boxWidth-borderWidth/2-len(scrollInfo)) + "╮"
	s.WriteString(borderStyle.Render(topBorder))

	iconW := 6
	descW := 22
	titleW := boxWidth - iconW - descW - 3

	for i := start; i < end; i++ {
		s.WriteString("\0338")
		fmt.Fprintf(&s, "\033[%dB", (i-start)+2)
		s.WriteString("\033[K")

		it := o.Items[i]
		rawIcon := fixedWidth(it.Icon, iconW)
		rawTitle := fixedWidth(it.Cmd, titleW)
		rawDesc := fixedWidth(it.Desc, descW)

		left := borderStyle.Render("│")
		right := borderStyle.Render("│")

		if i == o.Cursor {
			icon := iconSel.Render(" " + rawIcon + " ")
			title := titleSel.Render(rawTitle)
			pad := padSel.Render(" ")
			desc := descSel.Render(rawDesc)
			fmt.Fprintf(&s, "%s%s%s%s%s%s", left, icon, pad, title, desc, right)
		} else {
			icon := iconStyle.Render(" " + rawIcon + " ")
			title := titleStyle.Render(rawTitle)
			desc := descStyle.Render(rawDesc)
			fmt.Fprintf(&s, "%s%s %s%s%s", left, icon, title, desc, right)
		}
	}

	// Bottom border
	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", windowSize+2)
	s.WriteString("\033[K")
	bottomBorder := "╰" + strings.Repeat("─", boxWidth) + "╯"
	s.WriteString(borderStyle.Render(bottomBorder))

	s.WriteString("\0338") // restore to prompt
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

	if !o.Visible {
		return ""
	}

	o.Visible = false
	o.Items = nil

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
