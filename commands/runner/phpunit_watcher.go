package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "phpunit-watcher",
		Description: "Automatically rerun PHPUnit tests when source code changes",
		Options: []spec.Option{
			{Name: "--filter", Description: "Watch a specific test"},
		},
	})
}
