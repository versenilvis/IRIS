package cc

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "premake",
		Description: "The premake5.lua file",
		Subcommands: []core.Subcommand{
			{Name: "clean", Description: "Remove all binaries and generated files"},
			{Name: "vs2022", Description: "Generate Visual Studio 2022 project files"},
			{Name: "vs2019", Description: "Generate Visual Studio 2019 project files"},
			{Name: "vs2017", Description: "Generate Visual Studio 2017 project files"},
			{Name: "vs2015", Description: "Generate Visual Studio 2015 project files"},
			{Name: "vs2013", Description: "Generate Visual Studio 2013 project files"},
			{Name: "vs2012", Description: "Generate Visual Studio 2012 project files"},
			{Name: "vs2010", Description: "Generate Visual Studio 2010 project files"},
			{Name: "vs2008", Description: "Generate Visual Studio 2008 project files"},
			{Name: "vs2005", Description: "Generate Visual Studio 2005 project files"},
			{Name: "gmake", Description: "Generate GNU Makefiles (This generator is deprecated by gmake2)"},
			{Name: "gmake2", Description: "Generate GNU Makefiles (including Cygwin and MinGW)"},
			{Name: "xcode4", Description: "Generate Apple Xcode 4 project files"},
			{Name: "codelite", Description: "Generate CodeLite project files"},
		},
		Options: []core.Option{
			{Name: "--file", Description: "The premake5.lua file"},
			{Name: "--debugger", Description: "Start MobDebug remote debugger. Works with ZeroBrane Studio"},
			{Name: "--fatal", Description: "Treat warnings from project scripts as errors"},
			{Name: "--insecure", Description: "Forfeit SSH certification checks"},
			{Name: "--interactive", Description: "Interactive command prompt"},
			{Name: "--os", Description: "Generate files for a different operating system"},
			{Name: "--scripts", Description: "Search for additional scripts on the given path"},
			{Name: "--systemscript", Description: "Override default system script (premake5-system.lua)"},
			{Name: "--verbose", Description: "Generate extra debug text output"},
			{Name: "--cc", Description: "Choose a C/C++ compiler set"},
			{Name: "--dc", Description: "Choose a D compiler"},
			{Name: "--dotnet", Description: "Choose a .NET compiler set"},
			{Name: "--help", Description: "Shows a complete list of the actions supported"},
			{Name: "--version", Description: "Display version information"},
		},
	})
}
