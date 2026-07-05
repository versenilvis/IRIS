package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "typeorm",
		Description: "Show help for command",
		Subcommands: []core.Subcommand{
			{Name: "schema:sync", Description: "Synchronizes your entities with database schema"},
			{Name: "schema:log", Description: "Shows sql to be executed by schema:sync command"},
			{Name: "schema:drop", Description: "Drops all tables in the database on your default connection"},
			{Name: "query", Description: "Executes given SQL query on a default connection"},
			{Name: "entity:create", Description: "Generates a new entity"},
			{Name: "subscriber:create", Description: "Generates a new subscriber"},
			{Name: "migration:create", Description: "Creates a new migration file"},
			{Name: "migration:generate", Description: "Generates a new migration file with sql needs to be executed to update schema"},
			{Name: "migration:run", Description: "Runs all pending migrations"},
			{Name: "migration:show", Description: "Show all migrations and whether they have been run or not"},
			{Name: "migration:revert", Description: "Reverts last executed migration"},
			{Name: "version", Description: "Prints TypeORM version this project uses"},
			{Name: "cache:clear", Description: "Clears all data stored in query runner cache"},
			{Name: "init", Description: "Generates initial TypeORM project structure"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for command"},
			{Name: "-v", Description: "Show the version"},
			{Name: "-c", Description: "Name of the connection on which to run a query"},
			{Name: "-f", Description: "Name of the file with connection configuration"},
			{Name: "-n", Description: "Name of the entity class"},
			{Name: "-d", Description: "Directory where entity should be created"},
			{Name: "-o", Description: "Generate a migration file on Javascript instead of Typescript"},
			{Name: "-p", Description: "Pretty-print generated SQL"},
			{Name: "--dr", Description: "Prints out the contents of the migration instead of writing it to a file"},
			{Name: "--ch", Description: "Runs all pending migrations"},
			{Name: "-t", Description: "Indicates if transaction should be used or not for migration run"},
			{Name: "--db", Description: "Database type you'll use in your project"},
			{Name: "--express", Description: "Indicates if express should be included in the project"},
			{Name: "--docker", Description: "Set to true if docker-compose must be generated as well"},
			{Name: "--pm", Description: "Install packages"},
		},
	})
}
