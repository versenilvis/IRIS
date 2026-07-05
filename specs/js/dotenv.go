package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dotenv",
		Description: "Loads environment variables from .env",
		Options: []core.Option{
			{Name: "-f", Description: "List of env files to parse"},
			{Name: "-h", Description: "Display help"},
			{Name: "-v", Description: "Show version"},
			{Name: "-t", Description: "Create a template env file"},
		},
	})
}
