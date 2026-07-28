package integration

import (
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/versenilvis/iris/internal/logger"
	"github.com/versenilvis/iris/spec"
	"golang.org/x/term"
)

const (
	boxWidth = 76 // total visual width, corners included
	maxItems = 6
)

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
	var sb strings.Builder
	w := 0
	for _, r := range s {
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
