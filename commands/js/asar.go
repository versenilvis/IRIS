package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "asar",
		Description: "A simple extensive tar-like archive format with indexing",
		Subcommands: []spec.Subcommand{
			{Name: "pack", Description: "Create asar archive"},
			{Name: "directory", Description: "The directory you want to archive"},
			{Name: "output", Description: "The name of the output file"},
			{Name: "list", Description: "List files of asar archive"},
			{Name: "archive", Description: "The archive file"},
			{Name: "extract-file", Description: "Extract one file from archive"},
			{Name: "filename", Description: "The name of the file you want to extract"},
			{Name: "extract", Description: "Extract archive"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for asar"},
			{Name: "--V", Description: "Output the version number"},
		},
	})
}
