package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "wrangler",
		Description: "Path to configuration file [default: wrangler.toml]",
		Subcommands: []spec.Subcommand{
			{Name: "kv:namespace", Description: "Interact with your Workers KV Namespaces"},
			{Name: "create", Description: "Create a new namespace"},
			{Name: "delete", Description: "Delete namespace"},
			{Name: "list", Description: "List all namespaces on your Cloudflare account"},
			{Name: "kv:key", Description: "Individually manage Workers KV key-value"},
			{Name: "key", Description: "Key whose value to delete"},
			{Name: "get", Description: "Get a key's value from a namespace"},
			{Name: "help", Description: "Prints this message or the help of the given subcommand(s)"},
			{Name: "put", Description: "Put a key-value pair into a namespace"},
			{Name: "kv:bulk", Description: "Interact with multiple Workers KV key-value pairs at once"},
			{Name: "route", Description: "List or delete worker routes"},
			{Name: "secret", Description: "Generate a secret that can be referenced in the worker script"},
			{Name: "init", Description: "Create a wrangler.toml for an existing project"},
			{Name: "name", Description: "The name of your worker! [default: worker]"},
			{Name: "build", Description: "Build your worker"},
			{Name: "preview", Description: "Preview your code temporarily on cloudflareworkers.com"},
			{Name: "dev", Description: "Start a local server for developing your worker"},
			{Name: "publish", Description: "Publish your worker to the orange cloud"},
			{Name: "config", Description: "Authenticate Wrangler with a Cloudflare API Token or Global API Key"},
			{Name: "subdomain", Description: "Configure your workers.dev subdomain"},
			{Name: "whoami", Description: "Retrieve your user info and test your auth config"},
			{Name: "tail", Description: "Aggregate logs from production worker"},
			{Name: "login", Description: "Authenticate Wrangler with your Cloudflare username and password"},
			{Name: "report", Description: "Report an error caught by Wrangler to Cloudflare"},
		},
		Options: []spec.Option{
			{Name: "-c", Description: "Path to configuration file [default: wrangler.toml]"},
			{Name: "-e", Description: "Environment to perform a command on"},
			{Name: "-h", Description: "Prints help information"},
			{Name: "--verbose", Description: "Toggle verbose output (when applicable)"},
			{Name: "-s", Description: "The type of project you want generated"},
			{Name: "-b", Description: "The binding of the namespace this action applies to"},
			{Name: "-n", Description: "The ID of the namespace this action applies to"},
			{Name: "--headless", Description: "Don't open the browser on preview"},
			{Name: "--watch", Description: "Start a local server for developing your worker"},
			{Name: "-p", Description: "Port to listen on. Defaults to 8787"},
			{Name: "--delete-class", Description: "Delete all durable objects associated with a class in your script"},
			{Name: "--new-class", Description: "Allow durable objects to be created from a class in your script"},
			{Name: "--rename-class", Description: "Rename a durable object class in your script"},
			{Name: "--transfer-class", Description: "Authenticate Wrangler with a Cloudflare API Token or Global API Key"},
			{Name: "--api-key", Description: "Do not verify provided credentials before writing out Wrangler config file"},
			{Name: "-f", Description: "Specify an output format [default: json]  [possible values: json, pretty]"},
			{Name: "--metrics", Description: "Provides endpoint for cloudflared metrics. Used to retrieve tunnel url"},
			{Name: "-V", Description: "Prints version information"},
		},
	})
}
