package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "black",
		Description: "Version",
		Options: []core.Option{
			{Name: "--code", Description: "Format the code passed in as a string"},
			{Name: "--line-length", Description: "How many characters per line to allow"},
			{Name: "--target-version", Description: "Python versions that should be supported"},
			{Name: "--pyi", Description: "Format all input files regardless of file extension"},
			{Name: "--ipynb", Description: "Format all input files like Jupyter Notebooks regardless of file extension"},
			{Name: "--python-cell-magics", Description: "Add the given magic to the list of known python-magics"},
			{Name: "--skip-string-normalization", Description: "Don't normalize string quotes or prefixes"},
			{Name: "--skip-magic-trailing-comma", Description: "Don't use trailing commas as a reason to split lines"},
			{Name: "--preview", Description: "Enable potentially disruptive style changes"},
			{Name: "--check", Description: "Don't write the files back, just return the status"},
			{Name: "--diff", Description: "Don't write the files back, just output a diff for each file on stdout"},
			{Name: "--color", Description: "Show colored diff"},
			{Name: "--no-color", Description: "Show uncolored diff"},
			{Name: "--fast", Description: "Skip temporary sanity checks"},
			{Name: "--safe", Description: "Run temporary sanity checks"},
			{Name: "--required-version", Description: "Require a specific version of Black"},
			{Name: "--include", Description: "Additional exlusions"},
			{Name: "--force-exclude", Description: "Exlude matching files and folders even when passed explicitly"},
			{Name: "--stdin-filename", Description: "The name of the file when passing it through stdin"},
			{Name: "--workers", Description: "Number of parallel workers"},
			{Name: "--quiet", Description: "Don't emit non-error messages to stderr"},
			{Name: "--verbose", Description: "Show the version"},
			{Name: "--config", Description: "Read configuration from filepath"},
			{Name: "--help", Description: "Show usage information"},
		},
	})
}
