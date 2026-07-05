package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "more",
		Description: "Opposite of less",
		Options: []spec.Option{
			{Name: "-d", Description: "Do not pause after any line containing a ^L (form feed)"},
			{Name: "-f", Description: "Count logical lines, rather than screen lines"},
			{Name: "-p", Description: "Instead, clear the whole screen and then display the text"},
			{Name: "-c", Description: "Squeeze multiple blank lines into one"},
			{Name: "-u", Description: "Silently ignored as backwards compatibility"},
			{Name: "-n", Description: "Specify the number of lines per screenful"},
			{Name: "--help", Description: "Display help text"},
			{Name: "-V", Description: "Display version information"},
		},
	})
}
