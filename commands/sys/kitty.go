package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kitty",
		Description: "A cat like utility to display images in the terminal",
		Subcommands: []core.Subcommand{
			{Name: "close-tab", Description: "Close the specified tab(s)"},
			{Name: "close-window", Description: "Close the specified window(s)"},
			{Name: "detach-tab", Description: "Detach the specified tab"},
			{Name: "detach-window", Description: "Detach the specified window"},
			{Name: "disable-ligatures", Description: "Control ligature rendering for the specified windows/tabs"},
			{Name: "env", Description: "Change the environment variables seen by processing in newly launched windows"},
			{Name: "focus-tab", Description: "The active window in the specified tab will be focused"},
			{Name: "get-colors", Description: "Get the terminal colors for the specified window"},
			{Name: "set-tab-title", Description: "Set the title for the specified tab(s)"},
			{Name: "set-window-title", Description: "Set the title of the specified window(s)"},
			{Name: "signal-child", Description: "Send one or more signals to the foreground process in the specified window(s)"},
		},
		Options: []core.Option{
			{Name: "--align", Description: "Horizontal alignment for the displayed image"},
			{Name: "--place", Description: "Choose where on the screen to display the image"},
			{Name: "--scale-up", Description: "Mirror the image about a horizontal or vertical axis or both"},
			{Name: "--clear", Description: "Remove all images currently displayed on the screen"},
			{Name: "--transfer-mode", Description: "Which mechanism to use to transfer images to the terminal"},
			{Name: "--detect-support", Description: "Detect support for image display in the terminal"},
			{Name: "--detection-timeout", Description: "How long to wait for detection to complete before aborting"},
			{Name: "--print-window-size", Description: "Print the current terminal window size in pixels"},
			{Name: "--stdin", Description: "Read an image from stdin"},
			{Name: "-T", Description: "Set the OS window title"},
			{Name: "-C", Description: "Specify a path to the configuration file(s) to use"},
			{Name: "-o", Description: "Override individual configuration options"},
			{Name: "-d", Description: "Change to the specified directory when launching"},
			{Name: "--session", Description: "Path to a file containing the startup session"},
			{Name: "--hold", Description: "Remain open after child process exits"},
			{Name: "-1", Description: "Only a single instance of kitty will run"},
			{Name: "--instance-group", Description: "Kitty will open a new window in an existing instance and quit immediately"},
			{Name: "--listen-on", Description: "Tell kitty to listen on the specified address for control messages"},
			{Name: "--start-as", Description: "Control how the initial kitty window is created"},
			{Name: "-v", Description: "The current kitty version"},
			{Name: "-h", Description: "Display this help message"},
			{Name: "--to", Description: "An address for the kitty instance to control"},
			{Name: "-m", Description: "The tab to match"},
			{Name: "--self", Description: "Close the tab of the window this command is run in, rather than the active tab"},
			{Name: "--target-group", Description: "Close the specified group of tabs"},
			{Name: "-t", Description: "The tab to match"},
			{Name: "-a", Description: "Disable in all windows"},
			{Name: "--no-response", Description: "Don't wait for a response indicating the success of the action"},
			{Name: "-c", Description: "The window to match"},
			{Name: "--temporary", Description: "The title can be overwritten by escape sequences"},
		},
	})
}
