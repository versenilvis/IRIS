package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "klist",
		Description: "Credential cache to list",
		Options: []core.Option{
			{Name: "-c", Description: "Credential cache to list"},
			{Name: "--cache", Description: "Credential cache to list"},
			{Name: "-s", Description: "Display AFS tokens"},
			{Name: "-5", Description: "Display v5 cred cache (this is the default)"},
			{Name: "-f", Description: "Include ticket flags in short form, each character stands for a specific flag"},
			{Name: "-v", Description: "Verbose output. Include all possible information"},
			{Name: "-l", Description: "JSON formatted output"},
			{Name: "--hidden", Description: "Verbose output"},
		},
	})
}
