package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-nx-workspace",
		Description: "The name of the workspace",
		Options: []core.Option{
			{Name: "--name", Description: "Workspace name (e.g., org name)"},
			{Name: "--preset", Description: "Empty [an empty workspace with a layout that works best for building apps]"},
			{Name: "--appName", Description: "The name of the application created by some presets"},
			{Name: "--cli", Description: "CSS"},
			{Name: "--interactive", Description: "Enable interactive mode when using presets (boolean)"},
			{Name: "--packageManager", Description: "Package manager to use (npm, yarn, pnpm)"},
			{Name: "--nx-cloud", Description: "Use Nx Cloud (boolean)"},
			{Name: "--help", Description: "Show help for create-nx-workspace"},
		},
	})
}
