package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "osqueryi",
		Description: "Your OS as a high-performance relational database",
		Options: []spec.Option{
			{Name: "--flagfile", Description: "Line-delimited file of additional flags"},
			{Name: "--D", Description: "Run as a daemon process"},
			{Name: "--S", Description: "Run as a shell process"},
			{Name: "--alarm_timeout", Description: "Seconds to allow for shutdown. Minimum is 10"},
			{Name: "--carver_block_size", Description: "Size of blocks used for POSTing data back to remote endpoints"},
			{Name: "--carver_compression", Description: "Compress archives using zstd prior to upload (default false)"},
			{Name: "--carver_continue_endpoint", Description: "TLS/HTTPS endpoint that receives carved content after session creation"},
			{Name: "--carver_disable_function", Description: "Disable the osquery file carver function (default true)"},
			{Name: "--carver_expiry", Description: "Seconds to store successful carve result metadata (in carves table)"},
			{Name: "--carver_start_endpoint", Description: "TLS/HTTPS init endpoint for forensic carver"},
			{Name: "--config_accelerated_refresh", Description: "Interval to wait if reading a configuration fails"},
			{Name: "--config_check", Description: "Check the format of an osquery config and exit"},
			{Name: "--config_dump", Description: "Dump the contents of the configuration, then exit"},
			{Name: "--config_enable_backup", Description: "Backup config and use it when refresh fails"},
			{Name: "--config_path", Description: "Path to JSON config file"},
			{Name: "--config_plugin", Description: "Config plugin name"},
			{Name: "--config_refresh", Description: "Optional interval in seconds to re-read configuration"},
			{Name: "--config_tls_endpoint", Description: "TLS/HTTPS endpoint for config retrieval"},
			{Name: "--config_tls_max_attempts", Description: "Number of attempts to retry a TLS config request"},
			{Name: "--daemonize", Description: "Attempt to daemonize (POSIX only)"},
			{Name: "--database_dump", Description: "Dump the contents of the backing store"},
			{Name: "--database_path", Description: "If using a disk-based backing store, specify a path"},
			{Name: "--disable_carver", Description: "Disable the osquery file carver (default true)"},
			{Name: "--disable_enrollment", Description: "Disable enrollment functions on related config/logger plugins"},
			{Name: "--disable_extensions", Description: "Disable extension API"},
			{Name: "--disable_reenrollment", Description: "Disable re-enrollment attempts if related plugins return invalid"},
			{Name: "--disable_tables", Description: "Comma-delimited list of table names to be disabled"},
			{Name: "--disable_watchdog", Description: "Disable userland watchdog process"},
			{Name: "--enable_extensions_watchdog", Description: "Enable userland watchdog for extensions processes"},
			{Name: "--enable_tables", Description: "Comma-delimited list of table names to be enabled"},
			{Name: "--enroll_always", Description: "On startup, send a new enrollment request"},
			{Name: "--enroll_secret_env", Description: "Name of environment variable holding enrollment-auth secret"},
			{Name: "--enroll_secret_path", Description: "Path to an optional client enrollment-auth secret"},
			{Name: "--enroll_tls_endpoint", Description: "TLS/HTTPS endpoint for client enrollment"},
			{Name: "--extensions_autoload", Description: "Optional path to a list of autoloaded & managed extensions"},
			{Name: "--extensions_interval", Description: "Seconds delay between connectivity checks"},
			{Name: "--extensions_require", Description: "Comma-separated list of required extensions"},
			{Name: "--extensions_socket", Description: "Path to the extensions UNIX domain socket"},
			{Name: "--extensions_timeout", Description: "Seconds to wait for autoloaded extensions"},
			{Name: "--force", Description: "Force osqueryd to kill previously-running daemons"},
		},
	})
}
