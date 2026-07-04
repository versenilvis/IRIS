package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "blitz",
		Description: "Show help for command",
		Subcommands: []core.Subcommand{
			{Name: "build", Description: "Creates a production build"},
			{Name: "codegen", Description: "Generates Routes Manifest"},
			{Name: "console", Description: "Run the Blitz console REPL"},
			{Name: "db", Description: "Run database commands"},
			{Name: "command", Description: "Run specific db command"},
			{Name: "dev", Description: "Start a development server"},
			{Name: "export", Description: "Exports a static page"},
			{Name: "generate", Description: "Generate new files for your Blitz project"},
			{Name: "type", Description: "What files to generate"},
			{Name: "help", Description: "Display help for <%= config.bin %>"},
			{Name: "install", Description: "Install a Recipe into your Blitz app"},
			{Name: "new", Description: "Create a new Blitz project"},
			{Name: "name", Description: "Name of your new project"},
			{Name: "prisma", Description: "Loads env variables then proxies all args to Prisma CLI"},
			{Name: "routes", Description: "Display all Blitz URL Routes"},
			{Name: "start", Description: "Start the production server"},
			{Name: "autocomplete", Description: "Display autocomplete installation instructions"},
			{Name: "shell", Description: "Shell type"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for command"},
			{Name: "-p", Description: "Set port number"},
			{Name: "-H", Description: "Set server hostname"},
			{Name: "--inspect", Description: "Enable the Node.js inspector"},
			{Name: "--no-incremental-build", Description: "Disable incremental build and start from a fresh cache"},
			{Name: "-o", Description: "Set the output dir (defaults to 'out')"},
			{Name: "-c", Description: "Show what files will be created without writing them to disk"},
			{Name: "--all", Description: "See all commands in CLI"},
			{Name: "--npm", Description: "Use npm as the package manager"},
			{Name: "--yarn", Description: "Use yarn as the package manager"},
			{Name: "--form", Description: "A form library"},
			{Name: "-d", Description: "Show what files will be created without writing them to disk"},
			{Name: "--no-git", Description: "Skip git repository creation"},
			{Name: "--skip-upgrade", Description: "Skip blitz upgrade if outdated"},
			{Name: "-r", Description: "Refresh cache (ignores displaying instructions)"},
		},
	})
}
