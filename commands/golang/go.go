package golang

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	globalBuildOptions := []core.Option{
		{Name: "-v", Description: "print package names"},
		{Name: "-n", Description: "print commands but do not run"},
		{Name: "-x", Description: "print commands as they run"},
		{Name: "-race", Description: "enable data race detection"},
		{Name: "-tags", Description: "comma-separated build tags"},
		{Name: "-ldflags", Description: "arguments to pass to linker"},
		{Name: "-gcflags", Description: "arguments to pass to compiler"},
		{Name: "-asmflags", Description: "arguments to pass to assembler"},
		{Name: "-compiler", Description: "name of compiler to use (gccgo or gc)"},
		{Name: "-work", Description: "print temporary work directory"},
		{Name: "-mod", Description: "module download mode (readonly, vendor, or mod)"},
		{Name: "-trimpath", Description: "remove file system paths from executable"},
		{Name: "-p", Description: "number of programs to run in parallel"},
		{Name: "-a", Description: "force rebuilding of up-to-date packages"},
	}

	buildModeSubcommands := []core.Subcommand{
		{Name: "archive", Description: "build non-main packages into .a files"},
		{Name: "c-archive", Description: "build main package into C archive"},
		{Name: "c-shared", Description: "build main package into C shared library"},
		{Name: "exe", Description: "build main packages into executables"},
		{Name: "pie", Description: "build position independent executables"},
		{Name: "plugin", Description: "build main packages into Go plugin"},
		{Name: "shared", Description: "combine packages into single shared library"},
	}

	core.Register(&core.Spec{
		Name:        "go",
		Description: "tool for managing Go source code",
		Subcommands: []core.Subcommand{
			{
				Name:        "build",
				Description: "compile packages and dependencies",
				MaxArgs:     1,
				Generator:   core.FileGenerator(".go"),
				Options: append(globalBuildOptions, 
					core.Option{Name: ".", Description: "current package"},
					core.Option{Name: "./...", Description: "all packages"},
					core.Option{Name: "-o", Description: "output file or directory"},
					core.Option{Name: "-i", Description: "install dependency packages"},
					core.Option{Name: "-buildmode", Description: "build mode to use"},
				),
				Subcommands: buildModeSubcommands,
			},
			{
				Name:        "run",
				Description: "compile and run Go program",
				MaxArgs:     1,
				Generator:   core.FileGenerator(".go"),
				Options: append(globalBuildOptions,
					core.Option{Name: ".", Description: "current package"},
					core.Option{Name: "-exec", Description: "invoke binary using xprog"},
				),
			},
			{
				Name:        "test",
				Description: "test packages",
				MaxArgs:     1,
				Generator:   core.FileGenerator(".go"),
				Options: append(globalBuildOptions,
					core.Option{Name: ".", Description: "current package"},
					core.Option{Name: "./...", Description: "all packages"},
					core.Option{Name: "-c", Description: "compile test binary but do not run"},
					core.Option{Name: "-i", Description: "install test dependencies"},
					core.Option{Name: "-json", Description: "convert output to JSON"},
					core.Option{Name: "-bench", Description: "run benchmarks"},
					core.Option{Name: "-run", Description: "run specific test regex"},
					core.Option{Name: "-cover", Description: "enable coverage report"},
				),
			},
			{
				Name:        "mod",
				Description: "module maintenance",
				Subcommands: []core.Subcommand{
					{Name: "init", Description: "initialize new module"},
					{Name: "tidy", Description: "add missing and remove unused modules"},
					{Name: "download", Description: "download modules to local cache"},
					{Name: "vendor", Description: "make vendored copy of dependencies"},
					{Name: "edit", Description: "edit go.mod from tools or scripts"},
					{Name: "graph", Description: "print module requirement graph"},
					{Name: "verify", Description: "verify dependencies have expected content"},
					{Name: "why", Description: "explain why packages or modules are needed"},
				},
			},
			{
				Name:        "get",
				Description: "add dependencies to current module and install them",
				Options: append(globalBuildOptions,
					core.Option{Name: "-u", Description: "update to newer minor/patch releases"},
					core.Option{Name: "-t", Description: "download modules needed for tests"},
					core.Option{Name: "-d", Description: "only download, do not install"},
				),
			},
			{
				Name:        "install",
				Description: "compile and install packages and dependencies",
				Options:     globalBuildOptions,
			},
			{
				Name:        "list",
				Description: "list packages or modules",
				Options: append(globalBuildOptions,
					core.Option{Name: "-m", Description: "list modules instead of packages"},
					core.Option{Name: "-u", Description: "add upgrade information"},
					core.Option{Name: "-json", Description: "print in JSON format"},
					core.Option{Name: "-f", Description: "specify alternate format"},
				),
			},
			{
				Name:        "fmt",
				Description: "gofmt (reformat) package sources",
				Options: []core.Option{
					{Name: "-n", Description: "print commands that would be executed"},
					{Name: "-x", Description: "print commands as they are executed"},
				},
			},
			{
				Name:        "vet",
				Description: "report likely mistakes in packages",
				Options:     globalBuildOptions,
			},
			{
				Name:        "env",
				Description: "print Go environment information",
				Options: []core.Option{
					{Name: "-json", Description: "print environment in JSON format"},
					{Name: "-u", Description: "unset named environment variables"},
					{Name: "-w", Description: "change default settings of environment variables"},
				},
			},
			{
				Name:        "version",
				Description: "print Go version",
				Options: []core.Option{
					{Name: "-m", Description: "print embedded module version info"},
					{Name: "-v", Description: "report unrecognized files"},
				},
			},
			{
				Name:        "clean",
				Description: "remove object files and cached files",
				Options: append(globalBuildOptions,
					core.Option{Name: "-i", Description: "remove installed archive or binary"},
					core.Option{Name: "-cache", Description: "remove entire go build cache"},
					core.Option{Name: "-modcache", Description: "remove entire module cache"},
				),
			},
			{
				Name:        "tool",
				Description: "run specified go tool",
				Options: []core.Option{
					{Name: "-n", Description: "print command but do not execute"},
				},
			},
			{
				Name:        "work",
				Description: "workspace maintenance",
				Subcommands: []core.Subcommand{
					{Name: "init", Description: "initialize workspace file"},
					{Name: "edit", Description: "edit go.work from tools or scripts"},
					{Name: "sync", Description: "sync workspace build list to modules"},
					{Name: "use", Description: "add modules to workspace file"},
				},
			},
			{
				Name:        "doc",
				Description: "show documentation for package or symbol",
				Options: []core.Option{
					{Name: "-all", Description: "show all documentation"},
					{Name: "-cmd", Description: "treat package main like regular package"},
					{Name: "-src", Description: "show full source code"},
				},
			},
			{
				Name:        "bug",
				Description: "start a bug report",
			},
			{
				Name:        "fix",
				Description: "update packages to use new APIs",
			},
			{
				Name:        "generate",
				Description: "generate Go files by processing source",
				Options:     globalBuildOptions,
			},
		},
	})
}
