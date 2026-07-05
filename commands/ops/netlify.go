package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "netlify",
		Description: "Print debugging information",
		Subcommands: []spec.Subcommand{
			{Name: "help", Description: "List available sub-commands"},
			{Name: "subcommand", Description: "The command to display help for"},
		},
		Options: []spec.Option{
			{Name: "--debug", Description: "Print debugging information"},
			{Name: "--httpProxy", Description: "Proxy server address to route requests through"},
			{Name: "-f", Description: "Delete without prompting (useful for CI)"},
			{Name: "--json", Description: "Output return values as JSON"},
			{Name: "-d", Description: "Data to use"},
			{Name: "--list", Description: "List out available API methods"},
			{Name: "-o", Description: "Disables any features that require network access"},
			{Name: "--context", Description: "Build context"},
			{Name: "--dry", Description: "Dry run: show instructions without running them"},
			{Name: "-s", Description: "Name of shell (bash|fish|zsh)"},
			{Name: "-a", Description: "Netlify auth token to deploy with"},
			{Name: "-b", Description: "Specify a folder to deploy"},
			{Name: "-m", Description: "A short message to include in the deploy log"},
			{Name: "-p", Description: "Deploy to production"},
			{Name: "--alias", Description: "Run build command before deploying"},
			{Name: "--prodIfUnlocked", Description: "Deploy to production if unlocked, create a draft otherwise"},
			{Name: "--skip-functions-cache", Description: "Timeout to wait for deployment to finish"},
			{Name: "-c", Description: "Command to run"},
			{Name: "-l", Description: "Start a public live session"},
			{Name: "--framework", Description: "Framework to use. Defaults to #auto which automatically detects a framework"},
			{Name: "--targetPort", Description: "Port of target app server"},
			{Name: "-H", Description: "Specifies a custom request method [default: GET]"},
			{Name: "-B", Description: "Path to the publish directory"},
			{Name: "-r", Description: "Replace all existing variables instead of merging them with the current ones"},
			{Name: "-n", Description: "Function name"},
			{Name: "-u", Description: "Pull template from URL"},
			{Name: "-q", Description: "Querystring to add to your function invocation"},
			{Name: "--identity", Description: "Simulate Netlify Identity authentication JWT"},
			{Name: "--no-identity", Description: "Affirm unauthenticated request"},
			{Name: "--force", Description: "Reinitialize CI hooks if the linked site is already configured to use CI"},
			{Name: "--gitRemoteName", Description: "Link a local repo or project folder to an existing site on Netlify"},
			{Name: "--name", Description: "Name of site to link to"},
			{Name: "--auth", Description: "Netlify auth token"},
			{Name: "--new", Description: "Login to new Netlify account"},
			{Name: "--silent", Description: "Silence CLI output"},
			{Name: "--admin", Description: "Open Netlify site"},
			{Name: "--site", Description: "Open site"},
			{Name: "--telemetry-disable", Description: "Opt out of sharing usage data"},
			{Name: "--telemetry-enable", Description: "Allow your usage to help shape development"},
		},
	})
}
