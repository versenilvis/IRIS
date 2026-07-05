package jvm

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mvn",
		Description: "Maven - a Java based project management and comprehension tool",
		Options: []spec.Option{
			{Name: "--also-make", Description: "Also build projects required by project list"},
			{Name: "--also-make-dependents", Description: "Also build projects that depend on projects in the project list"},
			{Name: "--batch-mode", Description: "Run in non-interactive (batch)"},
			{Name: "--builder", Description: "Specify the build strategy to use"},
			{Name: "--strict-checksums", Description: "Fail if checksums do not match"},
			{Name: "--lax-checksums", Description: "Warn if checksums do not match"},
			{Name: "--color", Description: "Specify the color mode of the output"},
			{Name: "--check-plugin-updates", Description: "Ineffective. Only kept for backward compatibility"},
			{Name: "--define", Description: "Define a system property"},
			{Name: "--errors", Description: "Produce execution error messages"},
			{Name: "--encrypt-master-password", Description: "Encrypt the master security password"},
			{Name: "--encrypt-password", Description: "Encrypt the server password"},
			{Name: "--file", Description: "Force the use of an alternate POM file (or directory with pom.xml)"},
			{Name: "--fail-at-end", Description: "Only fail the build afterwards; allow all non-impacted builds to continue"},
			{Name: "--fail-fast", Description: "Stop at first failure in reactorized builds"},
			{Name: "--fail-never", Description: "Never fail the build, regardless of project result"},
			{Name: "--global-settings", Description: "Specify the global settings file to use"},
			{Name: "--global-toolchains", Description: "Specify the global toolchains file to use"},
			{Name: "--help", Description: "Display help information"},
			{Name: "--log-file", Description: "Specify the file to log to"},
			{Name: "--legacy-local-repository", Description: "Use the Maven2 legacy local repository behaviour"},
			{Name: "--non-recursive", Description: "Do not recurse into sub-projects"},
			{Name: "--no-plugin-registry", Description: "Ineffective. Only kept for backward compatibility"},
			{Name: "--no-plugin-updates", Description: "Ineffective. Only kept for backward compatibility"},
			{Name: "--no-snapshot-updates", Description: "Suppress SNAPSHOT updates"},
			{Name: "--no-transfer-progress", Description: "Do not display transfer progress when downloading or uploading"},
			{Name: "--offline", Description: "Work offline"},
			{Name: "--activate-profiles", Description: "Activate the specified profiles (comma delimited)"},
			{Name: "--projects", Description: "Specify the projects to build"},
			{Name: "--quiet", Description: "Quiet output - only shows errors"},
			{Name: "--resume-from", Description: "Resume from the specified project"},
			{Name: "--settings", Description: "Specify the user settings file to use"},
			{Name: "--toolchains", Description: "Specify the toolchains file to use"},
			{Name: "--threads", Description: "Specify the number of threads to use"},
			{Name: "--update-snapshots", Description: "Forces a check for missing releases and updated snapshots on remote repositories"},
			{Name: "--update-plugins", Description: "Ineffective. Only kept for backward compatibility"},
			{Name: "--version", Description: "Display version information"},
			{Name: "--show-version", Description: "Display version information"},
			{Name: "--debug", Description: "Produce execution debug output"},
		},
	})
}
