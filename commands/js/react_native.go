package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "react-native",
		Description: "Attempt to fix all diagnosed issues",
		Options: []core.Option{
			{Name: "--fix", Description: "Attempt to fix all diagnosed issues"},
			{Name: "--contributor", Description: "Output usage information"},
			{Name: "--entry-file", Description: "Path to the root JS file, either absolute or relative to JS root"},
			{Name: "--platform", Description: "Specify a custom transformer to be used"},
			{Name: "--dev", Description: "If false, warnings are disabled and the bundle is minified (default: true)"},
			{Name: "--minify", Description: "File name where to store the resulting bundle, ex. /tmp/groups.bundle"},
			{Name: "--bundle-encoding", Description: "Path to make sourcemap's sources entries relative to, ex. /root/dir"},
			{Name: "--sourcemap-use-absolute-path", Description: "Report SourceMapURL using its full path"},
			{Name: "--assets-dest", Description: "Directory name where to store assets referenced in the bundle"},
			{Name: "--unstable-transform-profile", Description: "Removes cached files"},
			{Name: "--read-global-cache", Description: "Try to fetch transformed JS code from the global cache, if configured"},
			{Name: "--config", Description: "Path to the CLI configuration file"},
			{Name: "-h", Description: "Output usage information"},
			{Name: "--version", Description: "Shortcut for `--template react-native@version`"},
			{Name: "--template", Description: "Forces using npm for initialization"},
			{Name: "--directory", Description: "Uses a custom directory instead of `<projectName>`"},
			{Name: "--title", Description: "Uses a custom app title name for application"},
			{Name: "--skip-install", Description: "Skips dependencies installation step"},
			{Name: "--port", Description: "Port on which to listen to"},
			{Name: "--host", Description: "Change the default host"},
			{Name: "--projectRoot", Description: "Path to a custom project root"},
			{Name: "--watchFolders", Description: "Specify any additional folders to be added to the watch list"},
			{Name: "--assetPlugins", Description: "Specify any additional asset plugins to be used by the packager by full filepath"},
			{Name: "--sourceExts", Description: "Specify any additional source extensions to be used by the packager"},
			{Name: "--max-workers", Description: "Specify a custom transformer to be used"},
			{Name: "--reset-cache", Description: "Removes cached files"},
			{Name: "--custom-log-reporter-path", Description: "Enables logging"},
			{Name: "--https", Description: "Enables https connections to the server"},
			{Name: "--key", Description: "Path to custom SSL key"},
			{Name: "--cert", Description: "Path to custom SSL cert"},
			{Name: "--no-interactive", Description: "Disables interactive mode"},
			{Name: "--indexed-ram-bundle", Description: "Output usage information"},
			{Name: "--platforms", Description: "Scope linking to specified platforms"},
			{Name: "--all", Description: "Link all native modules and assets"},
			{Name: "--filename", Description: "Pulls the original Hermes tracing profile without any transformation"},
			{Name: "--sourcemap-path", Description: "The local path to your source map file, eg. /tmp/sourcemap.json"},
			{Name: "--generate-sourcemap", Description: "Generates the JS bundle and source map"},
			{Name: "--root", Description: "Do not launch packager while building"},
			{Name: "--terminal", Description: "Output usage information"},
			{Name: "--simulator", Description: "Explicitly set Xcode scheme to use"},
		},
	})
}
