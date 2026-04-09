package commands

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cargo",
		Description: "rust toolchain",
		Subcommands: []core.Subcommand{
			{Name: "build", Description: "compile project",
				Options: []core.Option{
					{Name: "--release", Description: "optimized build"},
				},
			},
			{Name: "run", Description: "run project"},
			{Name: "test", Description: "run tests"},
			{Name: "new", Description: "create project"},
			{Name: "add", Description: "add dependency"},
			{Name: "fmt", Description: "format code"},
			{Name: "clippy", Description: "lint code"},
			{Name: "check", Description: "check without build"},
			{Name: "clean", Description: "remove build artifacts"},
			{Name: "publish", Description: "publish crate"},
			{Name: "doc", Description: "build documentation"},
		},
	})
}
