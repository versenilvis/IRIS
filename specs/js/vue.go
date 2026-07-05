package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vue",
		Description: "Vue cli tools",
		Subcommands: []core.Subcommand{
			{Name: "create", Description: "Create a new project powered by vue-cli-service"},
			{Name: "add", Description: "Install a plugin and invoke its generator in an already created project"},
			{Name: "invoke", Description: "Invoke the generator of a plugin in an already created project"},
			{Name: "inspect", Description: "Inspect the webpack config in a project with vue-cli-service"},
			{Name: "serve", Description: "Serve a .js or .vue file in development mode with zero config"},
			{Name: "build", Description: "Build a .js or .vue file in production mode with zero config"},
			{Name: "ui", Description: "Start and open the vue-cli ui"},
			{Name: "init", Description: "Generate a project from a remote template (legacy API, requires @vue/cli-init)"},
			{Name: "config", Description: "Inspect and modify the config"},
			{Name: "outdated", Description: "(experimental) check for outdated vue cli service / plugins"},
			{Name: "upgrade", Description: "(experimental) upgrade vue cli service / plugins"},
			{Name: "migrate", Description: "(experimental) run migrator for an already-installed cli plugin"},
			{Name: "info", Description: "Print debugging information about your environment"},
		},
		Options: []core.Option{
			{Name: "-p", Description: "Skip prompts and use saved or remote preset"},
			{Name: "-d", Description: "Skip prompts and use default preset"},
			{Name: "-i", Description: "Skip prompts and use inline JSON string as preset"},
			{Name: "-m", Description: "Use specified npm client when installing dependencies"},
			{Name: "-r", Description: "Use specified npm registry when installing dependencies (only for npm)"},
			{Name: "-g", Description: "Force git initialization with initial commit message"},
			{Name: "-n", Description: "Skip git initialization"},
			{Name: "-f", Description: "Overwrite target directory if it exists"},
			{Name: "--merge", Description: "Merge target directory if it exists"},
			{Name: "-c", Description: "Use git clone when fetching remote preset"},
			{Name: "-X", Description: "Use specified proxy when creating project"},
			{Name: "-b", Description: "Scaffold project without beginner instructions"},
			{Name: "--skipGetStarted", Description: "Output usage information"},
			{Name: "--registry", Description: "Use specified npm registry when installing dependencies (only for npm)"},
			{Name: "-h", Description: "Output usage information"},
			{Name: "--mode", Description: "Inspect a specific module rule"},
			{Name: "--plugin", Description: "Inspect a specific plugin"},
			{Name: "--rules", Description: "List all module rule names"},
			{Name: "--plugins", Description: "List all plugin names"},
			{Name: "-v", Description: "Show full function definitions in output"},
			{Name: "-o", Description: "Open browser"},
			{Name: "-t", Description: "Build target (app | lib | wc | wc-async, default: app)"},
			{Name: "-H", Description: "Host used for the UI server (default: localhost)"},
			{Name: "-D", Description: "Run in dev mode"},
			{Name: "--quiet", Description: "Don't output starting messages"},
			{Name: "--headless", Description: "Don't open browser on start and output port"},
			{Name: "--offline", Description: "Use cached template"},
			{Name: "-s", Description: "Set option value"},
			{Name: "-e", Description: "Open config with default editor"},
			{Name: "--json", Description: "Outputs JSON result only"},
			{Name: "--next", Description: "Also check for alpha / beta / rc versions when upgrading"},
			{Name: "--all", Description: "Upgrade all plugins"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
