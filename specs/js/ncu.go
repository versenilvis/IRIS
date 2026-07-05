package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ncu",
		Description: "Clear the default cache, or the cache file specified by --cacheFile",
		Options: []core.Option{
			{Name: "--cache", Description: "Clear the default cache, or the cache file specified by --cacheFile"},
			{Name: "--cacheExpiration", Description: "Cache expiration in minutes. Only works with --cache. (default: 10)"},
			{Name: "--cacheFile", Description: "Force color in terminal"},
			{Name: "--concurrency", Description: "Max number of concurrent HTTP requests to registry. (default: 8)"},
			{Name: "--configFileName", Description: "Config file name. (default: .ncurc.{json,yml,js,cjs})"},
			{Name: "--configFilePath", Description: "Directory of .ncurc config file. (default: directory of packageFile)"},
			{Name: "--cwd", Description: "Working directory in which npm will be executed"},
			{Name: "--deep", Description: "Include deprecated packages"},
			{Name: "--doctor", Description: "Specifies the install script to use in doctor mode. (default: npm install/yarn)"},
			{Name: "--doctorTest", Description: "Specifies the test script to use in doctor mode. (default: npm test)"},
			{Name: "--enginesNode", Description: "Include only packages that satisfy engines.node as specified in the package file"},
			{Name: "--errorLevel", Description: "Check global packages instead of in the current project"},
			{Name: "--groupFunction", Description: "Show help"},
			{Name: "--interactive", Description: "Output new package file instead of human-readable message"},
			{Name: "--jsonDeps", Description: "Output upgraded dependencies in json"},
			{Name: "--loglevel", Description: "Package file data (you can also use stdin)"},
			{Name: "--packageFile", Description: "Package file(s) location. (default: ./package.json)"},
			{Name: "--packageManager", Description: "Current working directory of npm"},
			{Name: "--registry", Description: "Remove version ranges from the final package version"},
			{Name: "--retry", Description: "Number of times to retry failed requests for package info. (default: 3)"},
			{Name: "--root", Description: "Don't output anything. Alias for --loglevel silent"},
			{Name: "--stdin", Description: "Read package.json from stdin"},
			{Name: "--target", Description: "Log additional information for debugging. Alias for --loglevel verbose"},
			{Name: "--version", Description: "Output the version number of npm-check-updates"},
			{Name: "--workspace", Description: "Run on all workspaces. Add --root to also upgrade the root project"},
		},
	})
}
