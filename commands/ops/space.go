package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "space",
		Description: "Deta Space CLI for mananging Deta Space projects",
		Subcommands: []core.Subcommand{
			{Name: "completion", Description: "Generate the autocompletion script for the specified shell"},
			{Name: "bash", Description: "Generate the autocompletion script for bash"},
			{Name: "fish", Description: "Generate the autocompletion script for fish"},
			{Name: "powershell", Description: "Generate the autocompletion script for powershell"},
			{Name: "zsh", Description: "Generate the autocompletion script for zsh"},
			{Name: "link", Description: "Link code to project"},
			{Name: "login", Description: "Login to space"},
			{Name: "new", Description: "Create new project"},
			{Name: "open", Description: "Open current project in browser"},
			{Name: "push", Description: "Push code for project"},
			{Name: "release", Description: "Create release for a project"},
			{Name: "validate", Description: "Validate spacefile in dir"},
			{Name: "version", Description: "Space CLI version"},
			{Name: "upgrade", Description: "Upgrade Space CLI version"},
			{Name: "help", Description: "Help about any command"},
		},
		Options: []core.Option{
			{Name: "--no-descriptions", Description: "Disable completion descriptions"},
			{Name: "--dir", Description: "Src of project to link"},
			{Name: "--id", Description: "Project id of project to link"},
			{Name: "--blank", Description: "Create blank project"},
			{Name: "--name", Description: "Project name"},
			{Name: "--open", Description: "Open builder instance/project in browser after push"},
			{Name: "--skip-logs", Description: "Skip following logs after push"},
			{Name: "--tag", Description: "Tag to identify this push"},
			{Name: "--confirm", Description: "Release latest revision"},
			{Name: "--listed", Description: "Listed on discovery"},
			{Name: "--notes", Description: "Release notes"},
			{Name: "--rid", Description: "Revision id for release"},
			{Name: "--version", Description: "Version for the release"},
			{Name: "--help", Description: "Display help"},
		},
	})
}
