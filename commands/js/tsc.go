package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "tsc",
		Description: "CLI tool for TypeScript compiler",
		Options: []core.Option{
			{Name: "--all", Description: "Show all compiler options"},
			{Name: "--generateTrace", Description: "Generates an event trace and a list of types"},
			{Name: "--help", Description: "Gives local information for help on the CLI"},
			{Name: "--init", Description: "Initializes a TypeScript project and creates a tsconfig.json file"},
			{Name: "--listFilesOnly", Description: "Print names of files that are part of the compilation and then stop processing"},
			{Name: "--locale", Description: "Set the language of the messaging from TypeScript. This does not affect emit"},
			{Name: "--project", Description: "Print the final configuration instead of building"},
			{Name: "--version", Description: "Print the compiler's version"},
			{Name: "--build", Description: "Build one or more projects and their dependencies, if out of date"},
			{Name: "--clean", Description: "Delete the outputs of all projects"},
			{Name: "--dry", Description: "Show what would be built (or deleted, if specified with '--clean')"},
			{Name: "--force", Description: "Build all projects, including those that appear to be up to date"},
			{Name: "--verbose", Description: "Enable verbose logging"},
			{Name: "--excludeDirectories", Description: "Remove a list of directories from the watch process"},
			{Name: "--excludeFiles", Description: "Remove a list of files from the watch mode's processing"},
			{Name: "--fallbackPolling", Description: "Watch input files"},
			{Name: "--watchDirectory", Description: "Specify how the TypeScript watch mode works"},
			{Name: "--allowJs", Description: "Allow 'import x from y' when a module doesn't have a default export"},
			{Name: "--allowUmdGlobalAccess", Description: "Allow accessing UMD globals from modules"},
			{Name: "--allowUnreachableCode", Description: "Disable error reporting for unreachable code"},
			{Name: "--allowUnusedLabels", Description: "Disable error reporting for unused label"},
			{Name: "--alwaysStrict", Description: "Ensure 'use strict' is always emitted"},
			{Name: "--checkJs", Description: "Enable error reporting in type-checked JavaScript files"},
			{Name: "--composite", Description: "Generate .d.ts files from TypeScript and JavaScript files in your project"},
			{Name: "--declarationDir", Description: "Specify the output directory for generated declaration files"},
			{Name: "--declarationMap", Description: "Create sourcemaps for d.ts files"},
			{Name: "--diagnostics", Description: "Output compiler performance information after building"},
			{Name: "--disableSizeLimit", Description: "Opt a project out of multi-project reference checking when editing"},
			{Name: "--emitBOM", Description: "Emit a UTF-8 Byte Order Mark (BOM) in the beginning of output files"},
			{Name: "--emitDeclarationOnly", Description: "Only output d.ts files and not JavaScript files"},
			{Name: "--emitDecoratorMetadata", Description: "Emit design-type metadata for decorated declarations in source files"},
			{Name: "--esModuleInterop", Description: "Differentiate between undefined and not present when type checking"},
			{Name: "--experimentalDecorators", Description: "Enable experimental support for TC39 stage 2 draft decorators"},
			{Name: "--explainFiles", Description: "Print files read during the compilation including why it was included"},
			{Name: "--extendedDiagnostics", Description: "Output more detailed compiler performance information after building"},
			{Name: "--generateCpuProfile", Description: "Emit a v8 CPU profile of the compiler run for debugging"},
			{Name: "--importHelpers", Description: "Specify emit/checking behavior for imports that are only used for types"},
			{Name: "--incremental", Description: "Save .tsbuildinfo files to allow for incremental compilation of projects"},
			{Name: "--inlineSourceMap", Description: "Include sourcemap files inside the emitted JavaScript"},
			{Name: "--inlineSources", Description: "Include source code in the sourcemaps inside the emitted JavaScript"},
		},
	})
}
