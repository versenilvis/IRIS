package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "browser-sync",
		Description: "Keep multiple browsers & devices in sync when building websites",
		Subcommands: []core.Subcommand{
			{Name: "start", Description: "Start Browsersync"},
			{Name: "init", Description: "Create a configuration file"},
			{Name: "reload", Description: "Send a reload event over HTTP protocol"},
			{Name: "recipe", Description: "Generate the files for a recipe"},
			{Name: "ls", Description: "List recipes"},
		},
		Options: []core.Option{
			{Name: "--version", Description: "Show version number"},
			{Name: "--help", Description: "Show help"},
			{Name: "-s", Description: "Run a Local server (uses your cwd as the web root)"},
			{Name: "--cwd", Description: "Working directory"},
			{Name: "--json", Description: "If true, certain logs will output as json only"},
			{Name: "--serveStatic", Description: "Directories to serve static files from"},
			{Name: "--port", Description: "Specify a port to use"},
			{Name: "-p", Description: "Proxy an existing server"},
			{Name: "--ws", Description: "Proxy mode only - enable websocket proxying"},
			{Name: "-b", Description: "Choose which browser should be auto-opened"},
			{Name: "-w", Description: "Watch files"},
			{Name: "--ignore", Description: "Ignore patterns for file watchers"},
			{Name: "-f", Description: "File paths to watch"},
			{Name: "--index", Description: "Specify which file should be used as the index page"},
			{Name: "--plugins", Description: "Load Browsersync plugins"},
			{Name: "--extensions", Description: "Specify file extension fallbacks"},
			{Name: "--startPath", Description: "Specify the start path for the opened browser"},
			{Name: "--single", Description: "If true, the connect-history-api-fallback middleware will be added"},
			{Name: "--https", Description: "Enable SSL for local development"},
			{Name: "--directory", Description: "Show a directory listing for the server"},
			{Name: "--xip", Description: "Use xip.io domain routing"},
			{Name: "--tunnel", Description: "Use a public URL"},
			{Name: "--open", Description: "Choose which URL is auto-opened (local, external or tunnel), or provide a url"},
			{Name: "--cors", Description: "Add Access Control headers to every request"},
			{Name: "-c", Description: "Specify a path to a configuration file"},
			{Name: "--host", Description: "Specify a hostname to use"},
			{Name: "--listen", Description: "Specify a hostname bind to (this will prevent binding to all interfaces)"},
			{Name: "--logLevel", Description: "Set the logger output level (silent, info or debug)"},
			{Name: "--reload-delay", Description: "Time in milliseconds to delay the reload event following file changes"},
			{Name: "--reload-debounce", Description: "Specify a port for the UI to use"},
			{Name: "--watchEvents", Description: "Specify which file events to respond to"},
			{Name: "--no-notify", Description: "Disable the notify element in browsers"},
			{Name: "--no-open", Description: "Don't open a new browser window"},
			{Name: "--no-snippet", Description: "Disable the snippet injection"},
			{Name: "--no-online", Description: "Force offline usage"},
			{Name: "--no-ui", Description: "Don't start the user interface"},
			{Name: "--no-ghost-mode", Description: "Disable Ghost Mode"},
			{Name: "--no-inject-changes", Description: "Reload on every file change"},
			{Name: "--no-reload-on-restart", Description: "Don't auto-reload all browsers following a restart"},
			{Name: "-u", Description: "Provide the full URL to the running browsersync isntance"},
		},
	})
}
