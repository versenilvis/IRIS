package jvm

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "gradle",
		Description: "Log all warnings",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Compute all outputs"},
			{Name: "run", Description: "Run applications"},
			{Name: "check", Description: "Run all checks"},
			{Name: "clean", Description: "Clear the contents of the build directory"},
			{Name: "projects", Description: "List of all sub-projects"},
			{Name: "tasks", Description: "List of main tasks of the selected project"},
			{Name: "help", Description: "Display task usage information"},
			{Name: "buildEnvironment", Description: "Visualises the buildscript dependencies of the selected project"},
			{Name: "properties", Description: "Gives you a list of the properties of the selected project"},
			{Name: "init", Description: "Create new Gradle builds, with new or existing projects"},
			{Name: "test", Description: "Run a test task"},
		},
		Options: []spec.Option{
			{Name: "-?", Description: "Shows a help message with all available CLI options"},
			{Name: "-v", Description: "Prints Gradle, Groovy, Ant, JVM, and operating system version information"},
			{Name: "-S", Description: "Print out the full (very verbose) stacktrace for any exceptions"},
			{Name: "-s", Description: "Print out the stacktrace also for user exceptions (e.g. compile error)"},
			{Name: "--scan", Description: "Debug Gradle Daemon process"},
			{Name: "--build-cache", Description: "Disables --parallel"},
			{Name: "--priority", Description: "Generate a build scan with detailed performance diagnostics"},
			{Name: "--watch-fs", Description: "Disables --daemon"},
			{Name: "--foreground", Description: "Starts the Gradle Daemon in a foreground process"},
			{Name: "-q", Description: "Log errors only"},
			{Name: "-w", Description: "Set log level to warn"},
			{Name: "-i", Description: "Set log level to info"},
			{Name: "-d", Description: "Log in debug mode (includes normal stacktrace)"},
			{Name: "--console", Description: "Specifies which type of console output to generate"},
			{Name: "--warning-mode", Description: "Specifies how to log warning"},
			{Name: "--include-build", Description: "Run the build as a composite, including the specified build"},
			{Name: "--offline", Description: "Specifies that the build should operate without accessing network resources"},
			{Name: "--refresh-dependencies", Description: "Refresh the state of dependencies"},
			{Name: "--dry-run", Description: "Specifies the build file. For example: gradle --build-file=foo.gradle"},
			{Name: "-c", Description: "Specifies the start directory for Gradle"},
			{Name: "--project-cache-dir", Description: "Sets a system property of the JVM, for example -Dmyprop=myvalue"},
			{Name: "-I", Description: "Specifies an initialization script"},
			{Name: "-P", Description: "Sets a project property of the root project, for example -Pmyprop=myvalue"},
			{Name: "--all", Description: "Display task usage information"},
			{Name: "--task", Description: "Visualises the buildscript dependencies of the selected project"},
			{Name: "--status", Description: "(Standalone command) Stop all Gradle Daemons of the same version"},
			{Name: "--type", Description: "Specify project type"},
			{Name: "--gradle-version", Description: "The Gradle version used for downloading and executing the Wrapper"},
			{Name: "--distribution-type", Description: "The Gradle distribution type used for the Wrapper"},
			{Name: "--gradle-distribution-url", Description: "The full URL pointing to Gradle distribution ZIP file"},
			{Name: "--continuous", Description: "This will force test and all task dependencies of test to execute"},
		},
	})
}
