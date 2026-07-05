package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ansible",
		Description: "Define and run a single Ansible task",
		Options: []spec.Option{
			{Name: "--ask-vault-pass", Description: "Ask for vault password"},
			{Name: "--list-hosts", Description: "Outputs a list of matching hosts; does not execute"},
			{Name: "--playbook-dir", Description: "Perform a syntax check on the playbook, but do not execute it"},
			{Name: "--vault-id", Description: "Specify the vault identity to use"},
			{Name: "--vault-password-file", Description: "Specify a vault password file"},
			{Name: "--version", Description: "Run asynchronously, failing after specified seconds"},
			{Name: "--check", Description: "When changing (small) files and templates, show the differences in those files"},
			{Name: "--module-path", Description: "Prepend colon-separated path(s) to module library"},
			{Name: "--poll", Description: "Set the poll interval if using -B"},
			{Name: "--args", Description: "Specify module arguments"},
			{Name: "--extra-vars", Description: "Set additional variables as key=value or YAML/JSON, if filename prepend with @"},
			{Name: "--forks", Description: "Specify number of parallel processes to use"},
			{Name: "--help", Description: "Show help for ansible"},
			{Name: "--inventory", Description: "Specify inventory host path or comma separated host list"},
			{Name: "--limit", Description: "Limit selected hosts to an additional pattern"},
			{Name: "--module-name", Description: "Specify the module name to execute"},
			{Name: "--one-line", Description: "Condense output"},
			{Name: "--tree", Description: "Log output to specific directory"},
			{Name: "--verbose", Description: "Enable verbose mode"},
			{Name: "-vvv", Description: "Enable very verbose mode"},
			{Name: "-vvvv", Description: "Enable connection debug mode"},
			{Name: "--become-method", Description: "Privilege escalation method to use"},
			{Name: "--become-user", Description: "Privilege escalation user to use"},
			{Name: "--ask-become-pass", Description: "Prompt for privilege escalation password"},
			{Name: "--become", Description: "Run operations with become"},
			{Name: "--private-key", Description: "Use this fole to authenticate the connection"},
			{Name: "--scp-extra-args", Description: "Extra arguments to pass to (only) scp"},
			{Name: "--sftp-extra-args", Description: "Extra arguments to pass to (only) sftp"},
			{Name: "-ssh-extra-args", Description: "Extra arguments to pass to (only) ssh"},
			{Name: "--ssh-common-args", Description: "Extra arguments to pass to sftp/scp/ssh"},
			{Name: "--timeout", Description: "Override the connection timeout in seconds"},
			{Name: "--connection", Description: "Connection type to use"},
			{Name: "--ask-pass", Description: "Ask for connection password"},
			{Name: "--user", Description: "Connect as this user"},
		},
	})
}
