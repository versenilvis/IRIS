package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "hostname",
		Description: "Set or print name of current host system",
		Options: []spec.Option{
			{Name: "-f", Description: "Include domain information in the printed name"},
			{Name: "-s", Description: "Trim off any domain information from the printed name"},
			{Name: "-d", Description: "Only print domain information"},
		},
	})
}
