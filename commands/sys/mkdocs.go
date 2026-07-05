package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mkdocs",
		Description: "Project documentation with Markdown",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Build the MkDocs documentation"},
			{Name: "gh-deploy", Description: "Deploy your documentation to GitHub Pages"},
			{Name: "new", Description: "Create a new MkDocs project"},
			{Name: "serve", Description: "Run the builtin development server"},
		},
		Options: []spec.Option{
			{Name: "-h", Description: "Provide a specific MkDocs config"},
			{Name: "-s", Description: "Enable strict mode. This will cause MkDocs to abort the build on any warnings"},
			{Name: "-t", Description: "The theme to use when building your documentation"},
			{Name: "--use-directory-urls", Description: "Use directory URLs when building pages (the default)"},
			{Name: "--no-directory-urls", Description: "Don't use directory URLs when building pages"},
			{Name: "-q", Description: "Silence warnings"},
			{Name: "-v", Description: "Enable verbose output"},
			{Name: "-c", Description: "Remove old files from the site directory before building (the default)"},
			{Name: "--dirty", Description: "Only rebuild pages that have been modified since last build"},
			{Name: "-d", Description: "The directory to output the result of the documentation build"},
			{Name: "-m", Description: "Force the push to the repository"},
			{Name: "--no-history", Description: "Replace the whole Git history with one new commit"},
			{Name: "--ignore-version", Description: "Ignore check that build is not being deployed with an older version of MkDocs"},
			{Name: "--shell", Description: "Use the shell when invoking Git"},
			{Name: "-a", Description: "IP address and port to serve documentation locally (default: localhost:8000)"},
			{Name: "--live-reload", Description: "Enable the live reloading in the development server (this is the default)"},
			{Name: "--no-reload", Description: "Disable the live reloading in the development server"},
			{Name: "--dirtyreload", Description: "A directory or file to watch for live reloading. Can be supplied multiple times"},
			{Name: "-V", Description: "Show the version and exit"},
		},
	})
}
