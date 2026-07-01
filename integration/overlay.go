package integration

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/commands/core"
	"golang.org/x/term"
)

const (
	boxWidth = 72
	maxItems = 6
)

type Overlay struct {
	mu            sync.Mutex
	Visible       bool
	Items         []core.Suggestion
	Cursor        int
	StartIdx      int
	LastGhostLen  int
	TypedQuery    string
	UserNavigated bool
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
		Visible:  false,
		Cursor:   0,
		StartIdx: 0,
	}
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
	runes := []rune(s)
	if len(runes) > width {
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

	// add extra padding to erase any stray characters left by fast backspaces
	// before the debounce timer fired, 10 spaces is safe and won't hit right prompts
	padLen += 10

	if ghostText != "" || padLen > 0 {
		s.WriteString("\0337") // save cursor at prompt
		if ghostText != "" {
			s.WriteString("\033[90m")
			s.WriteString(ghostText)
			s.WriteString("\033[0m")
		}
		if padLen > 0 {
			s.WriteString(strings.Repeat(" ", padLen))
		}
		s.WriteString("\0338") // restore cursor back to prompt
		o.LastGhostLen = len(ghostText)
	}

	return s.String()
}

func (o *Overlay) Render() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.Visible || len(o.Items) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString("\033[?7l")

	var offset int
	if o.UserNavigated && len(o.Items) > 0 && o.Cursor >= 0 && o.Cursor < len(o.Items) {
		currentCmd := o.Items[o.Cursor].Cmd
		typedLen := len([]rune(o.TypedQuery))
		currentLen := len([]rune(currentCmd))
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || width <= 0 {
			width = 120
		}
		if currentLen+4 < width {
			offset = currentLen - typedLen
		} else {
			curCol := (currentLen + 2) % width
			typedCol := (typedLen + 2) % width
			offset = curCol - typedCol
		}
	}

	s.WriteString("\0337")

	windowSize := maxItems
	if len(o.Items) < windowSize {
		windowSize = len(o.Items)
	}

	scrolloffUp := 1
	scrolloffDown := 0
	if windowSize <= 3 {
		scrolloffUp = 0
	}

	if o.Cursor < o.StartIdx+scrolloffUp {
		o.StartIdx = o.Cursor - scrolloffUp
	}
	if o.Cursor >= o.StartIdx+windowSize-scrolloffDown {
		o.StartIdx = o.Cursor - windowSize + scrolloffDown + 1
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

	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", 1)
	s.WriteString("\033[2K")
	if offset > 0 {
		fmt.Fprintf(&s, "\033[%dD", offset)
	} else if offset < 0 {
		fmt.Fprintf(&s, "\033[%dC", -offset)
	}

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
		s.WriteString("\033[2K")
		if offset > 0 {
			fmt.Fprintf(&s, "\033[%dD", offset)
		} else if offset < 0 {
			fmt.Fprintf(&s, "\033[%dC", -offset)
		}

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

	s.WriteString("\0338")
	fmt.Fprintf(&s, "\033[%dB", windowSize+2)
	s.WriteString("\033[2K")
	if offset > 0 {
		fmt.Fprintf(&s, "\033[%dD", offset)
	} else if offset < 0 {
		fmt.Fprintf(&s, "\033[%dC", -offset)
	}
	bottomBorder := "╰" + strings.Repeat("─", boxWidth) + "╯"
	s.WriteString(borderStyle.Render(bottomBorder))

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
