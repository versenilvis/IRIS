package dev

import (
	"github.com/versenilvis/iris/commands/core"
)

func makePythonSpec(name string) *core.Spec {
	return &core.Spec{
		Name:        name,
		Description: "python interpreter",
		Generator:   core.FileGenerator(".py"),
		Options: []core.Option{
			{Name: "-m", Description: "run module as script"},
			{Name: "-c", Description: "run inline command"},
			{Name: "-i", Description: "interactive after script"},
			{Name: "-u", Description: "unbuffered stdio"},
			{Name: "-O", Description: "optimize bytecode"},
			{Name: "-v", Description: "verbose"},
			{Name: "-W", Description: "warning control"},
			{Name: "-B", Description: "no .pyc files"},
		},
	}
}

func init() {
	core.Register(makePythonSpec("python"))
	core.Register(makePythonSpec("python3"))
	core.Register(makePythonSpec("py"))

	core.Register(&core.Spec{
		Name:        "uv",
		Description: "fast python package manager",
		Subcommands: []core.Subcommand{
			{Name: "add", Description: "add dependency", Options: []core.Option{
				{Name: "--dev", Description: "dev dependency"},
				{Name: "--optional", Description: "optional dependency"},
				{Name: "-e", Description: "editable install"},
				{Name: "--index-url", Description: "custom index"},
			}},
			{Name: "remove", Description: "remove dependency", Generator: pipPackageGenerator},
			{Name: "sync", Description: "sync environment", Options: []core.Option{
				{Name: "--frozen", Description: "no lockfile update"},
				{Name: "--dev", Description: "include dev deps"},
			}},
			{Name: "run", Description: "run command in env", Generator: core.FileGenerator(".py")},
			{Name: "pip", Description: "pip interface", Subcommands: []core.Subcommand{
				{Name: "install", Description: "install packages"},
				{Name: "uninstall", Description: "uninstall packages", Generator: pipPackageGenerator},
				{Name: "list", Description: "list packages"},
				{Name: "freeze", Description: "freeze packages"},
				{Name: "show", Description: "show package", Generator: pipPackageGenerator},
				{Name: "compile", Description: "compile requirements"},
			}},
			{Name: "venv", Description: "create virtual env", Generator: core.FileGenerator("/"), Options: []core.Option{
				{Name: "--python", Description: "python version"},
				{Name: "--seed", Description: "install seed packages"},
			}},
			{Name: "init", Description: "init project"},
			{Name: "build", Description: "build package"},
			{Name: "publish", Description: "publish to pypi"},
			{Name: "lock", Description: "update lockfile"},
			{Name: "tree", Description: "show dependency tree"},
			{Name: "tool", Description: "manage tools", Subcommands: []core.Subcommand{
				{Name: "install", Description: "install tool"},
				{Name: "run", Description: "run tool"},
				{Name: "list", Description: "list tools"},
				{Name: "uninstall", Description: "uninstall tool"},
			}},
			{Name: "python", Description: "manage python versions", Subcommands: []core.Subcommand{
				{Name: "install", Description: "install python"},
				{Name: "list", Description: "list versions"},
				{Name: "find", Description: "find python"},
				{Name: "pin", Description: "pin python version"},
			}},
		},
	})

	core.Register(&core.Spec{
		Name:        "poetry",
		Description: "python dependency manager",
		Subcommands: []core.Subcommand{
			{Name: "add", Description: "add dependency", Options: []core.Option{
				{Name: "--dev", Description: "dev dependency"},
				{Name: "--group", Description: "dependency group"},
			}},
			{Name: "remove", Description: "remove dependency", Generator: pipPackageGenerator},
			{Name: "install", Description: "install dependencies"},
			{Name: "update", Description: "update dependencies"},
			{Name: "run", Description: "run in env", Generator: core.FileGenerator(".py")},
			{Name: "shell", Description: "activate virtual env"},
			{Name: "build", Description: "build package"},
			{Name: "publish", Description: "publish to pypi"},
			{Name: "show", Description: "show packages"},
			{Name: "check", Description: "check pyproject.toml"},
			{Name: "init", Description: "init new project"},
			{Name: "new", Description: "create new project"},
			{Name: "env", Description: "manage envs", Subcommands: []core.Subcommand{
				{Name: "info", Description: "show env info"},
				{Name: "list", Description: "list envs"},
				{Name: "use", Description: "switch python"},
				{Name: "remove", Description: "remove env"},
			}},
		},
	})
}
