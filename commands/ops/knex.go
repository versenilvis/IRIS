package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "knex",
		Description: "SQL query builder for JavaScript",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Create a fresh knexfile"},
			{Name: "migrate:make", Description: "Create a named migration file"},
			{Name: "migrate:latest", Description: "Run all migrations that have not yet been run"},
			{Name: "migrate:up", Description: "Run the next or the specified migration that has not yet been run"},
			{Name: "migrate:rollback", Description: "Rollback the last batch of migrations performed"},
			{Name: "migrate:down", Description: "Undo the last or the specified migration that was already run"},
			{Name: "migrate:currentVersion", Description: "View the current version for the  migration"},
			{Name: "migrate:list|migrate:status", Description: "List all migrations files with status"},
			{Name: "migrate:unlock", Description: "Forcibly unlocks the migrations lock table"},
			{Name: "seed:make", Description: "Create a named seed file"},
			{Name: "seed:run", Description: "Run seed files"},
			{Name: "help", Description: "Display help for command"},
		},
		Options: []spec.Option{
			{Name: "--version", Description: "Output the version number"},
			{Name: "--debug", Description: "Run with debugging"},
			{Name: "--knexfile", Description: "Specify the knexfile path"},
			{Name: "--knexpath", Description: "Specify the path to knex instance"},
			{Name: "--cwd", Description: "Specify the working directory"},
			{Name: "--client", Description: "Set DB client without a knexfile"},
			{Name: "--connection", Description: "Set DB connection without a knexfile"},
			{Name: "--migrations-directory", Description: "Set migrations directory without a knexfile"},
			{Name: "--migrations-table-name", Description: "Set migrations table name without a knexfile"},
			{Name: "--env", Description: "Environment, default: process.env.NODE_ENV || development"},
			{Name: "--esm", Description: "Enable ESM interop"},
			{Name: "--specific", Description: "Specify one seed file to execute"},
			{Name: "--timestamp-filename-prefix", Description: "Enable a timestamp prefix on name of generated seed files"},
			{Name: "--help", Description: "Display help for command"},
		},
	})
}
