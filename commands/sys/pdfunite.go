package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pdfunite",
		Description: "Combine multiple pdfs",
		Options: []spec.Option{
			{Name: "-v", Description: "Print copyright and version info"},
			{Name: "-h", Description: "Print usage information"},
		},
	})
}
