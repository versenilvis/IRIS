package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "herd",
		Description: "Display this application version",
		Subcommands: []core.Subcommand{
			{Name: "phpVersion", Description: "The PHP version you want to use; e.g 8.1"},
			{Name: "nodeVersion", Description: "The Node version you want to use; e.g 21"},
			{Name: "links", Description: "Display all of the registered Herd links"},
			{Name: "list", Description: "List commands"},
			{Name: "namespace", Description: "The namespace name"},
			{Name: "site", Description: "The site to get the PHP executable path for"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Display this application version"},
			{Name: "--quiet", Description: "Do not output any message"},
			{Name: "--verbose", Description: "Force ANSI output"},
			{Name: "--no-ansi", Description: "Disable ANSI output"},
			{Name: "--no-interaction", Description: "Do not ask any interactive question"},
			{Name: "--debug", Description: "Tail the completion debug log"},
			{Name: "--path", Description: "Get the URL to the current share tunnel for Expose"},
			{Name: "--format", Description: "Change the version of PHP used by Herd to serve the current working directory"},
			{Name: "--site", Description: "Change the version of Node used by the CLI for the current working directory"},
			{Name: "--secure", Description: "Link the site with a trusted TLS certificate"},
			{Name: "--isolate", Description: "Isolate the site to the PHP version specified, for example 8.3"},
			{Name: "--raw", Description: "To output raw command list"},
			{Name: "--lines", Description: "Show the log viewer UI for the given site"},
			{Name: "--sites", Description: "Get or set the loopback address used for Herd sites"},
			{Name: "--json", Description: "Get all of the paths registered with Herd"},
			{Name: "--expireIn", Description: "Secure the given domain with a trusted TLS certificate"},
			{Name: "--expiring", Description: "Limits the results to only sites expiring within the next 60 days"},
			{Name: "--days", Description: "Share the current site via an Expose tunnel"},
			{Name: "--all", Description: "Change the version of PHP used by Herd"},
		},
	})
}
