package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "drush",
		Description: "Drush is a command line shell and Unix scripting interface for Drupal",
	})
}
