package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fastlane",
		Description: "Helps you with your initial fastlane setup",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Helps you with your initial fastlane setup"},
			{Name: "swift", Description: "Fastlane configuration written in Swift (Beta). Swift setup is still in beta"},
			{Name: "appleID", Description: "Your Apple ID"},
			{Name: "action", Description: "Shows more information for a specific command"},
			{Name: "actions", Description: "Lists all available fastlane actions"},
			{Name: "add_plugin", Description: "Add a new plugin to your fastlane setup"},
			{Name: "docs", Description: "Generate a markdown based documentation based on the Fastfile"},
			{Name: "enable_auto_complete", Description: "Enable tab auto completion"},
			{Name: "env", Description: "Print your fastlane environment, use this when you submit an issue on GitHub"},
			{Name: "help", Description: "Display global or [command] help documentation"},
			{Name: "install_plugins", Description: "Install all plugins for this project"},
			{Name: "lanes", Description: "Lists all available lanes and shows their description"},
			{Name: "list", Description: "Lists all available lanes without description"},
			{Name: "new_action", Description: "Create a new custom action for fastlane"},
			{Name: "new_plugin", Description: "Create a new plugin that can be used with fastlane"},
			{Name: "run", Description: "Run a fastlane one-off action without a full lane"},
			{Name: "search_plugins", Description: "Search for plugins, search query is optional"},
			{Name: "socket_server", Description: "Starts local socket server and enables only a single local connection"},
			{Name: "seconds", Description: "Connection timeout in seconds"},
			{Name: "port", Description: "The port on localhost"},
			{Name: "trigger", Description: "Run a specific lane. Pass the lane name and optionally the platform first"},
			{Name: "lane", Description: "Specific lane to trigger"},
			{Name: "update_fastlane", Description: "Update fastlane to the latest release"},
			{Name: "update_plugins", Description: "Update all plugin dependencies"},
		},
		Options: []core.Option{
			{Name: "-u", Description: "Only iOS projects Your Apple ID"},
			{Name: "-f", Description: "Overwrite the existing README.md in the ./fastlane folder"},
			{Name: "-c", Description: "Add custom command(s) for which tab auto complete should be enabled too"},
			{Name: "-j", Description: "Output the lanes in JSON instead of text"},
			{Name: "--name", Description: "Name of your new action"},
			{Name: "-s", Description: "Keeps socket server up even after error or disconnects, requires CTRL-C to kill"},
			{Name: "-p", Description: "Sets the port on localhost for the socket connection"},
			{Name: "--disable_runner_upgrades", Description: "Prevents fastlane from attempting to update FastlaneRunner swift project"},
			{Name: "--swift_server_port", Description: "Prevents fastlane from attempting to update FastlaneRunner swift project"},
			{Name: "--platform", Description: "Only show actions available on the given platform"},
			{Name: "-h", Description: "Show help for fastlane"},
			{Name: "-v", Description: "Show version information for fastlane"},
			{Name: "--verbose", Description: "Show version information for fastlane"},
			{Name: "--capture_output", Description: "Captures the output of the current run, and generates a markdown issue template"},
			{Name: "--troubleshoot", Description: "Add environment(s) to use with `dotenv`"},
		},
	})
}
