package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "prisma",
		Description: "Display this help message",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Setup Prisma for your app"},
			{Name: "generate", Description: "Generate artifacts (e.g. Prisma Client)"},
			{Name: "studio", Description: "Open Prisma Studio"},
			{Name: "format", Description: "Format your schema"},
			{Name: "migrate", Description: "Migrate your database"},
			{Name: "resolve", Description: "Resolve issues with database migrations in deployment databases"},
			{Name: "status", Description: "Check the status of your database migrations"},
			{Name: "db", Description: "Manage your database schema and lifecycle (Preview)"},
			{Name: "version", Description: "Print current version of Prisma components"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Display this help message"},
			{Name: "--schema", Description: "Custom path to your Prisma schema"},
			{Name: "--skip-seed", Description: "Skip triggering seed"},
			{Name: "--skip-generate", Description: "Skip triggering generators (e.g. Prisma Client)"},
			{Name: "--datasource-provider", Description: "Define the datasource provider to use"},
			{Name: "--url", Description: "Define a custom datasource url"},
			{Name: "--data-proxy", Description: "Enable the Data Proxy in the Prisma Client"},
			{Name: "--no-hints", Description: "Hides the hint messages but still outputs errors and warnings"},
			{Name: "--no-engine", Description: "Generate a client for use with Accelerate only"},
			{Name: "--watch", Description: "Watch the Prisma schema and rerun after a change"},
			{Name: "--allow-no-models", Description: "Allow generating a client without models"},
			{Name: "-p", Description: "Port to start Studio on"},
			{Name: "-b", Description: "Browser to open Studio in"},
			{Name: "-n", Description: "Hostname to bind the Express server to"},
			{Name: "--create-only", Description: "The name of the migration. If no name is provided, the CLI will prompt you"},
			{Name: "-f", Description: "Skip the confirmation prompt"},
			{Name: "--applied", Description: "Record a specific migration as applied"},
			{Name: "--rolled-back", Description: "Record a specific migration as rolled back"},
			{Name: "--from-url", Description: "A datasource url"},
			{Name: "--to-url", Description: "A datasource url"},
			{Name: "--from-empty", Description: "Flag to assume from is an empty datamodel"},
			{Name: "--to-empty", Description: "Flag to assume to is an empty datamodel"},
			{Name: "--from-schema-datamodel", Description: "Path to a Prisma schema file, uses the 'datamodel' for the diff"},
			{Name: "--to-schema-datamodel", Description: "Path to a Prisma schema file, uses the 'datamodel' for the diff"},
			{Name: "--from-schema-datasource", Description: "Path to a Prisma schema file, uses the 'datasource url' for the diff"},
			{Name: "--to-schema-datasource", Description: "Path to a Prisma schema file, uses the 'datasource url' for the diff"},
			{Name: "--from-migrations", Description: "Path to the Prisma Migrate migrations directory"},
			{Name: "--to-migrations", Description: "Path to the Prisma Migrate migrations directory"},
			{Name: "--shadow-database-url", Description: "Manage your database schema and lifecycle (Preview)"},
			{Name: "--force", Description: "Ignore current Prisma schema file"},
			{Name: "--print", Description: "Print the introspected Prisma schema to stdout"},
			{Name: "--composite-type-depth", Description: "Skip generation of artifacts such as Prisma Client"},
			{Name: "--force-reset", Description: "Seed your database"},
			{Name: "--file", Description: "Path to a file. The content will be sent as the script to be executed"},
			{Name: "--stdin", Description: "Use the terminal standard input as the script to be executed"},
			{Name: "--json", Description: "Output JSON"},
		},
	})
}
