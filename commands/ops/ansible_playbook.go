package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ansible-playbook",
		Description: "Runs Ansible playbooks, executing the defined tasks on the targeted hosts",
		Options: []core.Option{
			{Name: "--ask-vault-pass", Description: "Ask for vault password"},
			{Name: "--flush-cache", Description: "Clears the fact cache for every host in inventory"},
			{Name: "--force-handlers", Description: "Run handlers even if a task fails"},
			{Name: "--list-hosts", Description: "Outputs a list of matching hosts; does not execute"},
			{Name: "--list-tags", Description: "List all available tags"},
			{Name: "--list-tasks", Description: "List all tasks that would be executed"},
			{Name: "--skip-tags", Description: "Only run plays and tasks whose tags do not match these values"},
			{Name: "--start-at-task", Description: "Start the playbook at the task matching this name one-step-at-a-time"},
			{Name: "--step", Description: "Execute one-step-at-a-time"},
			{Name: "--syntax-check", Description: "Perform a syntax check on the playbook, but do not execute it"},
			{Name: "--vault-id", Description: "Specify the vault identity to use"},
			{Name: "--vault-password-file", Description: "Specify a vault password file"},
			{Name: "--version", Description: "When changing (small) files and templates, show the differences in those files"},
			{Name: "--module-path", Description: "Prepend colon-separated path(s) to module library"},
			{Name: "--extra-vars", Description: "Set additional variables as key=value or YAML/JSON, if filename prepend with @"},
			{Name: "--forks", Description: "Specify number of parallel processes to use"},
			{Name: "--help", Description: "Show help for ansible"},
			{Name: "--inventory", Description: "Specify inventory host path or comma separated host list"},
			{Name: "--limit", Description: "Limit selected hosts to an additional pattern"},
			{Name: "--tags", Description: "Only run plays and tasks tagged with these values"},
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
