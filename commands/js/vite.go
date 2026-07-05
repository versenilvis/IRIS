package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "vite",
		Description: "Native ESM-powered web dev build tool",
		Options: []spec.Option{
			{Name: "-h", Description: "Show help message"},
			{Name: "--host", Description: "Specify the hostname"},
			{Name: "--port", Description: "Specify the port"},
			{Name: "--https", Description: "Use TLS + HTTP/2"},
			{Name: "--open", Description: "Open browser on startup"},
			{Name: "--cors", Description: "Enable CORS"},
			{Name: "--strictPort", Description: "Exit if the port is already in use"},
			{Name: "--force", Description: "For the optimizer to re-bundle"},
			{Name: "-c", Description: "Use the specified config file"},
			{Name: "--base", Description: "Public base path"},
			{Name: "-l", Description: "Set the log level"},
			{Name: "--clearScreen", Description: "Enable screen clearing when logging"},
			{Name: "-d", Description: "Show debug logs"},
			{Name: "-f", Description: "Filter debug logs"},
			{Name: "-m", Description: "Set env mode"},
			{Name: "-v", Description: "Display version number"},
			{Name: "--target", Description: "Transpile target (must be a valid esbuild target)"},
			{Name: "--outDir", Description: "Output directory"},
			{Name: "--assetsDir", Description: "Directory under outDir to place assets in"},
			{Name: "--assetsInlineLimit", Description: "Static asset base64 inline threshold in bytes"},
			{Name: "--ssr", Description: "Build specified entry for server-side rendering"},
			{Name: "--sourcemap", Description: "Output sourcemaps for build"},
			{Name: "--minify", Description: "Enable minification"},
			{Name: "--manifest", Description: "Emit build manifest json"},
			{Name: "--ssrManifest", Description: "Emit ssr manifest json"},
			{Name: "--emptyOutDir", Description: "Force empty outDir when it's outside of root"},
			{Name: "-w", Description: "Rebuilds when modules have changed on disk"},
		},
	})
}
