package golang

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "goreleaser",
		Description: "Deliver Go binaries as fast and easily as possible",
		Subcommands: []spec.Subcommand{
			{Name: "b", Description: "Builds the current project"},
			{Name: "c", Description: "Checks if configuration is valid"},
			{Name: "completion", Description: "Generate the autocompletion script for the specified shell"},
			{Name: "bash", Description: "Generate the autocompletion script for bash"},
			{Name: "fish", Description: "Generate the autocompletion script for fish"},
			{Name: "powershell", Description: "Generate the autocompletion script for powershell"},
			{Name: "zsh", Description: "Generate the autocompletion script for zsh"},
			{Name: "hc", Description: "Checks if needed tools are installed"},
			{Name: "i", Description: "Generates a .goreleaser.yaml file"},
			{Name: "schema", Description: "Outputs goreleaser's JSON schema"},
			{Name: "r", Description: "Releases the current project"},
			{Name: "help", Description: "Help about any command"},
		},
		Options: []spec.Option{
			{Name: "--clean", Description: "Remove the dist folder before building"},
			{Name: "--config", Description: "Load configuration from file"},
			{Name: "--deprecated", Description: "Force print the deprecation message - tests only"},
			{Name: "--id", Description: "Builds only the specified build ids"},
			{Name: "--output", Description: "Amount tasks to run concurrently (default: number of CPUs)"},
			{Name: "--rm-dist", Description: "Remove the dist folder before building"},
			{Name: "--single-target", Description: "Skips global before hooks"},
			{Name: "--skip-post-hooks", Description: "Skips all post-build hooks"},
			{Name: "--skip-validate", Description: "Skips several sanity checks"},
			{Name: "--snapshot", Description: "Generate an unversioned snapshot build, skipping all validations"},
			{Name: "--timeout", Description: "Timeout to the entire build process"},
			{Name: "--quiet", Description: "Quiet mode: no output"},
			{Name: "--no-descriptions", Description: "Disable completion descriptions"},
			{Name: "--auto-snapshot", Description: "Automatically sets --snapshot if the repository is dirty"},
			{Name: "--fail-fast", Description: "Whether to abort the release publishing on the first error"},
			{Name: "--parallelism", Description: "Amount tasks to run concurrently (default: number of CPUs)"},
			{Name: "--release-footer", Description: "Load custom release notes footer from a markdown file"},
			{Name: "--release-footer-tmpl", Description: "Load custom release notes header from a markdown file"},
			{Name: "--release-header-tmpl", Description: "Removes the dist folder"},
			{Name: "--skip", Description: "Skips announcing releases (implies --skip=validate)"},
			{Name: "--skip-before", Description: "Skips global before hooks"},
			{Name: "--skip-docker", Description: "Skips Docker Images/Manifests builds"},
			{Name: "--skip-ko", Description: "Skips Ko builds"},
			{Name: "--skip-publish", Description: "Skips publishing artifacts (implies --skip=announce)"},
			{Name: "--skip-sbom", Description: "Skips cataloging artifacts"},
			{Name: "--skip-sign", Description: "Skips signing artifacts"},
			{Name: "--debug", Description: "Enable verbose mode"},
			{Name: "--verbose", Description: "Enable verbose mode"},
			{Name: "--help", Description: "Display help"},
		},
	})
}
