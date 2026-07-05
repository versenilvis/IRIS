package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "codesign",
		Description: "Create and manipulate code signatures",
		Options: []core.Option{
			{Name: "--all-architectures", Description: "Display information about the code at the path(s) given"},
			{Name: "-D", Description: "Constructs and prints the hosting chain of a running program"},
			{Name: "-i", Description: "Indicates the granularity of code signing. Pagesize must be a power of two"},
			{Name: "-r", Description: "Sign the code at the path(s) given using this identity"},
			{Name: "-v", Description: "Requests verification of code signatures"},
			{Name: "--continue", Description: "During static validation, do not validate the contents of the code's resources"},
		},
	})
}
