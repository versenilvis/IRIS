package cc

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "bazel",
		Description: "Bazel target",
		Subcommands: []core.Subcommand{
			{Name: "run", Description: "Runs the specified target"},
			{Name: "test", Description: "Builds and runs the specified test targets"},
			{Name: "build", Description: "Builds the specified targets"},
		},
		Options: []core.Option{
			{Name: "--autodetect_server_javabase", Description: "Back to the local JDK for running the bazel server and instead exits"},
			{Name: "--nobatch", Description: "Run with a server"},
			{Name: "--batch_cpu_scheduling", Description: "Only on Linux; use 'batch' CPU scheduling for Bazel"},
			{Name: "--nobatch_cpu_scheduling", Description: "Only on Linux; Bazel does not perform a system call"},
			{Name: "--bazelrc", Description: "Wait for a running command to complete"},
			{Name: "--noblock_for_lock", Description: "Don't log debug information from the client to stderr"},
			{Name: "--connect_timeout_secs", Description: "The amount of time the client waits for each attempt to connect to the server"},
			{Name: "--expand_configs_in_place", Description: "Changed the expansion of --config flags to be done in-place"},
			{Name: "--noexpand_configs_in_place", Description: "Look for the home bazelrc file at $HOME/.bazelrc"},
			{Name: "--nohome_rc", Description: "Don't look for the home bazelrc file at $HOME/.bazelrc"},
			{Name: "--idle_server_tasks", Description: "Run System.gc() when the server is idle"},
			{Name: "--noidle_server_tasks", Description: "Don't run System.gc() when the server is idle"},
			{Name: "--ignore_all_rc_files", Description: "Enables all rc files"},
			{Name: "--io_nice_level", Description: "The maximum amount of time the client waits to connect to the server"},
			{Name: "--macos_qos_class", Description: "Sets the QoS service class of the bazel server when running on macOS"},
			{Name: "--max_idle_secs", Description: "The number of seconds the build server will wait idling before shutting down"},
			{Name: "--output_base", Description: "Specifies the output location to which all build output will be written"},
			{Name: "--output_base_root", Description: "The user-specific directory beneath which all build outputs are written"},
			{Name: "--preemptible", Description: "If true, the command can be preempted if another command is started"},
			{Name: "--nopreemptible", Description: "If true, the command can be preempted if another command is started"},
			{Name: "--server_jvm_out", Description: "The location to write the server's JVM's output"},
			{Name: "--shutdown_on_low_sys_mem", Description: "Linux only. Don't shut down the server when the system is low on free RAM"},
			{Name: "--system_rc", Description: "Look for the system-wide bazelrc"},
			{Name: "--nosystem_rc", Description: "Don't look for the system-wide bazelrc"},
			{Name: "--unlimit_coredumps", Description: "Scan every file for a change"},
			{Name: "--windows_enable_symlinks", Description: "Real symbolic links will be created on Windows instead of file copying"},
			{Name: "--nowindows_enable_symlinks", Description: "Real symbolic links will be created via file copying"},
			{Name: "--workspace_rc", Description: "Look for the workspace bazelrc file at $workspace/.bazelrc"},
			{Name: "--noworkspace_rc", Description: "Don't look for the workspace bazelrc file at $workspace/.bazelrc"},
		},
	})
}
