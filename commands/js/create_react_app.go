package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-react-app",
		Description: "Creates a new React project",
		Options: []spec.Option{
			{Name: "--template", Description: "Use npm to install dependencies (default when Yarn is not installed)"},
			{Name: "--use-pnp", Description: "Use Yarn Plug'n'Play to create the app"},
			{Name: "--scripts-version", Description: "Print additional logs"},
			{Name: "-h", Description: "Output usage information"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
