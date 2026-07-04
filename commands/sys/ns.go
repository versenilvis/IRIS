package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ns",
		Description: "Forces rebuilding the native application",
		Subcommands: []core.Subcommand{
			{Name: "next", Description: "The latest development release"},
		},
		Options: []core.Option{
			{Name: "--clean", Description: "Forces rebuilding the native application"},
			{Name: "--timeout", Description: "Specifies that you want to debug the app in an emulator"},
			{Name: "--device", Description: "Specifies a connected device/emulator to start and run the app"},
			{Name: "--force", Description: "Disables Hot Module Replacement (HMR)"},
			{Name: "--framework-path", Description: "Show the output of the command in JSON format"},
			{Name: "--justlaunch", Description: "Does not print the application output in the console"},
			{Name: "--release", Description: "Bundle the application"},
			{Name: "--help", Description: "Specifies that the command will produce and deploy an Android App Bundle"},
			{Name: "--key-store-path", Description: "Provides the password for the keystore file specified with --key-store-path"},
			{Name: "--key-store-alias", Description: "Provides the alias for the keystore file specified with --key-store-path"},
			{Name: "--key-store-alias-password", Description: "Provides the password for the alias specified with --key-store-alias-password"},
			{Name: "--sdk", Description: "Specifies the target simulator's sdk"},
			{Name: "--debug-brk", Description: "Attaches the debug tools to a deployed and running app"},
			{Name: "--no-watch", Description: "Changes in your code will not be livesynced"},
			{Name: "--no-client", Description: "Sets the Android SDK that will be used to build the project"},
			{Name: "--for-device", Description: "Open the CLI's documentation in your web browser"},
			{Name: "--insecure", Description: "Enables anonymous usage reporting"},
			{Name: "--template", Description: "Create a project using a predefined template"},
			{Name: "--angular", Description: "Create a project using the Angular template"},
			{Name: "--react", Description: "Create a project using the React template"},
			{Name: "--vue", Description: "Create a project using the Vue template"},
			{Name: "--svelte", Description: "Create a project using the Svelte template"},
			{Name: "--typescript", Description: "Create a project using plain TypeScript"},
			{Name: "--javascript", Description: "Create a project using plain JavaScript"},
			{Name: "--path", Description: "Clean your Nativescript project"},
			{Name: "--framework", Description: "Sets the unit testing framework to install"},
			{Name: "--pluginName", Description: "Used to set the default file and class names in the plugin source"},
			{Name: "--includeTypeScriptDemo", Description: "Specifies if a TypeScript demo should be created"},
			{Name: "--includeAngularDemo", Description: "Specifies if an Angular demo should be created"},
			{Name: "--background", Description: "Sets the background color of the splashscreen"},
			{Name: "--hmr", Description: "Enables the hot module replacement (HMR) feature"},
			{Name: "--ipa", Description: "Use the provided .ipa file instead of building the project"},
			{Name: "--team-id", Description: "Lists all available Android devices"},
			{Name: "--available-devices", Description: "Lists all available iOS devices"},
			{Name: "-v", Description: "View your current Nativescript CLI version"},
		},
	})
}
