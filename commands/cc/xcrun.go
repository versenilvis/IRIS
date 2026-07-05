package cc

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xcrun",
		Description: "SceneKit CLI utilities",
		Subcommands: []spec.Subcommand{
			{Name: "scntool", Description: "SceneKit CLI utilities"},
		},
		Options: []spec.Option{
			{Name: "--convert", Description: "File to convert"},
			{Name: "--format", Description: "Format to convert to"},
			{Name: "--output", Description: "Path to the output file"},
			{Name: "--force-y-up", Description: "Convert objects to use y axis up"},
			{Name: "--force-interleaved", Description: "Use .ktx, .astc and .pvrtc files for textures if available in the asset catalog"},
			{Name: "--verbose", Description: "Get pretty error message"},
			{Name: "-h", Description: "Help message"},
			{Name: "--version", Description: "Show the xcrun version"},
			{Name: "-v", Description: "Show verbose logging output"},
			{Name: "--sdk", Description: "Find the tool for the given SDK name"},
			{Name: "--toolchain", Description: "Find the tool for the given toolchain"},
			{Name: "-l", Description: "Show command path to be executed (and --run)"},
			{Name: "-f", Description: "Find and print the tool path"},
			{Name: "--run", Description: "(Default) find and execute the tool"},
			{Name: "-n", Description: "Do not use the lookup cache"},
			{Name: "-k", Description: "Invalidate all existing cache entries"},
			{Name: "--show-sdk-path", Description: "Show selected SDK install path"},
			{Name: "--show-sdk-version", Description: "Show selected SDK version"},
			{Name: "--show-sdk-build-version", Description: "Show selected SDK build version"},
			{Name: "--show-sdk-platform-path", Description: "Show selected SDK platform path"},
			{Name: "--show-sdk-platform-version", Description: "Show selected SDK platform version"},
		},
	})
}
