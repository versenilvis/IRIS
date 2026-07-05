package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "nvm",
		Description: "Node version",
		Subcommands: []spec.Subcommand{
			{Name: "uninstall", Description: "Uninstall a version"},
			{Name: "use", Description: "Modify PATH to use <version>. Uses .nvmrc if available and version is omitted"},
			{Name: "exec", Description: "Run <command> on <version>. Uses .nvmrc if available and version is omitted"},
			{Name: "current", Description: "Display currently activated version of Node"},
			{Name: "ls", Description: "List installed versions, matching a given <version> if provided"},
			{Name: "version", Description: "Resolve the given description to a single local version"},
			{Name: "version-remote", Description: "Resolve the given description to a single remote version"},
			{Name: "deactivate", Description: "Undo effects of `nvm` on current shell"},
			{Name: "pattern or name", Description: "Pattern or name"},
			{Name: "unalias", Description: "Deletes the alias named <name>"},
			{Name: "install-latest-npm", Description: "Attempt to upgrade to the latest working `npm` on the current node version"},
			{Name: "reinstall-packages", Description: "Reinstall global `npm` packages contained in <version> to current version"},
			{Name: "unload", Description: "Unload `nvm` from shell"},
			{Name: "dir", Description: "Display path to the cache directory for nvm"},
			{Name: "clear", Description: "Empty cache directory for nvm"},
		},
		Options: []spec.Option{
			{Name: "--no-colors", Description: "Suppress colored output"},
			{Name: "--no-alias", Description: "Suppress `nvm alias` output"},
			{Name: "--silent", Description: "Silences stdout/stderr output"},
			{Name: "--lts", Description: "Uses automatic LTS (long-term support) alias `lts/*`, if available"},
			{Name: "-s", Description: "Skip binary download, install from source only"},
			{Name: "--reinstall-packages-from", Description: "When installing, reinstall packages installed in <version>"},
			{Name: "--skip-default-packages", Description: "When installing, skip the default-packages file if it exists"},
			{Name: "--latest-npm", Description: "Disable the progress bar on any downloads"},
			{Name: "--alias", Description: "Uninstall a version"},
			{Name: "--help", Description: "Show help page"},
			{Name: "--version", Description: "Print out the installed version of nvm"},
		},
	})
}
