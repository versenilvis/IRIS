package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pdfunite",
		Description: "Combine multiple pdfs",
		Options: []core.Option{
			{Name: "-v", Description: "Print copyright and version info"},
			{Name: "-h", Description: "Print usage information"},
		},
	})
}
