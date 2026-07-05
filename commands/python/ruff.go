package python

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ruff",
		Description: "Enable verbose logging",
		Subcommands: []spec.Subcommand{
			{Name: "graph", Description: "Generate a map of Python file dependencies or dependents"},
		},
		Options: []spec.Option{
			{Name: "-v", Description: "Enable verbose logging"},
			{Name: "-q", Description: "Print diagnostics, but nothing else"},
			{Name: "-s", Description: "Path to the `pyproject.toml` or `ruff.toml` file to use for configuration"},
			{Name: "--isolated", Description: "Ignore all configuration files"},
			{Name: "--help", Description: "Print help"},
			{Name: "--fix", Description: "Apply fixes to resolve lint violations"},
			{Name: "--unsafe-fixes", Description: "Include fixes that may not retain the original intent of the code"},
			{Name: "--show-fixes", Description: "Show an enumeration of all fixed lint violations"},
			{Name: "--diff", Description: "Run in watch mode by re-running whenever files change"},
			{Name: "--fix-only", Description: "Ignore any `# noqa` comments"},
			{Name: "--output-format", Description: "Specify file to write the linter output to (default: stdout)"},
			{Name: "--target-version", Description: "The minimum Python version that should be supported"},
			{Name: "--preview", Description: "Enable preview mode; checks will include unstable rules and fixes"},
			{Name: "--extension", Description: "Show counts for every rule with at least one violation"},
			{Name: "--add-noqa", Description: "Enable automatic additions of `noqa` directives to failing lines"},
			{Name: "--show-files", Description: "See the files Ruff will be run against with the current settings"},
			{Name: "--show-settings", Description: "See the settings Ruff will use to lint a given Python file"},
			{Name: "-h", Description: "Print help"},
			{Name: "--select", Description: "Comma-separated list of rule codes to enable (or ALL, to enable all rules)"},
			{Name: "--ignore", Description: "Comma-separated list of rule codes to disable"},
			{Name: "--extend-select", Description: "Like --select, but adds additional rule codes on top of the selected ones"},
			{Name: "--per-file-ignores", Description: "List of mappings from file pattern to code to exclude"},
			{Name: "--extend-per-file-ignores", Description: "Like --fixable, but adds additional rule codes on top of those already specified"},
			{Name: "--exclude", Description: "List of paths, used to omit files and/or directories from analysis"},
			{Name: "--extend-exclude", Description: "Respect file exclusions via `.gitignore` and other standard ignore files"},
			{Name: "--force-exclude", Description: "Enforce exclusions, even for paths passed to Ruff directly on the command-line"},
			{Name: "-n", Description: "Disable cache reads"},
			{Name: "--cache-dir", Description: "Path to the cache directory"},
			{Name: "--stdin-filename", Description: "The name of the file when passing it through stdin"},
			{Name: "-e", Description: "The minimum Python version that should be supported"},
			{Name: "--respect-gitignore", Description: "Respect file exclusions via `.gitignore` and other standard ignore files"},
			{Name: "--line-length", Description: "Set the line-length"},
			{Name: "--range", Description: "When specified, Ruff will try to only format the code in the given range"},
			{Name: "--all", Description: "Explain all rules"},
			{Name: "--no-preview", Description: "Disable preview mode"},
			{Name: "--direction", Description: "Attempt to detect imports from string literals"},
		},
	})
}
