package golang

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	globalBuildOptions := []spec.Option{
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

	buildModeSubcommands := []spec.Subcommand{
		{Name: "archive", Description: "build non-main packages into .a files"},
		{Name: "c-archive", Description: "build main package into C archive"},
		{Name: "c-shared", Description: "build main package into C shared library"},
		{Name: "exe", Description: "build main packages into executables"},
		{Name: "pie", Description: "build position independent executables"},
		{Name: "plugin", Description: "build main packages into Go plugin"},
		{Name: "shared", Description: "combine packages into single shared library"},
	}

	spec.Register(&spec.Spec{
		Name:        "go",
		Description: "tool for managing Go source code",
		Subcommands: []spec.Subcommand{
			{
				Name:        "build",
				Description: "compile packages and dependencies",
				MaxArgs:     1,
				Generator:   spec.FileGenerator(".go"),
				Options: append(globalBuildOptions, 
					spec.Option{Name: ".", Description: "current package"},
					spec.Option{Name: "./...", Description: "all packages"},
					spec.Option{Name: "-o", Description: "output file or directory"},
					spec.Option{Name: "-i", Description: "install dependency packages"},
					spec.Option{Name: "-buildmode", Description: "build mode to use"},
				),
				Subcommands: buildModeSubcommands,
			},
			{
				Name:        "run",
				Description: "compile and run Go program",
				MaxArgs:     1,
				Generator:   spec.FileGenerator(".go"),
				Options: append(globalBuildOptions,
					spec.Option{Name: ".", Description: "current package"},
					spec.Option{Name: "-exec", Description: "invoke binary using xprog"},
				),
			},
			{
				Name:        "test",
				Description: "test packages",
				MaxArgs:     1,
				Generator:   spec.FileGenerator(".go"),
				Options: append(globalBuildOptions,
					spec.Option{Name: ".", Description: "current package"},
					spec.Option{Name: "./...", Description: "all packages"},
					spec.Option{Name: "-c", Description: "compile test binary but do not run"},
					spec.Option{Name: "-i", Description: "install test dependencies"},
					spec.Option{Name: "-json", Description: "convert output to JSON"},
					spec.Option{Name: "-bench", Description: "run benchmarks"},
					spec.Option{Name: "-run", Description: "run specific test regex"},
					spec.Option{Name: "-cover", Description: "enable coverage report"},
				),
			},
			{
				Name:        "mod",
				Description: "module maintenance",
				Subcommands: []spec.Subcommand{
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
					spec.Option{Name: "-u", Description: "update to newer minor/patch releases"},
					spec.Option{Name: "-t", Description: "download modules needed for tests"},
					spec.Option{Name: "-d", Description: "only download, do not install"},
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
					spec.Option{Name: "-m", Description: "list modules instead of packages"},
					spec.Option{Name: "-u", Description: "add upgrade information"},
					spec.Option{Name: "-json", Description: "print in JSON format"},
					spec.Option{Name: "-f", Description: "specify alternate format"},
				),
			},
			{
				Name:        "fmt",
				Description: "gofmt (reformat) package sources",
				Options: []spec.Option{
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
				Options: []spec.Option{
					{Name: "-json", Description: "print environment in JSON format"},
					{Name: "-u", Description: "unset named environment variables"},
					{Name: "-w", Description: "change default settings of environment variables"},
				},
			},
			{
				Name:        "version",
				Description: "print Go version",
				Options: []spec.Option{
					{Name: "-m", Description: "print embedded module version info"},
					{Name: "-v", Description: "report unrecognized files"},
				},
			},
			{
				Name:        "clean",
				Description: "remove object files and cached files",
				Options: append(globalBuildOptions,
					spec.Option{Name: "-i", Description: "remove installed archive or binary"},
					spec.Option{Name: "-cache", Description: "remove entire go build cache"},
					spec.Option{Name: "-modcache", Description: "remove entire module cache"},
				),
			},
			{
				Name:        "tool",
				Description: "run specified go tool",
				Options: []spec.Option{
					{Name: "-n", Description: "print command but do not execute"},
				},
			},
			{
				Name:        "work",
				Description: "workspace maintenance",
				Subcommands: []spec.Subcommand{
					{Name: "init", Description: "initialize workspace file"},
					{Name: "edit", Description: "edit go.work from tools or scripts"},
					{Name: "sync", Description: "sync workspace build list to modules"},
					{Name: "use", Description: "add modules to workspace file"},
				},
			},
			{
				Name:        "doc",
				Description: "show documentation for package or symbol",
				Options: []spec.Option{
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
