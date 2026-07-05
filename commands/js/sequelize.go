package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "sequelize",
		Description: "The environment to run the command in",
		Options: []spec.Option{
			{Name: "--env", Description: "The environment to run the command in"},
			{Name: "--config", Description: "The path to the config file"},
			{Name: "--options-path", Description: "The path to a JSON file with additional options"},
			{Name: "--migrations-path", Description: "The path to the migrations folder"},
			{Name: "--seeders-path", Description: "The path to the seeders folder"},
			{Name: "--models-path", Description: "The path to the models folder"},
			{Name: "--url", Description: "The database connection string to use"},
			{Name: "--debug", Description: "When available show various debug information"},
			{Name: "--help", Description: "Show help"},
			{Name: "--version", Description: "Show version number"},
			{Name: "--charset", Description: "Pass charset option to dialect, MYSQL only"},
			{Name: "--collate", Description: "Pass collate option to dialect"},
			{Name: "--encoding", Description: "Pass encoding option to dialect, PostgreSQL only"},
			{Name: "--ctype", Description: "Pass ctype option to dialect, PostgreSQL only"},
			{Name: "--template", Description: "Pass template option to dialect, PostgreSQL only"},
			{Name: "--force", Description: "Will drop the existing config folder and re-create it"},
			{Name: "--to", Description: "Migration name to run migrations until"},
			{Name: "--from", Description: "Migration name to start migrations from (excluding)"},
			{Name: "--seed", Description: "List of seed files"},
			{Name: "--name", Description: "Name of the migration to undo"},
			{Name: "--attributes", Description: "A list of attributes"},
			{Name: "--underscored", Description: "Use snake case for the timestamp's attribute names"},
		},
	})
}
