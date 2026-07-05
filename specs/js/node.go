package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "node",
		Description: "Run the node interpreter",
		Options: []core.Option{
			{Name: "-e", Description: "Evaluate script"},
			{Name: "--watch", Description: "Watch input files"},
			{Name: "--watch-path", Description: "Specify a watch directory or file"},
			{Name: "--watch-preserve-output", Description: "Disable the clearing of the console when watch mode restarts the process"},
			{Name: "--env-file", Description: "Specify a file containing environment variables"},
			{Name: "-p", Description: "Evaluate script and print result"},
			{Name: "-c", Description: "Syntax check script without executing"},
			{Name: "-v", Description: "Print Node.js version"},
			{Name: "-i", Description: "Always enter the REPL even if stdin does not appear to be a terminal"},
			{Name: "-h", Description: "Print node command line options (currently set)"},
			{Name: "--inspect", Description: "Activate inspector on host:port (default: 127.0.0.1:9229)"},
			{Name: "--preserve-symlinks", Description: "Run AdonisJS command-line"},
			{Name: "-prod", Description: "Build for production"},
			{Name: "--assets", Description: "Build frontend assets when webpack encore is installed"},
			{Name: "--no-assets", Description: "Disable building assets"},
			{Name: "--ignore-ts-errors", Description: "Ignore typescript errors and complete the build process"},
			{Name: "--tsconfig", Description: "Path to the TypeScript project configuration file"},
			{Name: "--encore-args", Description: "CLI options to pass to the encore command line"},
			{Name: "--client", Description: "Select the package manager to decide which lock file to copy to the build folder"},
			{Name: "-w", Description: "Watch for file changes and re-start the HTTP server on change"},
			{Name: "--node-args", Description: "CLI options to pass to the node command line"},
			{Name: "-f", Description: "Define a custom set of seeders files names to run"},
			{Name: "-r", Description: "Add resourceful methods to the controller class"},
			{Name: "--connection", Description: "The connection flag is used to lookup the directory for the migration file"},
			{Name: "--folder", Description: "Pre-select a migration directory"},
			{Name: "--create", Description: "Define the table name for creating a new table"},
			{Name: "--table", Description: "Define the table name for altering an existing table"},
			{Name: "-m", Description: "Generate the migration for the model"},
			{Name: "--force", Description: "Explicitly force to run migrations in production"},
			{Name: "--dry-run", Description: "Print SQL queries, instead of running the migrations"},
			{Name: "--batch", Description: "Use 0 to rollback to initial state"},
		},
	})
}
