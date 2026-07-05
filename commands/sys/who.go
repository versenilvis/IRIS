package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "who",
		Description: "Display who is logged in",
		Subcommands: []spec.Subcommand{
			{Name: "am", Description: "Returns the invoker's real user name"},
		},
		Options: []spec.Option{
			{Name: "-a", Description: "Same as -bdlprTtu"},
			{Name: "-b", Description: "Time of last system boot"},
			{Name: "-d", Description: "Print dead processes"},
			{Name: "-H", Description: "Write column headings above the regular output"},
			{Name: "-l", Description: "Print system login processes (unsupported)"},
			{Name: "-m", Description: "Only print information about the current terminal"},
			{Name: "-p", Description: "Print active processes spawned by launchd(8) (unsupported)"},
			{Name: "-q", Description: "'Quick mode': List only names and number of users currently logged on"},
			{Name: "-r", Description: "Print the current runlevel"},
			{Name: "-s", Description: "List only the name, line and time fields (this is the default)"},
			{Name: "-T", Description: "Print last system clock change (unsupported)"},
			{Name: "-u", Description: "Print the idle time for each user and the associated process ID"},
		},
	})
}
