package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "create-react-native-app",
		Description: "Creates a new React Native project",
		Options: []core.Option{
			{Name: "--template", Description: "The path inside of a GitHub repo where the example lives"},
			{Name: "--yes", Description: "Use the default options for creating a project"},
			{Name: "--no-install", Description: "Skip installing npm packages or CocoaPods"},
			{Name: "--use-npm", Description: "Use npm to install dependencies. (default when Yarn is not installed)"},
			{Name: "-h", Description: "Output usage information"},
			{Name: "-V", Description: "Output the version number"},
		},
	})
}
