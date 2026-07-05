package python

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pipx",
		Description: "Installed package",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "Install a package"},
			{Name: "package_spec", Description: "Package name or pip installation spec"},
			{Name: "inject", Description: "Installs packages to an existing pipx-managed virtual environment"},
			{Name: "package", Description: "Name of the existing pipx-managed Virtual Environment to inject into"},
			{Name: "uninstall-all", Description: "Uninstall all pipx-managed packages"},
			{Name: "reinstall", Description: "Reinstalls a package"},
			{Name: "reinstall-all", Description: "Reinstalls all packages"},
			{Name: "list", Description: "List packages and apps installed with pipx"},
			{Name: "app", Description: "App/package name and any arguments to be passed to it"},
			{Name: "runpip", Description: "Run pip in an existing pipx-managed Virtual Environment"},
			{Name: "pipargs", Description: "Arguments to forward to pip command"},
			{Name: "ensurepath", Description: "Ensure directory where pipx stores apps is in your PATH environment variable"},
			{Name: "environment", Description: "Print a list of variables used in pipx.constants"},
			{Name: "completions", Description: "Print instructions on enabling shell completions for pipx"},
		},
		Options: []spec.Option{
			{Name: "--include-deps", Description: "Include apps of dependent packages"},
			{Name: "--force", Description: "Modify existing virtual environment and files in PIPX_BIN_DIR"},
			{Name: "--suffix", Description: "Optional suffix for virtual environment and executable names"},
			{Name: "--python", Description: "Give the virtual environment access to the system site-packages dir"},
			{Name: "--index-url", Description: "Base URL of Python Package Index"},
			{Name: "--editable", Description: "Install a project in editable mode"},
			{Name: "--pip-args", Description: "Arbitrary pip arguments to pass directly to pip install/upgrade commands"},
			{Name: "--verbose", Description: "Show verbose output"},
			{Name: "--include-apps", Description: "Add apps from the injected packages onto your PATH"},
			{Name: "--system-site-packages", Description: "Give the virtual environment access to the system site-packages dir"},
			{Name: "--include-injected", Description: "Also upgrade packages injected into the main app's environment"},
			{Name: "--skip", Description: "Skip these packages"},
			{Name: "--json", Description: "Output rich data in json format"},
			{Name: "--short", Description: "List packages only"},
			{Name: "--no-cache", Description: "Do not re-use cached virtual environment if it exists"},
			{Name: "--pypackages", Description: "Require app to be run from local __pypackages__ directory"},
			{Name: "--spec", Description: "The package name or specific installation source passed to pip"},
			{Name: "--value", Description: "Print the value of the variable"},
			{Name: "--help", Description: "Show help for pipx"},
		},
	})
}
