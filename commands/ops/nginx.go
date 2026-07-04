package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nginx",
		Description: "Nginx (pronounced",
		Options: []core.Option{
			{Name: "-c", Description: "Use an alternative configuration file"},
			{Name: "-e", Description: "Set global configuration directives"},
			{Name: "-p", Description: "Set the prefix path.  The default value is %%PREFIX%%"},
			{Name: "-q", Description: "Suppress non-error messages during configuration testing"},
			{Name: "-T", Description: "Same as -t, but additionally dump configuration files to standard output"},
			{Name: "-t", Description: "Print the nginx version, compiler version, and configure script parameters"},
			{Name: "-v", Description: "Print the nginx version"},
			{Name: "-?", Description: "Print help"},
			{Name: "-s", Description: "Sends SIGTERM"},
		},
	})
}
