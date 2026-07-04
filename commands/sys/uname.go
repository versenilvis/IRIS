package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "uname",
		Description: "Print operating system name",
		Options: []core.Option{
			{Name: "-a", Description: "Print all available system information"},
			{Name: "-m", Description: "Print the machine hardware name"},
			{Name: "-n", Description: "Print the system hostname"},
			{Name: "-p", Description: "Print the machine processor architecture name"},
			{Name: "-r", Description: "Print the operating system release"},
			{Name: "-s", Description: "Print the operating system name"},
			{Name: "-v", Description: "Print the operating system version"},
		},
	})
}
