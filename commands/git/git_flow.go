package git

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "git-flow",
		Description: "Git extensions to provide high-level repository operations for Vincent Driessen's branching model",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Initialize a new git repo with support for the branching model"},
			{Name: "feature", Description: "List all feature branches"},
			{Name: "start", Description: "Create a new feature branch"},
			{Name: "name", Description: "The name of the new feature branch"},
			{Name: "finish", Description: "Merge a feature branch into develop"},
			{Name: "publish", Description: "Push a feature branch to the remote repository"},
			{Name: "pull", Description: "Pull a feature branch from the remote repository"},
			{Name: "origin", Description: "The name of the remote feature branch"},
			{Name: "release", Description: "List all release branches"},
			{Name: "hotfix", Description: "List all hotfix branches"},
			{Name: "support", Description: "List all support branches"},
		},
		Options: []spec.Option{
			{Name: "-d", Description: "Use default branch naming conventions"},
			{Name: "-f", Description: "Force setting of gitflow branches, even if already configured"},
		},
	})
}
