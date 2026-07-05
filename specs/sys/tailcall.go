package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tailcall",
		Description: "TailCall CLI for managing and optimizing GraphQL configurations",
		Subcommands: []core.Subcommand{
			{Name: "check", Description: "Validate a composition spec"},
			{Name: "start", Description: "Launch the GraphQL Server for the specific configuration"},
			{Name: "init", Description: "Bootstrap a new TailCall project"},
			{Name: "gen", Description: "Generate GraphQL configurations from various sources"},
		},
		Options: []core.Option{
			{Name: "--n-plus-one-queries", Description: "Detect N+1 issues"},
			{Name: "--schema", Description: "Display the schema of the composition spec"},
			{Name: "--format", Description: "Change the format of the input file"},
		},
	})
}
