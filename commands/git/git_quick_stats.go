package git

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "git-quick-stats",
		Description: "Show help for git-quick-stats",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for git-quick-stats"},
			{Name: "--suggest-reviewers", Description: "Show the best people to contact to review code"},
			{Name: "--detailed-git-stats", Description: "Give a detailed list of git stats"},
			{Name: "---git-stats-by-branch", Description: "Show detailed list of git stats by branch"},
			{Name: "--changelogs", Description: "Show changelogs"},
			{Name: "--changelogs-by-author", Description: "Show changelogs by author"},
			{Name: "--my-daily-stats", Description: "Show your current daily stats"},
			{Name: "--json-output", Description: "Save git log as a JSON formatted file to a specified area"},
			{Name: "--branch-tree", Description: "Show an ASCII graph of the git repo branch history"},
			{Name: "--branches-by-date", Description: "Show branches by date"},
			{Name: "--contributors", Description: "See a list of everyone who contributed to the repo"},
			{Name: "--commits-per-author", Description: "Displays a list of commits per author"},
			{Name: "--commits-per-day", Description: "Displays a list of commits per day"},
			{Name: "--commits-by-month", Description: "Displays a list of commits per month"},
			{Name: "--commits-by-year", Description: "Displays a list of commits per year"},
			{Name: "--commits-by-weekday", Description: "Displays a list of commits per weekday"},
			{Name: "--commits-by-hour", Description: "Displays a list of commits per hour"},
			{Name: "--commits-by-author-by-hour", Description: "Displays a list of commits per hour by author"},
			{Name: "--commits-by-timezone", Description: "Displays a list of commits per timezone"},
		},
	})
}
