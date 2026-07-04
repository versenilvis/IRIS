package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "speedtest-cli",
		Description: "Command line interface for testing internet bandwidth using speedtest.net",
		Options: []core.Option{
			{Name: "--help", Description: "Show this help message and exit"},
			{Name: "--no-download", Description: "Do not perform download test"},
			{Name: "--no-upload", Description: "Do not perform upload test"},
			{Name: "--single", Description: "Suppress verbose output, only show basic information"},
			{Name: "--csv", Description: "Single character delimiter to use in CSV ouput. Default ','"},
			{Name: "--csv-header", Description: "Print CSV headers"},
			{Name: "--json", Description: "Display a list of speedtest.net servers sorted by distance"},
			{Name: "--server", Description: "Specify a server ID to test against. Can be supplied multiple times"},
			{Name: "--exclude", Description: "Exclude a server from selection. Can be supplied multiple times"},
			{Name: "--mini", Description: "URL for the Speedtest Mini server"},
			{Name: "--source", Description: "Source IP address to bind to"},
			{Name: "--timeout", Description: "HTTP timeout in seconds. Default 10"},
			{Name: "--secure", Description: "Use HTTPS instead of HTTP when communicating with speedtest.net operated servers"},
			{Name: "--no-pre-allocate", Description: "Show the version number and exit"},
		},
	})
}
