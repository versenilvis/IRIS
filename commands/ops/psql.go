package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "psql",
		Description: "Psql is a terminal-based front-end to PostgreSQL",
		Options: []spec.Option{
			{Name: "-a", Description: "Put all query output into file filename. This is equivalent to the command \\\\o"},
			{Name: "-p", Description: "Print the psql version and exit"},
			{Name: "-w", Description: "Show help about psql and exit"},
		},
	})
}
