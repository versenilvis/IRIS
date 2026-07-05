package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "fdisk",
		Description: "Manipulate disk partition table",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for fdisk"},
			{Name: "--version", Description: "Show version for lsblk"},
			{Name: "--sector-size", Description: "Specify the sector size of the disk"},
			{Name: "--protect-boot", Description: "Specify the compatibility mode, 'dos' or 'nondos'"},
			{Name: "--color", Description: "Colorize the output"},
			{Name: "--list", Description: "List the partition tables for the specified devices and then exit"},
			{Name: "--list-details", Description: "Like --list, but provides more details"},
			{Name: "--lock", Description: "Use exclusive BSD lock for device or file it operates"},
			{Name: "--noauto-pt", Description: "Don't automatically create a default partition table on empty device"},
			{Name: "--output", Description: "Desc"},
			{Name: "--getsz", Description: "This option is DEPRECATED in favour of blockdev(8)"},
			{Name: "--type", Description: "When listing partition tables, show sizes in 'sectors' or in 'cylinders'"},
			{Name: "--cylinders", Description: "Specify the number of cylinders of the disk"},
			{Name: "--heads", Description: "The argument when can be auto, never or always"},
			{Name: "--wipe-partitions", Description: "The argument when can be auto, never or always"},
		},
	})
}
