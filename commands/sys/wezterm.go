package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "wezterm",
		Description: "Wez",
		Subcommands: []spec.Subcommand{
			{Name: "start", Description: "Start the GUI, optionally running an alternative program"},
			{Name: "ssh", Description: "Establish an ssh session"},
			{Name: "serial", Description: "Open a serial port"},
			{Name: "connect", Description: "Connect to wezterm multiplexer"},
			{Name: "ls-fonts", Description: "Display information about fonts"},
			{Name: "show-keys", Description: "Show key assignments"},
			{Name: "cli", Description: "Interact with experimental mux server"},
			{Name: "list", Description: "List windows, tabs and panes"},
			{Name: "list-clients", Description: "List clients"},
			{Name: "proxy", Description: "Start rpc proxy pipe"},
			{Name: "tlscreds", Description: "Obtain tls credentials"},
			{Name: "move-pane-to-new-tab", Description: "Move a pane into a new tab"},
			{Name: "help", Description: "Print this message or the help of the given subcommand(s)"},
			{Name: "imgcat", Description: "Output an image to the terminal"},
			{Name: "record", Description: "Record a terminal session as an asciicast"},
			{Name: "replay", Description: "Replay an asciicast terminal session"},
			{Name: "shell-completion", Description: "Generate shell completion information"},
		},
		Options: []spec.Option{
			{Name: "--cwd", Description: "Specify the current working directory for the initially spawned program"},
			{Name: "--class", Description: "Override the position for the initial window launched by this process"},
			{Name: "--no-auto-connect", Description: "Print help information"},
			{Name: "-o", Description: "Override the position for the initial window launched by this process"},
			{Name: "-v", Description: "Print help information"},
			{Name: "--baud", Description: "Set the baud rate.  The default is 9600 baud"},
			{Name: "-h", Description: "Print help information"},
			{Name: "--text", Description: "Explain which fonts are used to render the supplied text string"},
			{Name: "--list-system", Description: "Whether to list all fonts available to the system"},
			{Name: "--rasterize-ascii", Description: "Show rasterized glyphs for the text in --text using ascii blocks"},
			{Name: "--key-table", Description: "In lua mode, show only the named key table"},
			{Name: "--lua", Description: "Show the keys as lua config statements"},
			{Name: "--format", Description: "Print help information"},
			{Name: "--pane-id", Description: "Create tab in a new window, rather than the window currently containing the pane"},
			{Name: "--move-pane-id", Description: "Split horizontally, with the new pane on the left"},
			{Name: "--right", Description: "Split horizontally, with the new pane on the right"},
			{Name: "--top", Description: "Split vertically, with the new pane on the top"},
			{Name: "--bottom", Description: "Split vertically, with the new pane on the bottom"},
			{Name: "--top-level", Description: "Rather than splitting the active pane, split the entire window"},
			{Name: "--workspace", Description: "Spawn into a new window, rather than a new tab"},
			{Name: "--prefer-mux", Description: "Print help information"},
			{Name: "--width", Description: "Do not respect the aspect ratio.  The default is to respect the aspect ratio"},
			{Name: "--explain", Description: "Explain what is being sent/received"},
			{Name: "--shell", Description: "Which shell to generate for"},
			{Name: "--config-file", Description: "Override specific configuration values"},
			{Name: "-V", Description: "Print version information"},
			{Name: "-n", Description: "Skip loading wezterm.lua"},
		},
	})
}
