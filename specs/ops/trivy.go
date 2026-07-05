package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "trivy",
		Description: "Skip updating built-in policies [$TRIVY_SKIP_POLICY_UPDATE]",
		Subcommands: []core.Subcommand{
			{Name: "image", Description: "Scan an image"},
			{Name: "filesystem", Description: "Scan local filesystem for language-specific dependencies and config files"},
			{Name: "rootfs", Description: "Scan rootfs"},
			{Name: "sbom", Description: "Generate SBOM for an artifact"},
			{Name: "repository", Description: "Scan remote repository"},
			{Name: "client", Description: "Client mode"},
			{Name: "server", Description: "Server mode"},
			{Name: "config", Description: "Scan config files"},
			{Name: "plugin", Description: "Manage plugins"},
			{Name: "install", Description: "Install a plugin"},
			{Name: "uninstall", Description: "Uninstall plugin"},
			{Name: "list", Description: "List installed plugin"},
			{Name: "info", Description: "Information about a plugin"},
			{Name: "run", Description: "Run a plugin on the fly"},
			{Name: "update", Description: "Update an existing plugin"},
			{Name: "help", Description: "Shows a list of commands or help for one command"},
			{Name: "version", Description: "Print the version"},
		},
		Options: []core.Option{
			{Name: "--skip-policy-update", Description: "Skip updating built-in policies [$TRIVY_SKIP_POLICY_UPDATE]"},
			{Name: "--removed-pkgs", Description: "Input file path instead of image name [$TRIVY_INPUT]"},
			{Name: "--config-policy", Description: "Clear image caches without scanning [$TRIVY_CLEAR_CACHE]"},
			{Name: "--ignorefile", Description: "Timeout (default: 5m0s) [$TRIVY_TIMEOUT]"},
			{Name: "--offline-scan", Description: "Do not issue API requests to identify dependencies [$TRIVY_OFFLINE_SCAN]"},
			{Name: "--skip-files", Description: "Specify the file paths to skip traversal [$TRIVY_SKIP_FILES]"},
			{Name: "--skip-dirs", Description: "Allow insecure server connections when using SSL [$TRIVY_INSECURE]"},
			{Name: "--severity", Description: "Output file name [$TRIVY_OUTPUT]"},
			{Name: "--skip-db-update", Description: "Skip updating vulnerability database [$TRIVY_SKIP_UPDATE, $TRIVY_SKIP_DB_UPDATE]"},
			{Name: "--cache-backend", Description: "Output template [$TRIVY_TEMPLATE]"},
			{Name: "--format", Description: "Exit code when vulnerabilities were found (default: 0) [$TRIVY_EXIT_CODE]"},
			{Name: "--vuln-type", Description: "Specify the Rego file to evaluate each vulnerability [$TRIVY_IGNORE_POLICY]"},
			{Name: "--list-all-pkgs", Description: "Suppress progress bar [$TRIVY_NO_PROGRESS]"},
			{Name: "--token", Description: "For authentication in client/server mode [$TRIVY_TOKEN]"},
			{Name: "--token-header", Description: "Custom headers in client/server mode [$TRIVY_CUSTOM_HEADERS]"},
			{Name: "--insecure", Description: "Allow insecure server connections when using SSL [$TRIVY_INSECURE]"},
			{Name: "--listen", Description: "Display only fixed vulnerabilities [$TRIVY_IGNORE_UNFIXED]"},
			{Name: "--remote", Description: "Remove all caches and database [$TRIVY_RESET]"},
			{Name: "--light", Description: "Deprecated [$TRIVY_LIGHT]"},
			{Name: "--server", Description: "Server address [$TRIVY_SERVER]"},
			{Name: "--artifact-type", Description: "Scan remote repository"},
			{Name: "--quiet", Description: "Suppress progress bar and log output [$TRIVY_QUIET]"},
			{Name: "--file-patterns", Description: "Specify file patterns [$TRIVY_FILE_PATTERNS"},
			{Name: "--include-non-failures", Description: "Enable more verbose trace output for custom queries [$TRIVY_TRACE]"},
			{Name: "--trace", Description: "Enable more verbose trace output for custom queries [$TRIVY_TRACE]"},
			{Name: "--debug", Description: "Enable debug output [$TRIVY_DEBUG]"},
			{Name: "--cache-dir", Description: "Cache directory [$TRIVY_CACHE_DIR]"},
			{Name: "--help", Description: "Show help"},
			{Name: "--version", Description: "Print the version"},
		},
	})
}
