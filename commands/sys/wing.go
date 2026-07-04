package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "wing",
		Description: "Runs a Wing executable in the Wing Console",
		Subcommands: []core.Subcommand{
			{Name: "run", Description: "Runs a Wing executable in the Wing Console"},
			{Name: "executable", Description: "Executable .wx file"},
			{Name: "compile", Description: "Compiles a Wing program"},
			{Name: "entrypoint", Description: "Program .w entrypoint"},
			{Name: "upgrade", Description: "Upgrades the Wing toolchain to the latest version"},
			{Name: "help", Description: "Display help for command"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Display help for command"},
			{Name: "-o", Description: "Output directory"},
			{Name: "-t", Description: "Target platform (options: 'tf-aws', 'sim')"},
		},
	})
}
