package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pmset",
		Description: "Display sleep timer (value in minutes, or 0 to disable)",
		Subcommands: []spec.Subcommand{
			{Name: "live", Description: "Display the settings currently in use. (default if no argument given)"},
			{Name: "custom", Description: "Display custom settings for all power sources"},
			{Name: "cap", Description: "Display which power management features the machine supports"},
			{Name: "sched", Description: "Display scheduled startup/wake and shutdown/sleep events"},
			{Name: "ups", Description: "Display UPS emergency thresholds"},
			{Name: "ps", Description: "Display status of batteries and UPSs"},
			{Name: "pslog", Description: "Display an ongoing log of power source (battery and UPS)state"},
			{Name: "rawlog", Description: "Display an ongoing log of battery state as read directly from battery"},
			{Name: "profiles", Description: "Display the settings associated with each Energy Saver profile. 10.5+"},
			{Name: "assertions", Description: "Display a summary of power assertions. 10.6+"},
			{Name: "assertionslog", Description: "Show a log of assertion creations and releases. 10.6+"},
			{Name: "sysloadlog", Description: "Display an ongoing log of lives changes to the system load advisory. 10.6+"},
			{Name: "ac", Description: "Display details about an attached AC power adapter"},
			{Name: "log", Description: "Display a history of sleeps, wakes, and other power management events"},
			{Name: "uuid", Description: "Display the currently active sleep/wake UUID"},
			{Name: "history", Description: "A debugging tool. Prints a timeline of system sleeplwake UUIDs"},
			{Name: "historydetailed", Description: "Prints driver-level timings for a sleep/wake. Pass a UUID as an argument"},
			{Name: "everything", Description: "Print output from every argument under the GETTING header 10.8+"},
			{Name: "UUID", Description: "Used for historydetailed"},
			{Name: "schedule", Description: "For setting up one-time power events"},
			{Name: "repeat", Description: "For setting up daily/weekly power on and power off events"},
		},
		Options: []spec.Option{
			{Name: "-g", Description: "GETTING"},
			{Name: "-a", Description: "Settings for all"},
			{Name: "-b", Description: "Settings for battery"},
			{Name: "-c", Description: "Settings for charger"},
			{Name: "-u", Description: "Settings for UPS"},
		},
	})
}
