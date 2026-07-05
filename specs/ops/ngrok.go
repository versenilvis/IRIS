package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ngrok",
		Description: "Path to log file,",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Shows a list of commands or help for one command"},
			{Name: "authtoken", Description: "Save authtoken to configuration file"},
			{Name: "credits", Description: "Prints author and licensing information"},
			{Name: "http", Description: "Start an HTTP tunnel"},
			{Name: "start", Description: "Start tunnels by name from the configuration file"},
			{Name: "tcp", Description: "Start a TCP tunnel"},
			{Name: "tls", Description: "Start a TLS tunnel"},
			{Name: "update", Description: "Update ngrok to the latest version"},
			{Name: "version", Description: "Print the version string"},
			{Name: "8080", Description: "Port"},
		},
		Options: []core.Option{
			{Name: "--all", Description: "Start all tunnels in the configuration file"},
			{Name: "--none", Description: "Start running no tunnels"},
			{Name: "--remote-addr", Description: "Bind remote address (requires you reserve an address)"},
			{Name: "--client-cas", Description: "Update ngrok to the latest version"},
			{Name: "--channel", Description: "Update channel (stable, beta)"},
		},
	})
}
