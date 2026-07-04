package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "hexo",
		Description: "Draft for",
		Subcommands: []core.Subcommand{
			{Name: "config", Description: "Get or set configurations"},
			{Name: "help", Description: "Displays a help for each command"},
			{Name: "command", Description: "The command to display help for"},
			{Name: "init", Description: "Create a new Hexo folder"},
			{Name: "destination", Description: "Folder path. Initialize in current folder if not specified"},
			{Name: "new", Description: "Create a new article"},
			{Name: "path", Description: "The path of the post"},
			{Name: "slug", Description: "The slug of the post"},
			{Name: "layout", Description: "The layout to use"},
			{Name: "title", Description: "The title of the post"},
			{Name: "generate", Description: "Generate static files"},
			{Name: "concurrency", Description: "The number of files to generate in parallel"},
			{Name: "publish", Description: "Publish a draft"},
			{Name: "filename", Description: "The name of the post"},
			{Name: "server", Description: "Start a local server. By default, this is at http://localhost:4000/"},
			{Name: "ip", Description: "The IP address to bind to"},
			{Name: "port", Description: "Only serve static files"},
			{Name: "deploy", Description: "Deploy your website"},
			{Name: "output directory", Description: "The path to output directory"},
			{Name: "file", Description: "The file to render"},
			{Name: "migrate", Description: "Migrate content from other blog systems"},
			{Name: "type", Description: "The type of migration. check https://hexo.io/docs/migration for more details"},
			{Name: "clean", Description: "Clean the cache file (`db.json`) and generated files (`public`)"},
			{Name: "list", Description: "List all routes"},
			{Name: "version", Description: "Display version information"},
		},
		Options: []core.Option{
			{Name: "--no-clone", Description: "Copy files instead of cloning from GitHub"},
			{Name: "--no-install", Description: "Skip npm install"},
			{Name: "-p", Description: "Post path. Customize the path of the post"},
			{Name: "-r", Description: "Replace the current post if it existed"},
			{Name: "-s", Description: "Post slug. Customize the URL of the post"},
			{Name: "-d", Description: "Deploy after generation finishes"},
			{Name: "-f", Description: "Force regenerate"},
			{Name: "-w", Description: "Watch file changes"},
			{Name: "-b", Description: "Raise an error if any unhandled exception is thrown during generation"},
			{Name: "-c", Description: "Maximum number of files to be generated in parallel. Default is infinity"},
			{Name: "-i", Description: "Override the default server IP. Bind to all IP address by default"},
			{Name: "-l", Description: "Enable logger. Override logger format"},
			{Name: "-o", Description: "Immediately open the server url in your default web browser"},
			{Name: "--setup", Description: "Setup without deployment"},
			{Name: "-g", Description: "Generate static files before deploying"},
			{Name: "--engine", Description: "Specify render engine"},
			{Name: "--pretty", Description: "Prettify JSON output"},
			{Name: "--config", Description: "Disable loading plugins and scripts"},
			{Name: "--debug", Description: "Log verbose messages to the terminal and to `debug.log`"},
			{Name: "--silent", Description: "Silence output to the terminal"},
			{Name: "--draft", Description: "Display draft posts (stored in the `source/_drafts` folder)"},
			{Name: "--cwd", Description: "Customize the path of current working directory"},
		},
	})
}
