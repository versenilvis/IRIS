package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "vsce",
		Description: "The Visual Studio Code Extension Manager",
		Subcommands: []spec.Subcommand{
			{Name: "ls", Description: "Lists all the files that will be published"},
			{Name: "package", Description: "Packages an extension"},
			{Name: "publish", Description: "Publishes an extension"},
			{Name: "unpublish", Description: "Unpublishes an extension. Example extension id: microsoft.csharp"},
			{Name: "ls-publishers", Description: "List all known publishers"},
			{Name: "delete-publisher", Description: "Deletes a publisher"},
			{Name: "login", Description: "Add a publisher to the known publishers list"},
			{Name: "logout", Description: "Remove a publisher from the known publishers list"},
			{Name: "verify-pat", Description: "Verify if the Personal Access Token has publish rights for the publisher"},
			{Name: "show", Description: "Show extension metadata"},
			{Name: "search", Description: "Search extension gallery"},
			{Name: "help", Description: "Display help for command"},
		},
		Options: []spec.Option{
			{Name: "--yarn", Description: "Use yarn instead of npm (default inferred from presence of yarn.lock or .yarnrc)"},
			{Name: "--no-yarn", Description: "Use npm instead of yarn (default inferred from lack of yarn.lock or .yarnrc)"},
			{Name: "--packagedDependencies", Description: "Select packages that should be published only (includes dependencies)"},
			{Name: "--ignoreFile", Description: "Indicate alternative .vscodeignore"},
			{Name: "--no-dependencies", Description: "Disable dependency detection via npm or yarn"},
			{Name: "-o", Description: "Target architecture"},
			{Name: "-m", Description: "Commit message used when calling `npm version`"},
			{Name: "--no-git-tag-version", Description: "Do not update `package.json`. Valid only when [version] is provided"},
			{Name: "--githubBranch", Description: "Skip rewriting relative links"},
			{Name: "--baseContentUrl", Description: "Prepend all relative links in README.md with this url"},
			{Name: "--baseImagesUrl", Description: "Prepend all relative image links in README.md with this url"},
			{Name: "--no-gitHubIssueLinking", Description: "Disable automatic expansion of GitHub-style issue syntax into links"},
			{Name: "--no-gitLabIssueLinking", Description: "Disable automatic expansion of GitLab-style issue syntax into links"},
			{Name: "--pre-release", Description: "Mark this package as a pre-release"},
			{Name: "-p", Description: "Personal Access Token (defaults to VSCE_PAT environment variable)"},
			{Name: "-t", Description: "Target architecture"},
			{Name: "-i", Description: "Publish the provided VSIX packages"},
			{Name: "--noVerify", Description: "Indicate alternative .vscodeignore"},
			{Name: "-f", Description: "Forces Unpublished Extension"},
			{Name: "--json", Description: "Output data in json format (default: false)"},
			{Name: "-h", Description: "Display help for command"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
