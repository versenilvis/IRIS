package python

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "googler",
		Description: "Google from the command-line",
		Options: []spec.Option{
			{Name: "-h", Description: "Show this help message and exit"},
			{Name: "-s", Description: "Start at the Nth result"},
			{Name: "-n", Description: "Show N results"},
			{Name: "-N", Description: "Show results from news section"},
			{Name: "-V", Description: "Show results from videos section"},
			{Name: "-c", Description: "Country-specific search with top-level domain .TLD, e.g., 'in' for India"},
			{Name: "-l", Description: "Display in language"},
			{Name: "-g", Description: "Country-specific geolocation search with country code CC, e.g. 'in' for India"},
			{Name: "-x", Description: "Disable automatic spelling correction"},
			{Name: "--colorize", Description: "Whether to colorize output which enables color when stdout is a tty device"},
			{Name: "-C", Description: "Equivalent to --colorize=never"},
			{Name: "--colors", Description: "Set output colors"},
			{Name: "-j", Description: "Open the first result in web browser and exit"},
			{Name: "-t", Description: "American date format with slashes, e.g., 2/24/2020, 2/2020, 2020"},
			{Name: "--to", Description: "Ending date/month/year of date range"},
			{Name: "-w", Description: "Search a site using Google"},
			{Name: "-e", Description: "Exclude site from results"},
			{Name: "--unfilter", Description: "Do not omit similar results"},
			{Name: "-p", Description: "Tunnel traffic through an HTTP proxy"},
			{Name: "--notweak", Description: "Disable TCP optimizations and forced TLS 1.2"},
			{Name: "--json", Description: "Output in JSON format; implies --noprompt"},
			{Name: "--url-handler", Description: "Custom script or cli utility to open results"},
			{Name: "--show-browser-logs", Description: "Do not suppress browser output (stdout and stderr)"},
			{Name: "--np", Description: "Search and exit, do not prompt"},
			{Name: "-4", Description: "Only connect over IPv6"},
			{Name: "-u", Description: "Perform in-place self-upgrade"},
			{Name: "--include-git", Description: "When used with --upgrade, get latest git master"},
			{Name: "-v", Description: "Show program's version number and exit"},
			{Name: "-d", Description: "Enable debugging"},
		},
	})
}
