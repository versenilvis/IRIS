package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pnpm",
		Description: "fast node packages",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages", Options: []core.Option{
				{Name: "--frozen-lockfile", Description: "no lockfile update"},
				{Name: "--prod", Description: "production only"},
			}},
			{Name: "add", Description: "add package", Options: []core.Option{
				{Name: "-D", Description: "save as devDependency"},
				{Name: "-g", Description: "install globally"},
				{Name: "-E", Description: "exact version"},
			}},
			{Name: "remove", Description: "remove package"},
			{Name: "run", Description: "run script", Generator: NpmScriptGenerator},
			{Name: "build", Description: "build project"},
			{Name: "dev", Description: "start dev server"},
			{Name: "test", Description: "run tests"},
			{Name: "update", Description: "update packages", Options: []core.Option{
				{Name: "--latest", Description: "update to latest"},
				{Name: "-i", Description: "interactive"},
			}},
			{Name: "store", Description: "manage store", Subcommands: []core.Subcommand{
				{Name: "prune", Description: "remove unreferenced"},
				{Name: "status", Description: "check store"},
				{Name: "path", Description: "print store path"},
			}},
			{Name: "exec", Description: "exec shell command"},
			{Name: "dlx", Description: "fetch and run package"},
			{Name: "ls", Description: "list packages"},
			{Name: "audit", Description: "security audit"},
			{Name: "patch", Description: "patch package"},
			{Name: "init", Description: "create package.json"},
			{Name: "publish", Description: "publish package"},
		},
		Options: []core.Option{
			{Name: "--filter", Description: "filter by package"},
			{Name: "-r", Description: "recursive"},
			{Name: "--workspace-root", Description: "run in workspace root"},
		},
	})
}
