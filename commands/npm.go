package commands

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "npm",
		Description: "node packages",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages"},
			{Name: "run", Description: "run script",
				Subcommands: []core.Subcommand{
					{Name: "dev", Description: "development server"},
					{Name: "build", Description: "build for production"},
					{Name: "start", Description: "start application"},
					{Name: "test", Description: "run test suite"},
					{Name: "lint", Description: "run linter"},
				},
			},
			{Name: "test", Description: "run tests"},
			{Name: "init", Description: "create package.json"},
			{Name: "publish", Description: "publish package"},
			{Name: "update", Description: "update packages"},
			{Name: "uninstall", Description: "remove package"},
			{Name: "ls", Description: "list installed"},
			{Name: "audit", Description: "security audit"},
			{Name: "cache", Description: "manage cache",
				Subcommands: []core.Subcommand{
					{Name: "clean", Description: "clear cache"},
					{Name: "verify", Description: "verify cache"},
				},
			},
		},
	})
}
