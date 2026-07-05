package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fvm",
		Description: "Print this usage information",
		Subcommands: []core.Subcommand{
			{Name: "config", Description: "Set configuration for FVM"},
			{Name: "path", Description: "Path to the Flutter versions cache"},
			{Name: "dart", Description: "Proxies Dart commands"},
			{Name: "doctor", Description: "Shows information about environment, and project configuration"},
			{Name: "flavor", Description: "Switches between different project flavors"},
			{Name: "flavor-name", Description: "The flavor to switch to"},
			{Name: "flutter", Description: "Proxies Flutter commands"},
			{Name: "global", Description: "Sets Flutter SDK version as global"},
			{Name: "version", Description: "Flutter SDK to set for global flutter command"},
			{Name: "install", Description: "Installs Flutter SDK version"},
			{Name: "stable", Description: "Latest stable release of Flutter"},
			{Name: "beta", Description: "Latest beta release of Flutter"},
			{Name: "dev", Description: "Latest dev release of Flutter (master)"},
			{Name: "list", Description: "Lists installed Flutter SDK versions"},
			{Name: "releases", Description: "View all Flutter SDK releases available for install"},
			{Name: "remove", Description: "Removes Flutter SDK version"},
			{Name: "spawn", Description: "Spawn a Flutter SDK version command"},
			{Name: "use", Description: "Sets a Flutter SDK version to be used in a project"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Print this usage information"},
			{Name: "--verbose", Description: "Print verbose output"},
			{Name: "--version", Description: "Current FVM version"},
			{Name: "-c", Description: "Set the path which FVM will cache the version. Priority over FVM_HOME"},
			{Name: "-s", Description: "Skip setup after a version install"},
			{Name: "-g", Description: "ADVANCED: Will cache a local version of Flutter repo for faster version install"},
			{Name: "--force", Description: "Skips version global check"},
			{Name: "-f", Description: "Skips command guards that does Flutter project checks"},
			{Name: "-p", Description: "If version provided is a channel. Will pin the latest release of the channel"},
			{Name: "--flavor", Description: "Sets version for a project flavor"},
		},
	})
}
