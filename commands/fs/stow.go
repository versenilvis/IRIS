package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "stow",
		Description: "Manage farms of symbolic links",
		Options: []spec.Option{
			{Name: "-n", Description: "Set the target directory to 'DIR' instead of the parent of the stow directory"},
			{Name: "-v", Description: "Ignore files ending in this Perl regex"},
			{Name: "--defer", Description: "Show Stow version, and exit"},
			{Name: "-h", Description: "Show Stow command syntax, and exit"},
		},
	})
}
