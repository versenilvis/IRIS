package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "whereis",
		Description: "Locate the binary, source, and manual page files for a command",
		Options: []spec.Option{
			{Name: "-b", Description: "Search only for binaries"},
			{Name: "-m", Description: "Search only for manual sections"},
			{Name: "-s", Description: "Search only for sources"},
			{Name: "-u", Description: "Search for unusual entries"},
			{Name: "-B", Description: "Search for binaries only in the specified directory"},
			{Name: "-M", Description: "Search for manual pages only in the specified directory"},
			{Name: "-S", Description: "Search for sources only in the specified directory"},
			{Name: "-f", Description: "Terminate the -B, -M, and -S options"},
		},
	})
}
