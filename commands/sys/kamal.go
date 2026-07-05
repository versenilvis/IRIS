package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "kamal",
		Description: "Skip image build and push",
		Subcommands: []spec.Subcommand{
			{Name: "boot", Description: "Boot new accessory service on host (use NAME=all to boot all accessories)"},
			{Name: "upload", Description: "Upload accessory files to host"},
			{Name: "directories", Description: "Create accessory directories on host"},
			{Name: "start", Description: "Start existing accessory container on host"},
			{Name: "stop", Description: "Stop existing accessory container on host"},
			{Name: "restart", Description: "Restart existing accessory container on host"},
			{Name: "details", Description: "Show details about accessory on host (use NAME=all to show all accessories)"},
			{Name: "exec", Description: "Execute a custom command on servers"},
			{Name: "logs", Description: "Show log lines from accessory on host"},
			{Name: "status", Description: "Show status of accessory on host"},
			{Name: "remove", Description: "Remove accessory container, image and data directory from host"},
		},
		Options: []spec.Option{
			{Name: "-P", Description: "Skip image build and push"},
			{Name: "--since", Description: "Number of lines to show from each server"},
			{Name: "--grep", Description: "Show lines with grep match only (use this to fetch specific requests by id)"},
			{Name: "--follow", Description: "Follow log on primary server (or specific host set by --hosts)"},
			{Name: "--verbose", Description: "Detailed logging"},
			{Name: "--quiet", Description: "Minimal logging"},
			{Name: "--version", Description: "Run commands against a specific app version"},
			{Name: "--primary", Description: "Run commands only on primary host instead of all"},
			{Name: "--hosts", Description: "Run commands on these hosts instead of all (separate by comma)"},
			{Name: "--roles", Description: "Run commands on these roles instead of all (separate by comma)"},
			{Name: "--config_file", Description: "Path to config file"},
			{Name: "-d", Description: "Specify destination to use"},
			{Name: "--skip_hooks", Description: "Don't run hooks"},
			{Name: "-i", Description: "Execute command over ssh for an interactive shell (use for console/bash)"},
			{Name: "--reuse", Description: "Reuse currently running container instead of starting a new one"},
			{Name: "-y", Description: "Proceed without confirmation question"},
			{Name: "--interactive", Description: "Execute command over ssh for an interactive shell (use for console/bash)"},
			{Name: "--stop", Description: "Stop the stale containers found"},
			{Name: "--rolling", Description: "Start existing Traefik container on servers"},
			{Name: "--json", Description: "Output as JSON"},
			{Name: "--confirmed", Description: "Proceed without confirmation question"},
		},
	})
}
