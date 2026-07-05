package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ansible-lint",
		Description: "Ansible static code analysis",
		Options: []spec.Option{
			{Name: "-f", Description: "Ansible static code analysis"},
			{Name: "--help", Description: "Show help for ansible-lint"},
			{Name: "--list-rules", Description: "List all the rules"},
			{Name: "--list-tags", Description: "List all the tags and the rules they cover"},
			{Name: "--format", Description: "Stdout formatting"},
			{Name: "-q", Description: "Quieter, reduce verbosity, can be specified twice"},
			{Name: "-p", Description: "Parseable output, same as '-f pep8'"},
			{Name: "--progressive", Description: "Specify custom rule directories"},
			{Name: "-R", Description: "Keep using embedded rules when using '-r'"},
			{Name: "--write", Description: "Allow ansible-lint to reformat YAML files and run rule transforms"},
			{Name: "--show-relpath", Description: "Display path relative to CWD"},
			{Name: "--tags", Description: "Only check rules whose id/tags match these values"},
			{Name: "-v", Description: "Increase verbosity level (-vv for more)"},
			{Name: "--skip-list", Description: "Only check rules whose id/tags do not match these values"},
			{Name: "--warn-list", Description: "Activate optional rules by their tag name"},
			{Name: "--nocolor", Description: "Disable colored output, same as NO_COLOR=1"},
			{Name: "--force-color", Description: "Force colored output, same as FORCE_COLOR=1"},
			{Name: "--exclude-paths", Description: "Path to directories or files to skip. This option is repeatable"},
			{Name: "--config-file", Description: "Disable installation of requirements.yml"},
			{Name: "--version", Description: "Show version of ansible-lint"},
		},
	})
}
