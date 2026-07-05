package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "networkQuality",
		Description: "Measure the different aspects of network quality",
		Options: []core.Option{
			{Name: "-h", Description: "Show help for networkQuality"},
			{Name: "-c", Description: "Produce computer readable output"},
			{Name: "-s", Description: "Run tests sequentially instead of in parallel"},
			{Name: "-v", Description: "Verbose output"},
			{Name: "-C", Description: "Use a custom configuration URL"},
			{Name: "-I", Description: "Bind test to interface"},
		},
	})
}
