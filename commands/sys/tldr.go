package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tldr",
		Description: "Tldr page",
		Options: []spec.Option{
			{Name: "-h", Description: "Display help for command"},
			{Name: "-s", Description: "Search all pages for the query"},
			{Name: "--linux", Description: "Show command page for Linux"},
			{Name: "--osx", Description: "Show command page for OSX"},
			{Name: "--sunos", Description: "Show command page for SunOS"},
			{Name: "-l", Description: "Show all pages for current platform"},
			{Name: "-u", Description: "Download the latest pages and generate search index"},
			{Name: "-c", Description: "Delete the entire local cache"},
			{Name: "--platform", Description: "Select platform"},
		},
	})
}
