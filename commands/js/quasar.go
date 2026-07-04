package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "quasar",
		Description: "Quasar Framework CLI",
		Subcommands: []core.Subcommand{
			{Name: "create", Description: "Create a new Quasar project"},
			{Name: "project-name", Description: "Name of the project"},
			{Name: "info", Description: "Gather environment information for issue reporting"},
			{Name: "upgrade", Description: "Upgrade Quasar Framework packages"},
			{Name: "serve", Description: "Start development server"},
		},
		Options: []core.Option{
			{Name: "-h", Description: "Output usage information"},
			{Name: "--kit", Description: "Use specific starter kit"},
			{Name: "--branch", Description: "Use specific branch of the starter kit"},
			{Name: "--clone", Description: "Use git clone"},
			{Name: "--offline", Description: "Use a cached starter kit"},
			{Name: "--install", Description: "Also perform package upgrades"},
			{Name: "--prerelease", Description: "Allow pre-release versions (alpha/beta)"},
			{Name: "--major", Description: "Allow newer major versions (breaking changes)"},
			{Name: "--port", Description: "Port to use"},
			{Name: "--hostname", Description: "Address to use"},
			{Name: "--gzip", Description: "Compress content"},
			{Name: "--silent", Description: "Suppress log message"},
			{Name: "--colors", Description: "Log messages with colors"},
			{Name: "--open", Description: "Open browser window after starting"},
			{Name: "--cache", Description: "Cache time (max-age) in seconds"},
			{Name: "--micro", Description: "Use micro-cache"},
			{Name: "--history", Description: "Use history api fallback"},
			{Name: "--index", Description: "History mode (only!) index url path"},
			{Name: "--https", Description: "Enable HTTPS"},
			{Name: "--cert", Description: "Path to SSL cert file (Optional)"},
			{Name: "--key", Description: "Path to SSL key file (Optional)"},
			{Name: "--proxy", Description: "Proxy specific requests defined in file"},
			{Name: "--cors", Description: "Enable CORS for all requests"},
		},
	})
}
