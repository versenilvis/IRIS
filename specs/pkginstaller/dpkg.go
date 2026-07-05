package pkginstaller

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dpkg",
		Description: "Debian package management system",
		Subcommands: []core.Subcommand{
			{Name: "query", Description: "Query the dpkg database"},
			{Name: "install", Description: "Install or upgrade packages"},
			{Name: "remove", Description: "Remove packages"},
			{Name: "purge", Description: "Remove packages and their configuration files"},
			{Name: "configure", Description: "Configure a package after installation"},
			{Name: "reconfigure", Description: "Reconfigure a package"},
			{Name: "list", Description: "List packages in the dpkg database"},
			{Name: "builddeb", Description: "Build Debian package files from sources"},
			{Name: "check", Description: "Check the dependencies of packages"},
			{Name: "compare-versions", Description: "Compare package versions"},
		},
		Options: []core.Option{
			{Name: "-l", Description: "List packages matching a pattern"},
			{Name: "-s", Description: "Display the status of available packages"},
			{Name: "-L", Description: "List files in a package"},
			{Name: "-S", Description: "Search for a package owning a file"},
			{Name: "-p", Description: "Display details about a package in the dpkg database"},
			{Name: "-W", Description: "Show a package in the dpkg database"},
			{Name: "-R", Description: "Recursively handle packages"},
			{Name: "-B", Description: "Uninstall packages that depend on the target package"},
			{Name: "--skip-same-version", Description: "Don't upgrade if the same version is already installed"},
			{Name: "--force-depends", Description: "Ignore dependency problems"},
			{Name: "--force-confnew", Description: "Always install the new version of configuration files"},
			{Name: "--force-confold", Description: "Always install the old version of configuration files"},
			{Name: "--force-confdef", Description: "Always install the default version of configuration files"},
			{Name: "--force-confmiss", Description: "Always install missing configuration files"},
			{Name: "--no-triggers", Description: "Skip processing triggers"},
			{Name: "--no-act", Description: "Simulate the action, but don't execute"},
			{Name: "--pending", Description: "Configure all unconfigured packages"},
			{Name: "--installed", Description: "List installed packages"},
			{Name: "--avail", Description: "List available packages"},
			{Name: "--hold", Description: "List packages on hold"},
			{Name: "--deferred", Description: "List deferred packages"},
			{Name: "-us", Description: "Build unsigned .changes and .dsc files"},
			{Name: "-uc", Description: "Build unsigned .changes file"},
			{Name: "-sa", Description: "Build source package"},
			{Name: "-rfakeroot", Description: "Use fakeroot when building the package"},
			{Name: "-b", Description: "Build binary package from source"},
			{Name: "--force-sign", Description: "Force signing of changes file"},
			{Name: "-m", Description: "Specify the package maintainer"},
			{Name: "-c", Description: "Specify the changes file to use"},
			{Name: "-v", Description: "Specify the version to use"},
			{Name: "--increment", Description: "Increment the version number in the changelog"},
			{Name: "-a", Description: "Check all installed packages"},
			{Name: "-d", Description: "Check for unmet dependencies"},
			{Name: "-i", Description: "Check installed packages"},
			{Name: "-U", Description: "Check unpacked packages"},
			{Name: "-r", Description: "Check reverse dependencies"},
			{Name: "-h", Description: "Help for dpkg"},
		},
	})
}
