package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "yarn",
		Description: "yarn package manager",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "install packages", Options: []spec.Option{
				{Name: "--frozen-lockfile", Description: "error on lockfile change"},
				{Name: "--production", Description: "production only"},
			}},
			{Name: "add", Description: "add package", Options: []spec.Option{
				{Name: "-D", Description: "dev dependency"},
				{Name: "-P", Description: "peer dependency"},
				{Name: "-E", Description: "exact version"},
				{Name: "-g", Description: "global"},
			}},
			{Name: "remove", Description: "remove package"},
			{Name: "run", Description: "run script", Generator: NpmScriptGenerator},
			{Name: "build", Description: "build project"},
			{Name: "dev", Description: "start dev server"},
			{Name: "test", Description: "run tests"},
			{Name: "upgrade", Description: "upgrade packages", Options: []spec.Option{
				{Name: "--latest", Description: "upgrade to latest"},
				{Name: "-i", Description: "interactive"},
			}},
			{Name: "workspace", Description: "run command in workspace"},
			{Name: "workspaces", Description: "list workspaces", Subcommands: []spec.Subcommand{
				{Name: "info", Description: "workspace info"},
				{Name: "run", Description: "run in all workspaces"},
			}},
			{Name: "link", Description: "symlink package"},
			{Name: "unlink", Description: "unlink package"},
			{Name: "audit", Description: "security audit"},
			{Name: "cache", Description: "manage cache", Subcommands: []spec.Subcommand{
				{Name: "clean", Description: "clear cache"},
				{Name: "list", Description: "list cache"},
				{Name: "dir", Description: "print cache dir"},
			}},
			{Name: "init", Description: "create package.json", Options: []spec.Option{
				{Name: "-2", Description: "berry (v2+)"},
			}},
			{Name: "publish", Description: "publish package"},
			{Name: "info", Description: "show package info"},
		},
		Options: []spec.Option{
			{Name: "--cwd", Description: "working directory"},
			{Name: "--verbose", Description: "verbose output"},
		},
	})
}
