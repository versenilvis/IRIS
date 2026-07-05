package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "install",
		Description: "Use suffix as the backup suffix if -b is given",
		Options: []core.Option{
			{Name: "-B", Description: "Use suffix as the backup suffix if -b is given"},
			{Name: "-b", Description: "Create directories.  Missing parent directories are created as required"},
			{Name: "-f", Description: "Specify a group. A numeric GID is allowed"},
			{Name: "-M", Description: "Disable all use of mmap(2)"},
			{Name: "-m", Description: "Specify an owner. A numeric UID is allowed"},
			{Name: "-p", Description: "Causes install to show when -C actually installs something"},
		},
	})
}
