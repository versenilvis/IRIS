package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-web3-frontend",
		Description: "Quickly create a Next.js project with wagmi and TailwindCSS ready to go",
		Options: []spec.Option{
			{Name: "--ts", Description: "Initialize as a TypeScript project"},
			{Name: "--use-npm", Description: "Explicitly tell the CLI to bootstrap the app using npm"},
			{Name: "--use-pnpm", Description: "Explicitly tell the CLI to bootstrap the app using pnpm"},
			{Name: "-h", Description: "Display help for command"},
		},
	})
}
