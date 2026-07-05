package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "firefox",
		Description: "Free open-source web browser developer by Mozilla",
		Options: []spec.Option{
			{Name: "--display", Description: "Specify an X display to use"},
			{Name: "--sync", Description: "Make X calls synchronous"},
			{Name: "--g-fatal-warnings", Description: "Make all warnings fatal"},
			{Name: "-h", Description: "Print help message and exit"},
			{Name: "-v", Description: "Print version information and exit"},
			{Name: "--full-version", Description: "Print full version information and exit"},
			{Name: "-P", Description: "Specify profile to use"},
			{Name: "--profile", Description: "Specify profile to use by folder"},
			{Name: "--migration", Description: "Start with migration wizard"},
			{Name: "--ProfileManager", Description: "Start with ProfileManager"},
			{Name: "--no-remote", Description: "Do not accept or send remote commands; implies --new-instance"},
			{Name: "--new-instance", Description: "Open new instance, not a new window in running instance"},
			{Name: "--safe-mode", Description: "Disables extensions and themes for this session"},
			{Name: "--allow-downgrade", Description: "Allows downgrading a profile"},
			{Name: "--MOZ_LOG", Description: "Treated as MOZ_LOG=<modules> environment variable, overrides it"},
			{Name: "--MOZ_LOG_FILE", Description: "Run without a GUI"},
			{Name: "--jsdebugger", Description: "Start the devtools server on a TCP port or Unix domain socket path"},
			{Name: "--browser", Description: "Open a browser window"},
			{Name: "--new-window", Description: "Open a URL in a new window"},
			{Name: "--new-tab", Description: "Open a URL in a new tab"},
			{Name: "--private-window", Description: "Open a URL in a new private window"},
			{Name: "--preferences", Description: "Open the preferences dialog"},
			{Name: "--screenshot", Description: "Take a screenshot"},
			{Name: "--window-size", Description: "Size of your screenshot"},
			{Name: "--search", Description: "Search for a term in your default search engine"},
			{Name: "--setDefaultBrowser", Description: "Set Firefox as the default browser"},
			{Name: "--first-startup", Description: "Run post-install actions before opening a new window"},
			{Name: "--kiosk", Description: "Start the browser in kiosk mode"},
			{Name: "--disable-pinch", Description: "Disable touch-screen and touch-pad pinch gestures"},
			{Name: "--jsconsole", Description: "Open the Browser Console"},
			{Name: "--devtools", Description: "Open DevTools on initial load"},
			{Name: "--marionette", Description: "Enable remote debugging server"},
			{Name: "--remote-debugging-port", Description: "Defaults to port 9222"},
			{Name: "--allow-remote-hosts", Description: "Values of the Host header to allow for incoming requests"},
			{Name: "--allow-remote-origins", Description: "Values of the Origin header to allow for incoming requests"},
		},
	})
}
