package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "wasm-bindgen",
		Description: "Generate bindings between WebAssembly and JavaScript",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for wasm-bindgen"},
			{Name: "--version", Description: "Show version for wasm-bindgen"},
			{Name: "--out-dir", Description: "Output directory"},
			{Name: "--out-name", Description: "Set a custom output filename (Without extension. Defaults to crate name)"},
			{Name: "--target", Description: "What type of output to generate"},
			{Name: "--no-modules-global", Description: "Name of global to assign generated bindings to"},
			{Name: "--browser", Description: "Hint that JS should only be compatible with a browser"},
			{Name: "--typescript", Description: "Output a TypeScript definition file (on by default)"},
			{Name: "--no-typescript", Description: "Don't emit a *.d.ts file"},
			{Name: "--omit-imports", Description: "Don't emit imports in generated JavaScript"},
			{Name: "--debug", Description: "Include otherwise-extraneous debug checks in output"},
			{Name: "--no-demangle", Description: "Don't demangle Rust symbol names"},
			{Name: "--keep-debug", Description: "Keep debug sections in wasm files"},
			{Name: "--remove-name-section", Description: "Remove the debugging `name` section of the file"},
			{Name: "--remove-producers-section", Description: "Remove the telemetry `producers` section"},
			{Name: "--omit-default-module-path", Description: "Don't add WebAssembly fallback imports in generated JavaScript"},
			{Name: "--encode-into", Description: "Whether or not to use TextEncoder#encodeInto()"},
			{Name: "--nodejs", Description: "Deprecated, use `--target nodejs`"},
			{Name: "--web", Description: "Deprecated, use `--target web`"},
			{Name: "--no-modules", Description: "Deprecated, use `--target no-modules`"},
			{Name: "--weak-refs", Description: "Enable usage of the JS weak references proposal"},
			{Name: "--reference-types", Description: "Enable usage of WebAssembly reference types"},
		},
	})
}
