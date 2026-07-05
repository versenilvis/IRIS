package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pocketbase",
		Description: "PocketBase CLI",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "Help about any command"},
			{Name: "migrate", Description: "Executes DB migration scripts"},
			{Name: "folder", Description: "Migrations folder"},
			{Name: "create", Description: "Creates new migration template file"},
			{Name: "name", Description: "Migration file name"},
			{Name: "down", Description: "Reverts the last [number] applied migrations"},
			{Name: "number", Description: "Number of applied migrations to revert"},
			{Name: "up", Description: "Runs all available migrations"},
			{Name: "serve", Description: "Starts the web server (default to 127.0.0.1:8090)"},
			{Name: "string", Description: "API HTTP server address"},
			{Name: "strings", Description: "CORS allowed domain origins list"},
		},
		Options: []spec.Option{
			{Name: "--http", Description: "API HTTP server address"},
			{Name: "--https", Description: "API HTTPS server address (auto TLS via Let's Encrypt)"},
			{Name: "--origins", Description: "CORS allowed domain origins list (default [*])"},
			{Name: "--debug", Description: "Enable debug mode, aka showing more detailed logs"},
			{Name: "--dir", Description: "PocketBase data directory"},
			{Name: "--encryptionEnv", Description: "Encryption environment variable name"},
			{Name: "-h", Description: "Show help for pocketbase"},
			{Name: "-v", Description: "Show version for pocketbase"},
		},
	})
}
