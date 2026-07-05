package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "drush",
		Description: "Drush is a command line shell and Unix scripting interface for Drupal",
	})
}
