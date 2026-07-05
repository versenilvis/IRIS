package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xcodes",
		Description: "Manage the Xcode versions installed on your Mac",
		Subcommands: []spec.Subcommand{
			{Name: "download", Description: "Download a specific version of Xcode"},
			{Name: "version", Description: "The version to install"},
			{Name: "install", Description: "Download and install a specific version of Xcode"},
			{Name: "installed", Description: "List the versions of Xcode that are installed"},
			{Name: "list", Description: "List all versions of Xcode available to install"},
			{Name: "select", Description: "Change the selected Xcode"},
			{Name: "version-or-path", Description: "Version or path of Xcode to select"},
			{Name: "uninstall", Description: "Uninstall a version of Xcode"},
			{Name: "update", Description: "Update the list of available versions of Xcode"},
			{Name: "signout", Description: "Clears the stored username and password"},
		},
		Options: []spec.Option{
			{Name: "--latest", Description: "Update and then install the latest non-prerelease version available"},
			{Name: "--latest-prerelease", Description: "The path to an aria2 executable. Searches $PATH by default"},
			{Name: "--no-aria2", Description: "Don't use aria2 to download Xcode, even if it's available"},
			{Name: "--directory", Description: "The directory where your Xcodes are installed. Defaults to /Applications"},
			{Name: "--data-source", Description: "The data source for available Xcode versions. (default: xcodereleases)"},
			{Name: "--color", Description: "Color the output"},
			{Name: "--no-color", Description: "Do not color the output"},
			{Name: "--path", Description: "Local path to Xcode.xip"},
			{Name: "--experimental-unxip", Description: "Use the experimental unxip functionality. May speed up unarchiving by up to 2-3x"},
			{Name: "-p", Description: "Print the path of the selected Xcode"},
			{Name: "--help", Description: "Show help information"},
		},
	})
}
