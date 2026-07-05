package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "zellij",
		Description: "Change where zellij looks for the configuration file",
		Options: []core.Option{
			{Name: "-c", Description: "Change where zellij looks for the configuration file"},
			{Name: "--config-dir", Description: "Change where zellij looks for the configuration directory"},
			{Name: "-d", Description: "Specify emitting additional debug information"},
			{Name: "--data-dir", Description: "Change where zellij looks for plugins"},
			{Name: "-h", Description: "Print help information"},
			{Name: "-l", Description: "Maximum panes on screen, caution: opening more panes will close old ones"},
			{Name: "-n", Description: "Specify name of a new session"},
			{Name: "-V", Description: "Print version information"},
			{Name: "-b", Description: "Create a detached session in the background if one does not exist"},
			{Name: "-f", Description: "If resurrecting a dead session, immediately run all its commands on startup"},
			{Name: "--index", Description: "Number of the session index in the active sessions ordered creation date"},
			{Name: "-y", Description: "Automatic yes to prompts"},
			{Name: "--cwd", Description: "Change the working directory of the editor"},
			{Name: "--height", Description: "Open the new pane in place of the current pane, temporarily suspending it"},
			{Name: "--width", Description: "The width if the pane is floating as a bare integer (eg. 1) or percent (eg. 10%)"},
			{Name: "-x", Description: "Print this message or the help of the given subcommand(s)"},
			{Name: "-r", Description: "List the sessions in reverse order"},
			{Name: "-s", Description: "Print just the session name"},
			{Name: "--attach-to-session", Description: "Whether to attach to a session specified in 'session-name' if it exists"},
			{Name: "--auto-layout", Description: "Whether to lay out panes in a predefined set of layouts whenever possible"},
			{Name: "--copy-clipboard", Description: "OSC52 destination clipboard"},
			{Name: "--copy-command", Description: "Switch to using a user supplied command for clipboard instead of OSC52"},
			{Name: "--copy-on-select", Description: "Automatically copy when selecting text (true or false)"},
			{Name: "--default-cwd", Description: "Set the default cwd"},
			{Name: "--default-layout", Description: "Set the default layout"},
			{Name: "--default-mode", Description: "Set the default mode"},
			{Name: "--default-shell", Description: "Set the default shell"},
			{Name: "--disable-mouse-mode", Description: "Disable handling of mouse events"},
			{Name: "--disable-session-metadata", Description: "If true, will disable writing session metadata to disk"},
			{Name: "--layout-dir", Description: "Set the layout_dir, defaults to subdirectory of config dir"},
			{Name: "--mirror-session", Description: "Mirror session when multiple users are connected"},
			{Name: "--mouse-mode", Description: "Disable display of pane frames"},
			{Name: "--on-force-close", Description: "Set behaviour on force close (quit or detach)"},
			{Name: "--pane-frames", Description: "Set display of the pane frames (true or false)"},
			{Name: "--scroll-buffer-size", Description: "Explicit full path to open the scrollback editor (default is $EDITOR or $VISUAL)"},
			{Name: "--serialize-pane-viewport", Description: "Whether pane viewports are serialized along with the session, default is false"},
			{Name: "--session-name", Description: "The name of the session to create when starting Zellij"},
			{Name: "--session-serialization", Description: "Whether sessions should be serialized to the HD"},
			{Name: "--simplified-ui", Description: "Allow plugins to use a more simplified layout that is compatible with more fonts"},
			{Name: "--styled-underlines", Description: "Whether to use ANSI styled underlines"},
		},
	})
}
