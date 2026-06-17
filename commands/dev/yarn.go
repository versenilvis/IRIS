package dev

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "yarn",
		Description: "yarn package manager",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages", Options: []core.Option{
				{Name: "--frozen-lockfile", Description: "error on lockfile change"},
				{Name: "--production", Description: "production only"},
			}},
			{Name: "add", Description: "add package", Options: []core.Option{
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
			{Name: "upgrade", Description: "upgrade packages", Options: []core.Option{
				{Name: "--latest", Description: "upgrade to latest"},
				{Name: "-i", Description: "interactive"},
			}},
			{Name: "workspace", Description: "run command in workspace"},
			{Name: "workspaces", Description: "list workspaces", Subcommands: []core.Subcommand{
				{Name: "info", Description: "workspace info"},
				{Name: "run", Description: "run in all workspaces"},
			}},
			{Name: "link", Description: "symlink package"},
			{Name: "unlink", Description: "unlink package"},
			{Name: "audit", Description: "security audit"},
			{Name: "cache", Description: "manage cache", Subcommands: []core.Subcommand{
				{Name: "clean", Description: "clear cache"},
				{Name: "list", Description: "list cache"},
				{Name: "dir", Description: "print cache dir"},
			}},
			{Name: "init", Description: "create package.json", Options: []core.Option{
				{Name: "-2", Description: "berry (v2+)"},
			}},
			{Name: "publish", Description: "publish package"},
			{Name: "info", Description: "show package info"},
		},
		Options: []core.Option{
			{Name: "--cwd", Description: "working directory"},
			{Name: "--verbose", Description: "verbose output"},
		},
	})
}
