package integration

import "github.com/charmbracelet/lipgloss"

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
