package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ansible-doc",
		Description: "Displays information on modules installed in Ansible libraries",
		Options: []spec.Option{
			{Name: "--metadata-dump", Description: "For internal testing only Dump json metadata for all plugins"},
			{Name: "--playbook-dir", Description: "Sets the relative path for many features including roles/ group_vars/ etc"},
			{Name: "--version", Description: "Show plugin names and their source files without summaries (implies --list)"},
			{Name: "--module-path", Description: "Prepend colon-separated path(s) to module library"},
			{Name: "--entry-point", Description: "Select the entry point for role(s)"},
			{Name: "--help", Description: "Show help and exit"},
			{Name: "--json", Description: "Change output into json format"},
			{Name: "--list", Description: "The path to the directory containing your roles"},
			{Name: "--snippet", Description: "Show playbook snippet for these plugin types: inventory, lookup, module"},
			{Name: "--type", Description: "Verbose mode (-vvv for more, -vvvv to enable connection debugging)"},
			{Name: "-v", Description: "Verbose mode (-vvv for more, -vvvv to enable connection debugging)"},
		},
	})
}
