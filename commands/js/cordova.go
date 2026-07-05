package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cordova",
		Description: "Print out the version of your cordova-cli install",
		Subcommands: []spec.Subcommand{
			{Name: "create", Description: "Create a project"},
			{Name: "help", Description: "Get help for a command"},
			{Name: "telemetry", Description: "Turn telemetry collection on or off"},
			{Name: "on", Description: "Turn telemetry collection on"},
			{Name: "off", Description: "Turn telemetry collection off"},
			{Name: "config", Description: "Set, get, delete, edit, and list global cordova options"},
			{Name: "ls", Description: "Lists contents of config file"},
			{Name: "edit", Description: "Opens the config file in an editor"},
			{Name: "get", Description: "Echo the config value to stdout"},
			{Name: "delete", Description: "Deletes the key from all configuration files"},
			{Name: "info", Description: "Generate project information"},
			{Name: "requirements", Description: "Checks and print out all the installation requirements for platforms specified"},
			{Name: "platform", Description: "Platform name(s) to build. If not specified, all platforms are built"},
			{Name: "add", Description: "Add specified platforms"},
			{Name: "remove", Description: "Remove specified platforms"},
			{Name: "update", Description: "Update specified platforms"},
			{Name: "list", Description: "List all installed and available platforms"},
			{Name: "plugin", Description: "Manage project plugins"},
			{Name: "prepare", Description: "Copy files into platform(s) for building"},
			{Name: "compile", Description: "Compile project for platform(s)"},
			{Name: "build", Description: "Build project for platform(s) (prepare + compile)"},
			{Name: "clean", Description: "Cleanup project from build artifacts"},
			{Name: "run", Description: "Run project (including prepare && compile)"},
			{Name: "serve", Description: "Run project with a local webserver (including prepare)"},
			{Name: "port", Description: "Port to use for local web server"},
		},
		Options: []spec.Option{
			{Name: "-d", Description: "Print out the version of your cordova-cli install"},
			{Name: "--no-update-notifier", Description: "Will disable updates check"},
			{Name: "--nohooks", Description: "Suppress executing hooks (taking RegExp hook patterns as parameters)"},
			{Name: "--no-telemetry", Description: "Disable telemetry collection for the current command"},
			{Name: "--template", Description: "Use a custom template located locally, in NPM, or GitHub"},
			{Name: "--nosave", Description: "Do not save"},
			{Name: "--link", Description: "Remove specified platforms"},
			{Name: "--save", Description: "Updates the version specified in **config.xml**"},
			{Name: "--searchpath", Description: "Don't search the registry for plugins"},
			{Name: "--debug", Description: "Build it for a device"},
			{Name: "--emulator", Description: "Platform name(s) to build. If not specified, all platforms are built"},
			{Name: "--list", Description: "Deploy a debug build. This is the default behavior unless --release is specified"},
			{Name: "--release", Description: "Deploy a release build"},
			{Name: "--noprepare", Description: "Skip preparing (available in Cordova v6.2 or later)"},
			{Name: "--nobuild", Description: "Skip building"},
			{Name: "--device", Description: "Deploy to a device"},
			{Name: "--target", Description: "Run project with a local webserver (including prepare)"},
		},
	})
}
