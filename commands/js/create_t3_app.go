package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "create-t3-app",
		Description: "The name of the application, as well as the name of the directory to create",
		Options: []spec.Option{
			{Name: "--noGit", Description: "Boolean value if we're running in CI (default: false)"},
			{Name: "--tailwind", Description: "Install Tailwind CSS"},
			{Name: "--nextAuth", Description: "Install NextAuth.js"},
			{Name: "--prisma", Description: "Install Prisma"},
			{Name: "--trpc", Description: "Install tRPC"},
			{Name: "-i", Description: "Explicitly tell the CLI to use a custom import alias"},
			{Name: "-v", Description: "Display the version number"},
			{Name: "--help", Description: "Display help for command"},
		},
	})
}
