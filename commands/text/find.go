package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "find",
		Description: "search for files",
		Generator:   spec.FileGenerator("/"),
		Options: []spec.Option{
			// BSD find; these came from a second `find` spec in commands/fs that
			// silently overwrote or was overwritten by this one depending on
			// package init order.
			{Name: "-E", Description: "Permit find to be safely used in conjunction with xargs"},
			{Name: "-d", Description: "Cause find to perform a depth-first traversal"},
			{Name: "-f", Description: "Specify a file hierarch for find to traverse"},
			{Name: "-s", Description: "Cause find to traverse the file hierarchies in lexicographical order"},
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
