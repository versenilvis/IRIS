package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "top",
		Description: "Display Linux tasks",
		Options: []spec.Option{
			{Name: "-h", Description: "Show library version and usage prompt"},
			{Name: "-b", Description: "Starts top in Batch mode"},
			{Name: "-c", Description: "Starts top with last remembered c state reversed"},
			{Name: "-i", Description: "Starts top with secure mode forced"},
			{Name: "-pid", Description: "Monitor pids"},
		},
	})
}
