package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "date",
		Description: "Display or set date and time",
		Options: []core.Option{
			{Name: "-d", Description: "Set the kernel's value for daylight saving time"},
			{Name: "-f", Description: "The format with which to parse the new date value"},
			{Name: "-j", Description: "Don't try to set the date"},
			{Name: "-n", Description: "Only set time on the current machine, instead of all machines in the local group"},
			{Name: "-R", Description: "Use RFC 2822 date and time output format"},
			{Name: "-r", Description: "Number of seconds since the Epoch (00:00:00 UTC, January 1, 1970)"},
			{Name: "-t", Description: "Set the system's value for minutes west of GMT"},
			{Name: "-u", Description: "Display or set the date in UTC (Coordinated Universal) time"},
			{Name: "-v", Description: "[+|-]val[ymwdHMS]"},
		},
	})
}
