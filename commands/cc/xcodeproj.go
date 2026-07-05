package cc

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xcodeproj",
		Description: "Xcodeproj lets you create and modify Xcode projects",
		Options: []spec.Option{
			{Name: "--ignore", Description: "A key to ignore in the comparison. Can be specified multiple times"},
			{Name: "--format", Description: "YAML output format"},
			{Name: "--group-option", Description: "Shows the difference between two targets"},
			{Name: "--project", Description: "The Xcode project document to use"},
			{Name: "--verbose", Description: "Show more debugging information"},
			{Name: "--version", Description: "Show the version of the tool"},
			{Name: "--no-ansi", Description: "Show output without ANSI codes"},
			{Name: "--help", Description: "Show help banner of specified command"},
		},
	})
}
