package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "composer",
		Description: "Composer Command",
		Options: []spec.Option{
			{Name: "-h", Description: "Display this help message"},
			{Name: "-q", Description: "Do not output any message"},
			{Name: "-V", Description: "Display this application version"},
			{Name: "--ansi", Description: "Force ANSI output"},
			{Name: "--no-ansi", Description: "Disable ANSI output"},
			{Name: "-n", Description: "Do not ask any interactive question"},
			{Name: "--profile", Description: "Display timing and memory usage information"},
			{Name: "--no-plugins", Description: "Whether to disable plugins"},
			{Name: "-d", Description: "If specified, use the given directory as working directory"},
			{Name: "--no-cache", Description: "Prevent use of the cache"},
			{Name: "-v", Description: "Verbosity of messages: 1 for normal output"},
			{Name: "-vv", Description: "Verbosity of messages: 2 for more verbose output"},
			{Name: "-vvv", Description: "Verbosity of messages: 3 for debug"},
			{Name: "-o", Description: "Show only recipes that are outdated"},
			{Name: "--force", Description: "Overwrite existing files when a new version of a recipe is available"},
		},
	})
}
