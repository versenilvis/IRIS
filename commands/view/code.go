package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "code",
		Description: "Read from stdin (e.g.",
		Options: []core.Option{
			{Name: "-d", Description: "Compare two files with each other"},
			{Name: "-m", Description: "Add folder(s) to the last active window"},
			{Name: "-g", Description: "Open a file at the path on the specified line and character position"},
			{Name: "-n", Description: "Force to open a new window"},
			{Name: "-r", Description: "Force to open a file or folder in an already opened window"},
			{Name: "-w", Description: "Wait for the files to be closed before returning"},
			{Name: "--locale", Description: "The locale to use (e.g. en-US or zh-TW)"},
			{Name: "--user-data-dir", Description: "Print usage"},
			{Name: "--extensions-dir", Description: "Set the root path for extensions"},
			{Name: "--list-extensions", Description: "List the installed extensions"},
			{Name: "--show-versions", Description: "Show versions of installed extensions, when using --list-extensions"},
			{Name: "--category", Description: "Filters installed extensions by provided category, when using --list-extensions"},
			{Name: "--install-extension", Description: "Uninstalls an extension"},
			{Name: "--enable-proposed-api", Description: "Print version"},
			{Name: "--verbose", Description: "Print verbose output (implies --wait)"},
			{Name: "--log", Description: "Log level to use. Default is 'info' when unspecified"},
			{Name: "-s", Description: "Print process usage and diagnostics information"},
			{Name: "--prof-startup", Description: "Run CPU profiler during startup"},
			{Name: "--disable-extensions", Description: "Disable all installed extensions"},
			{Name: "--disable-extension", Description: "Disable an extension"},
			{Name: "--sync", Description: "Turn sync on or off"},
			{Name: "--inspect-extensions", Description: "Disable GPU hardware acceleration"},
			{Name: "--max-memory", Description: "Max memory size for a window (in Mbytes)"},
			{Name: "--telemetry", Description: "Shows all telemetry events which VS code collects"},
		},
	})
}
