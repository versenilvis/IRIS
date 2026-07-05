package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "find",
		Description: "search for files",
		Generator:   core.FileGenerator("/"),
		Options: []core.Option{
			{Name: "-name", Description: "match by name"},
			{Name: "-iname", Description: "match by name (case insensitive)"},
			{Name: "-type", Description: "match by type (f=file, d=dir, l=link)"},
			{Name: "-size", Description: "match by size (+1M, -10k)"},
			{Name: "-mtime", Description: "modified n days ago"},
			{Name: "-atime", Description: "accessed n days ago"},
			{Name: "-ctime", Description: "changed n days ago"},
			{Name: "-newer", Description: "newer than file"},
			{Name: "-maxdepth", Description: "max depth to descend"},
			{Name: "-mindepth", Description: "min depth to start"},
			{Name: "-exec", Description: "execute command on match"},
			{Name: "-execdir", Description: "exec in matched dir"},
			{Name: "-delete", Description: "delete matched files"},
			{Name: "-print", Description: "print matched paths"},
			{Name: "-print0", Description: "null-terminated output"},
			{Name: "-empty", Description: "match empty files/dirs"},
			{Name: "-perm", Description: "match by permissions"},
			{Name: "-user", Description: "match by owner"},
			{Name: "-group", Description: "match by group"},
			{Name: "-not", Description: "negate expression"},
			{Name: "-or", Description: "logical OR"},
			{Name: "-and", Description: "logical AND"},
			{Name: "-prune", Description: "exclude from search"},
			{Name: "-ls", Description: "list in ls format"},
			{Name: "-xdev", Description: "don't cross filesystems"},
		},
	})
}
