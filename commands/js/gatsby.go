package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "gatsby",
		Description: "Set host. Defaults to localhost",
		Subcommands: []core.Subcommand{
			{Name: "develop", Description: "Start the Gatsby development server"},
			{Name: "build", Description: "Compile your application and make it ready for deployment"},
			{Name: "serve", Description: "Serve the production build of your site for testing"},
			{Name: "clean", Description: "Wipe out the cache (.cache folder) and public directories"},
			{Name: "info", Description: "Get helpful environment information"},
			{Name: "plugin", Description: "Run commands pertaining to gatsby plugins"},
			{Name: "docs", Description: "Documentation about using and creating plugins"},
			{Name: "repl", Description: "Gatsby will prompt you to type in commands and explore"},
		},
		Options: []core.Option{
			{Name: "-H", Description: "Set host. Defaults to localhost"},
			{Name: "-p", Description: "Set port. Defaults to env.PORT or 8000"},
			{Name: "-o", Description: "Open the site in your (default) browser for you"},
			{Name: "-S", Description: "Use HTTPS"},
			{Name: "--inspect", Description: "Opens a port for debugging"},
			{Name: "--prefix-paths", Description: "Build site with link paths prefixed (set pathPrefix in your config)"},
			{Name: "--no-uglify", Description: "Build site without uglifying JS bundles (for debugging)"},
			{Name: "--profile", Description: "Build site with react profiling"},
			{Name: "--open-tracing-config-file", Description: "Use Tracer configuration file"},
			{Name: "--graphql-tracing", Description: "Trace every graphql resolver, may have performance implications"},
			{Name: "--no-color", Description: "Disables colored terminal output"},
			{Name: "-C", Description: "Copy environment information to your clipboard"},
			{Name: "-v", Description: "View your current Gatsby CLI version"},
		},
	})
}
