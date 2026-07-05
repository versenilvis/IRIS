package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "cal",
		Description: "Displays a calendar and the date of Easter",
		Options: []spec.Option{
			{Name: "-h", Description: "Turns off highlighting of today"},
			{Name: "-j", Description: "Display Julian days (days one-based, numbered from January 1)"},
			{Name: "-m", Description: "Display a calendar for the specified year"},
		},
	})
}
