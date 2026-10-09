package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "sudoedit",
		Description: "Edit files as another user",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-A", Description: "Use helper program for password prompting"},
			{Name: "--askpass", Description: "Use helper program for password prompting"},
			{Name: "-C", Description: "Close all file descriptors greater than or equal to num"},
			{Name: "--close-from", Description: "Close all file descriptors greater than or equal to num"},
			{Name: "-D", Description: "Change directory before running command"},
			{Name: "--chdir", Description: "Change directory before running command"},
			{Name: "-g", Description: "Run command as the specified group name or ID"},
			{Name: "--group", Description: "Run command as the specified group name or ID"},
			{Name: "-h", Description: "Display a short help message"},
			{Name: "--help", Description: "Display a short help message"},
			{Name: "-k", Description: "Invalidate the user's cached credentials"},
			{Name: "--reset-timestamp", Description: "Invalidate the user's cached credentials"},
			{Name: "-n", Description: "Avoid prompting the user for input"},
			{Name: "--non-interactive", Description: "Avoid prompting the user for input"},
			{Name: "-p", Description: "Use the specified password prompt"},
			{Name: "--prompt", Description: "Use the specified password prompt"},
			{Name: "-R", Description: "Change root directory before running command"},
			{Name: "--chroot", Description: "Change root directory before running command"},
			{Name: "-S", Description: "Read password from standard input"},
			{Name: "--stdin", Description: "Read password from standard input"},
			{Name: "-u", Description: "Run command as the specified user name or ID"},
			{Name: "--user", Description: "Run command as the specified user name or ID"},
			{Name: "-V", Description: "Display version information and exit"},
			{Name: "--version", Description: "Display version information and exit"},
		},
	})
}
