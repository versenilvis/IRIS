package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "truncate",
		Description: "Shrink or extend the size of a file to the specified size",
		Options: []core.Option{
			{Name: "--no-create", Description: "Do not create any files"},
			{Name: "--io-blocks", Description: "Treat SIZE as number of IO blocks instead of bytes"},
			{Name: "--reference", Description: "Base size on RFILE"},
			{Name: "--size", Description: "Set or adjust the file size by SIZE bytes"},
			{Name: "--help", Description: "Show help for truncate"},
			{Name: "--version", Description: "Output version information and exit"},
		},
	})
}
