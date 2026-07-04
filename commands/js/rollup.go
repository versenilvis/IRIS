package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "rollup",
		Description: "Next-generation ES module bundler",
		Options: []core.Option{
			{Name: "-c", Description: "Directory for chunks (if absent, prints to stdout)"},
			{Name: "-e", Description: "Comma-separate list of module IDs to exclude"},
			{Name: "-f", Description: "Type of output (amd, cjs, es, iife, umd, system)"},
			{Name: "-g", Description: "Comma-separate list of `moduleID:Global` pairs"},
			{Name: "-h", Description: "Show help message"},
			{Name: "-i", Description: "Input (alternative to <entry file>)"},
			{Name: "-m", Description: "Generate sourcemap (`-m inline` for inline map)"},
			{Name: "-n", Description: "Name for UMD export"},
			{Name: "-o", Description: "Single output file (if absent, prints to stdout)"},
			{Name: "-p", Description: "Use the plugin specified (may be repeated)"},
			{Name: "-v", Description: "Show version number"},
			{Name: "-w", Description: "Watch files in bundle and rebuild on changes"},
			{Name: "--assetFileNames", Description: "Name pattern for emitted assets"},
			{Name: "--banner", Description: "Code to insert at top of bundle (outside wrapper)"},
			{Name: "--chunkFileNames", Description: "Name pattern for emitted secondary chunks"},
			{Name: "--compact", Description: "Minify wrapper code"},
			{Name: "--context", Description: "Specify top-level 'this' value"},
			{Name: "--entryFileNames", Description: "Name pattern for emitted entry chunks"},
			{Name: "--environment", Description: "Settings passed to config file"},
			{Name: "--no-esModule", Description: "Do not add __esmodule property"},
			{Name: "--exports", Description: "Specify export mode (auto, default, named, none)"},
			{Name: "--extend", Description: "Extend global variable defined by --name"},
			{Name: "--no-externalLiveBindings", Description: "Do not generate code to support live bindings"},
			{Name: "--failAfterWarnings", Description: "Exit with an error if the build produced warnings"},
			{Name: "--footer", Description: "Code to insert at end of bundle (outside wrapper)"},
			{Name: "--no-freeze", Description: "Do not freeze namespace objects"},
			{Name: "--no-hoistTransistiveImports", Description: "Do not hoist transistive imports into entry chunks"},
			{Name: "--no-indent", Description: "Don't indent result"},
			{Name: "--no-interop", Description: "Do not include interop block"},
			{Name: "--inlineDynamicImports", Description: "Create a single bundle when using dynamic imports"},
			{Name: "--intro", Description: "Code to insert at top of bundle (inside wrapper)"},
			{Name: "--minifyInternalImports", Description: "Force or disable minification of internal imports"},
			{Name: "--namespaceToStringTag", Description: "Create proper '.toString' methods for namespaces"},
			{Name: "--noConflict", Description: "Generate a noConflict method for UMD globals"},
			{Name: "--outro", Description: "Code to insert at end of bundle (inside wrapper)"},
			{Name: "--preferConst", Description: "Use 'const' instead of 'var' for exports"},
			{Name: "--no-preserveEntrySignatures", Description: "Avoid facade chunks for entry points"},
			{Name: "--preserveModules", Description: "Preserve module structure"},
			{Name: "--preserveModulesRoot", Description: "Put preserved modules under this path at root level"},
			{Name: "--preserveSymlinks", Description: "Do not follow symlinks when resolving files"},
		},
	})
}
