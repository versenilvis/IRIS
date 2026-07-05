package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ls",
		Description: "list directory contents",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-a", Description: "all files"},
			{Name: "-l", Description: "long format"},
			{Name: "-h", Description: "human readable"},
			{Name: "-la", Description: "long + all"},
			{Name: "-R", Description: "recursive"},
			{Name: "-t", Description: "sort by time"},
			{Name: "-S", Description: "sort by size"},
			{Name: "-1", Description: "one per line"},
			{Name: "--color", Description: "colorize output"},
			{Name: "-d", Description: "list directories only"},
			{Name: "-i", Description: "show inode numbers"},
			{Name: "-s", Description: "show block size"},
			{Name: "-r", Description: "reverse order"},
		},
	})
}
