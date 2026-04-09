package commands

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "go",
		Description: "go toolchain",
		Subcommands: []core.Subcommand{
			{
				Name: "build",
				Description: "compile packages",
				Options: []core.Option{
					{Name: "-o", Description: "output file"},
					{Name: "-v", Description: "verbose"},
					{Name: "-race", Description: "enable race detector"},
					{Name: "-ldflags", Description: "linker flags"},
				},
			},
			{
				Name: "run",
				Description: "compile and run",
				Options: []core.Option{
					{Name: "-race", Description: "enable race detector"},
				},
			},
			{
				Name: "test",
				Description: "run tests",
				Options: []core.Option{
					{Name: "-v", Description: "verbose output"},
					{Name: "-run", Description: "run specific test"},
					{Name: "-bench", Description: "run benchmarks"},
					{Name: "-cover", Description: "coverage report"},
					{Name: "-count", Description: "run count times"},
					{Name: "-race", Description: "enable race detector"},
				},
			},
			{
				Name: "mod",
				Description: "module maintenance",
				Subcommands: []core.Subcommand{
					{Name: "init", Description: "initialize module"},
					{Name: "tidy", Description: "clean up modules"},
					{Name: "download", Description: "download modules"},
					{Name: "vendor", Description: "make vendor copy"},
					{Name: "edit", Description: "edit go.mod"},
					{Name: "graph", Description: "print module graph"},
					{Name: "verify", Description: "verify dependencies"},
				},
			},
			{
				Name: "get",
				Description: "add dependency",
				Options: []core.Option{
					{Name: "-u", Description: "update module"},
				},
			},
			{
				Name: "fmt",
				Description: "format source",
			},
			{
				Name: "vet",
				Description: "examine code",
			},
			{
				Name: "install",
				Description: "compile and install",
			},
			{
				Name: "generate",
				Description: "run go generate",
			},
			{
				Name: "clean",
				Description: "remove object files",
				Options: []core.Option{
					{Name: "-cache", Description: "clean build cache"},
					{Name: "-testcache", Description: "clean test cache"},
				},
			},
			{
				Name: "env",
				Description: "print go environment",
			},
			{
				Name: "doc",
				Description: "show documentation",
			},
			{
				Name: "tool",
				Description: "run go tool",
			},
		},
	})
}
