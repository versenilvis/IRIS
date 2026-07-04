package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rush",
		Description: "Projects",
		Subcommands: []core.Subcommand{
			{Name: "add", Description: "Adds one or more dependencies to the package.json and runs rush update"},
			{Name: "init", Description: "Initializes a new repository to be managed by Rush"},
			{Name: "init-autoinstaller", Description: "Initializes a new autoinstaller"},
			{Name: "link", Description: "Create node_modules symlinks for all projects"},
			{Name: "list", Description: "List package information for all projects in the repo"},
			{Name: "unlink", Description: "Delete node_modules symlinks for all projects in the repo"},
			{Name: "update-autoinstaller", Description: "Updates autoinstaller package dependenices"},
			{Name: "update-cloud-credentials", Description: "(EXPERIMENTAL) Update the credentials used by the build cache provider"},
			{Name: "version", Description: "Manage package versions in the repo"},
			{Name: "rebuild", Description: "Clean and rebuild the entire set of projects"},
			{Name: "tab-complete", Description: "Provides tab completion"},
		},
		Options: []core.Option{
			{Name: "-t", Description: "Multi-Project Build Tool"},
			{Name: "-p", Description: "Specify package to the dependencies of the current project"},
			{Name: "--exact", Description: "If specified, the dependency will be added to all projects"},
			{Name: "-v", Description: "Verify the change file has been generated and that it is a valid JSON file"},
			{Name: "--no-fetch", Description: "The message to apply to all changed projects if the --bulk flag is provided"},
			{Name: "--bump-type", Description: "The bump type to apply to all changed projects if the --bulk flag is provided"},
			{Name: "--variant", Description: "If this flag is specified, output will be in JSON format"},
			{Name: "--overwrite-existing", Description: "Initializes a new autoinstaller"},
			{Name: "--name", Description: "Overrides the default maximum number of install attempts. The default value is 3"},
			{Name: "--ignore-hooks", Description: "Only check the validity of the shrinkwrap file without performing an install"},
			{Name: "-f", Description: "List package information for all projects in the repo"},
			{Name: "-a", Description: "If this flag is specified, applied changes will be published to the NPM registry"},
			{Name: "--add-commit-details", Description: "Adds commit author and hash to the changelog.json files for each change"},
			{Name: "--regenerate-changelogs", Description: "Regenerates all changelog files based on the current JSON content"},
			{Name: "-r", Description: "Append a suffix to all changed versions. Cannot be used with --prerelease-name"},
			{Name: "--force", Description: "Skips execution of all git hooks. Make sure you know what you are skipping"},
			{Name: "--unsafe", Description: "If this flag is specified, output will be in JSON format"},
			{Name: "--all", Description: "If this flag is specified, output will list all detected dependencies"},
			{Name: "-i", Description: "A static credential, to be cached"},
			{Name: "-d", Description: "If specified, delete stored credentials"},
			{Name: "-b", Description: "Updates package versions if needed to satisfy version policies"},
			{Name: "--override-version", Description: "Bumps package version based on version policies"},
			{Name: "--bypass-policy", Description: "The name of the version policy"},
			{Name: "--override-bump", Description: "Skips execution of all git hooks. Make sure you know what you are skipping"},
			{Name: "--word", Description: "The position in the word to be completed. The default value is 0"},
			{Name: "-h", Description: "Show this help message and exit"},
		},
	})
}
