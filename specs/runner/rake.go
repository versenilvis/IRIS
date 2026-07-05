package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rake",
		Description: "A ruby build program with capabilities similar to make",
		Options: []core.Option{
			{Name: "-n", Description: "Do a dry run without executing actions"},
			{Name: "-h", Description: "Display this help message"},
			{Name: "-I", Description: "Include LIBDIR in the search path for required modules"},
			{Name: "-P", Description: "Display the tasks and dependencies, then exit"},
			{Name: "-q", Description: "Do not log messages to standard output"},
			{Name: "-f", Description: "Use FILE as the rakefile"},
			{Name: "-r", Description: "Require MODULE before executing rakefile"},
			{Name: "-s", Description: "Like --quiet, but also suppresses the 'in directory' announcement"},
			{Name: "-T", Description: "Display the tasks and dependencies, then exit"},
			{Name: "-t", Description: "Turn on invoke/execute tracing, enable full backtrace"},
			{Name: "-v", Description: "Log message to standard output (default)"},
			{Name: "-V", Description: "Display the program version"},
		},
	})
}
