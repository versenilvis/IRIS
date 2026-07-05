package jvm

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "jenv",
		Description: "Executable file",
		Subcommands: []spec.Subcommand{
			{Name: "commands", Description: "List all available JEnv commands"},
			{Name: "help", Description: "Parses and displays help contents from a command's source file"},
			{Name: "info", Description: "Show information about which command will be executed"},
			{Name: "enable-plugin", Description: "Activate a jEnv plugin"},
			{Name: "pluginName", Description: "Plugin Name"},
			{Name: "disable-plugin", Description: "Deactivate a jEnv plugin"},
			{Name: "doctor", Description: "Run jEnv diagnostics"},
			{Name: "global", Description: "Sets the global Java version"},
			{Name: "global-options", Description: "Sets the global Java options"},
			{Name: "local-options", Description: "Sets the local application-specific Java options"},
			{Name: "shell", Description: "Sets a shell-specific Java version by setting the `JENV_VERSION'"},
			{Name: "shell-options", Description: "Sets the shell-specific Java options"},
			{Name: "hooks", Description: "List hook scripts for a given jenv command"},
			{Name: "init", Description: "Configure the shell environment for jenv"},
			{Name: "javahome", Description: "Display path to selected JAVA_HOME"},
			{Name: "options", Description: "Show the current Java options"},
			{Name: "options-file", Description: "Detect the file that sets the current jenv jvm options"},
			{Name: "options-file-read", Description: "Read options from file"},
			{Name: "options-file-write", Description: "Write options to a file"},
			{Name: "prefix", Description: "Displays the directory where a Java version is installed"},
			{Name: "refresh-plugins", Description: "Refresh plugins links"},
			{Name: "refresh-versions", Description: "Refresh alias names"},
			{Name: "rehash", Description: "Rehash jenv shims (run this after installing executables)"},
			{Name: "remove", Description: "Remove JDK installations"},
			{Name: "root", Description: "Display the root directory where versions and shims are kept"},
			{Name: "shims", Description: "List existing jenv shims"},
			{Name: "version", Description: "Shows the currently selected Java version and how it was selected"},
			{Name: "versions", Description: "Lists all Java versions found in `$JENV_ROOT/versions/*'"},
			{Name: "whence", Description: "List all Java versions that contain the given executable"},
		},
		Options: []spec.Option{
			{Name: "--usage", Description: "Show information about which command will be executed"},
			{Name: "--unset", Description: "Remove local jEnv settings"},
			{Name: "--short", Description: "Show only files without path"},
			{Name: "--bare", Description: "Display only version"},
			{Name: "--verbose", Description: "Display verbose output"},
			{Name: "--path", Description: "Show help for jEnv"},
			{Name: "--version", Description: "Show version for jEnv"},
		},
	})
}
