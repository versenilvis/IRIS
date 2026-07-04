package text

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "shred",
		Description: "Overwrite a file to hide its contents, and optionally delete it",
		Options: []core.Option{
			{Name: "--force", Description: "Change permissions to allow writing if necessary"},
			{Name: "--iterations", Description: "Overwrite N times instead of the default (3)"},
			{Name: "--random-source", Description: "Get random bytes from FILE"},
			{Name: "--size", Description: "Shred this many bytes (suffixes like K, M, G accepted)"},
			{Name: "--remove", Description: "Like -u but give control on HOW to delete"},
			{Name: "--verbose", Description: "Show progress"},
			{Name: "--exact", Description: "Add a final overwrite with zeros to hide shredding"},
			{Name: "--help", Description: "Display this help and exit"},
			{Name: "--version", Description: "Output version information and exit"},
		},
	})
}
