package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "aws-vault",
		Description: "Add credentials to the secure keystore",
		Options: []spec.Option{
			{Name: "-f", Description: "Force-remove the profile without a prompt"},
			{Name: "--profiles", Description: "Show only the profile names"},
			{Name: "--sessions", Description: "Show only the session names"},
			{Name: "--credentials", Description: "Show only the profiles with stored credential"},
			{Name: "-n", Description: "Use master credentials, no session or role used"},
			{Name: "-d", Description: "Duration of the temporary or assume-role session. Defaults to 1h"},
			{Name: "--region", Description: "The AWS region"},
			{Name: "-t", Description: "The MFA token to use"},
			{Name: "--help", Description: "Show context-sensitive help (also try --help-long and --help-man)"},
			{Name: "--version", Description: "Show application version"},
			{Name: "--debug", Description: "Show debugging output"},
		},
	})
}
