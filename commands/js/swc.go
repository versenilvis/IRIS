package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "swc",
		Description: "Path to the file",
		Options: []spec.Option{
			{Name: "--filename", Description: "Path to the file"},
			{Name: "--config-file", Description: "Path to a .swcrc file to use"},
			{Name: "--env-name", Description: "Whether or not to look up .swcrc files"},
			{Name: "--ignore", Description: "List of glob paths to not compile"},
			{Name: "--only", Description: "List of glob paths to only compile"},
			{Name: "--watch", Description: "Watch for changes and recompile"},
			{Name: "--quiet", Description: "Suppress compilation output"},
			{Name: "--source-maps", Description: "Generate source maps"},
			{Name: "--source-map-target", Description: "Define the file for the source map"},
			{Name: "--source-file-name", Description: "Set sources[0] on returned source map"},
			{Name: "--source-root", Description: "The root from which all sources are relative"},
			{Name: "--out-file", Description: "Compile all input files into a single file"},
			{Name: "--out-dir", Description: "Compile an input directory of modules into an output directory"},
			{Name: "--copy-files", Description: "When compiling a directory, copy over non-compilable files"},
			{Name: "--include-dotfiles", Description: "Include dotfiles when compiling and copying non-compilable files"},
			{Name: "--config", Description: "Override a config from .swcrc file"},
			{Name: "--sync", Description: "Invoke swc synchronously. Useful for debugging"},
			{Name: "--log-watch-compilation", Description: "Log a message when a watched file is successfully compiled"},
			{Name: "--extensions", Description: "Log a message when a watched file is successfully compiled"},
		},
	})
}
