package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "df",
		Description: "Display free disk space",
		Options: []spec.Option{
			{Name: "-a", Description: "Show all mount points"},
			{Name: "-b", Description: "Use 512-byte blocks (default)"},
			{Name: "-g", Description: "Use 1073741824-byte (1-Gbyte) blocks"},
			{Name: "-m", Description: "Use 1048576-byte (1-Mbyte) blocks"},
			{Name: "-k", Description: "Use 1024-byte (1-Kbyte) blocks"},
			{Name: "-H", Description: "Include the number of free inodes"},
			{Name: "-l", Description: "Only display information about locally-mounted filesystems"},
			{Name: "-n", Description: "Print out the previously obtained statistics"},
		},
	})
}
