package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "visudo",
		Description: "Checking existing sudoers file for syntax errors",
		Options: []spec.Option{
			{Name: "--check", Description: "Checking existing sudoers file for syntax errors"},
			{Name: "--file", Description: "Set an alternative sudoers file location"},
			{Name: "--help", Description: "Display a short help message"},
			{Name: "--quiet", Description: "Enable quiet mode (syntax error not printed)"},
			{Name: "--strict", Description: "Enable strict checking of the sudoers file"},
			{Name: "--version", Description: "Display version and exit"},
			{Name: "--export", Description: "Export JSON and write it to output_file"},
			{Name: "--perms", Description: "Enforce default mode for the sudoers file"},
			{Name: "--owner", Description: "Enforce the default ownership for the sudoers file"},
		},
	})
}
