package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ncal",
		Description: "Displays a calendar and the date of Easter",
		Options: []spec.Option{
			{Name: "-h", Description: "Turns off highlighting of today"},
			{Name: "-J", Description: "Display date of Easter (for western churches)"},
			{Name: "-j", Description: "Display Julian days (days one-based, numbered from January 1)"},
			{Name: "-m", Description: "Display a calendar for the specified year"},
			{Name: "-o", Description: "Display date of Orthodox Easter (Greek and Russian Orthodox Churches)"},
			{Name: "-p", Description: "Print the number of the week below each week column"},
			{Name: "-3", Description: "Display the previous, current and next month surrounding today"},
			{Name: "-A", Description: "Display the number of months after the current month"},
			{Name: "-B", Description: "Display the number of months before the current month"},
			{Name: "-C", Description: "Switch to cal mode"},
			{Name: "-N", Description: "Switch to ncal mode"},
			{Name: "-d", Description: "Use yyyy-mm as the current date (for debugging of date selection)"},
			{Name: "-H", Description: "Use yyyy-mm-dd as the current date (for debugging of highlighting)"},
		},
	})
}
