package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-next-app",
		Description: "Output the version number",
		Options: []spec.Option{
			{Name: "-V", Description: "Output the version number"},
			{Name: "--ts", Description: "Initialize as a TypeScript project"},
			{Name: "--use-npm", Description: "Explicitly tell the CLI to bootstrap the app using npm"},
			{Name: "--use-pnpm", Description: "Explicitly tell the CLI to bootstrap the app using pnpm"},
			{Name: "-e", Description: "Display help for command"},
		},
	})
}
