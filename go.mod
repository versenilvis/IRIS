module github.com/versenilvis/iris

go 1.26.1

require (
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/creack/pty v1.1.24
	github.com/spf13/cobra v1.10.2
	golang.org/x/sys v0.42.0
	golang.org/x/term v0.41.0
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.4.1 // indirect
	github.com/charmbracelet/x/ansi v0.11.6 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.15 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.9.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/versenilvis/fuzzy v0.1.0-rc // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
)

// this is my personal project, I want to point to my local folder to debug and develop
// the github version is the stable one
replace github.com/versenilvis/fuzzy => ../fuzzy
