package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "phpunit-watcher",
		Description: "Automatically rerun PHPUnit tests when source code changes",
		Options: []core.Option{
			{Name: "--filter", Description: "Watch a specific test"},
		},
	})
}
