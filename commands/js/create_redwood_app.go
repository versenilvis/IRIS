package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-redwood-app",
		Description: "Name of your Redwood project",
		Options: []core.Option{
			{Name: "--help", Description: "Show help"},
			{Name: "--typescript", Description: "Generate a TypeScript project"},
			{Name: "--overwrite", Description: "Create even if target directory isn't empty"},
			{Name: "--telemetry", Description: "Enables sending telemetry events for this create"},
			{Name: "--git-init", Description: "Initialize a git repository"},
			{Name: "--commit-message", Description: "Commit message for the initial commit"},
			{Name: "--yes", Description: "Skip prompts and use defaults"},
			{Name: "--version", Description: "Show version number"},
			{Name: "--yarn-install", Description: "Install node modules. Skip via --no-yarn-install"},
		},
	})
}
