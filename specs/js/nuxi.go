package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nuxi",
		Description: "The directory of the target application",
		Options: []core.Option{
			{Name: "--help", Description: "Show help"},
			{Name: "--verbose", Description: "Log information about the installation process"},
			{Name: "--template", Description: "Specify a Git repository to use as a template"},
			{Name: "--force", Description: "Force clone to any existing directory"},
			{Name: "--prefer-offline", Description: "Try local cache first to download templates"},
			{Name: "--shell", Description: "Open shell in cloned directory (experimental)"},
			{Name: "--cwd", Description: "The current working directory of the target application"},
			{Name: "--dotenv", Description: "Point to another .env file to load, relative to the root directory"},
			{Name: "--prerender", Description: "Point to another .env file to load, relative to the root directory"},
			{Name: "--clipboard", Description: "Copy URL to clipboard"},
			{Name: "--open", Description: "Open URL in browser"},
			{Name: "--no-clear", Description: "Does not clear the console after startup"},
			{Name: "--port", Description: "Port to listen"},
			{Name: "--host", Description: "Hostname of the server"},
			{Name: "--https", Description: "Listen with https protocol with a self-signed certificate by default"},
			{Name: "--ssl-cert", Description: "Specify a certificate for https"},
			{Name: "--ssl-key", Description: "Specify the key for the https certificate"},
			{Name: "--dev", Description: "Run tests in development mode"},
			{Name: "--watch", Description: "Actively watch for changes and rerun tests"},
		},
	})
}
