package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pg_dump",
		Description: "Dumps a database as a text file or to other formats",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for pg_dump"},
			{Name: "--file", Description: "Output file or directory name"},
			{Name: "--format", Description: "Output file format"},
			{Name: "--jobs", Description: "Number of parallel job to dump"},
			{Name: "--verbose", Description: "Verbose mode"},
			{Name: "--version", Description: "Output version information"},
			{Name: "--compress", Description: "Compression level for compressed formats"},
			{Name: "--lock-wait-timeout", Description: "Fail after waiting <timeout> for a table lock"},
			{Name: "--no-sync", Description: "Do not wait for changes to be written safely to disk"},
			{Name: "--data-only", Description: "Dump only the data, not the schema"},
			{Name: "--blobs", Description: "Include large objects in dump"},
			{Name: "--no-blobs", Description: "Exclude large objects in dump"},
			{Name: "--clean", Description: "Clean (drop) database objects before recreating"},
			{Name: "--create", Description: "Include commands to create database in dump"},
			{Name: "--extension", Description: "Dump the specified extension(s) only"},
			{Name: "--encoding", Description: "Dump the data in encoding <encoding>"},
			{Name: "--schema", Description: "Dump the specified schema(s) only"},
			{Name: "--exclude-schema", Description: "Do NOT dump the specified schema(s)"},
			{Name: "--no-owner", Description: "Skip restoration of object ownership in plain-text format"},
			{Name: "--schema-only", Description: "Dump only the schema, no data"},
			{Name: "--superuser", Description: "Superuser user name to use in plain-text format"},
			{Name: "--table", Description: "Dump the specified table(s) only"},
			{Name: "--exclude-table", Description: "Do NOT dump the specified table(s)"},
			{Name: "--no-privileges", Description: "Do not dump privileges (grant/revoke)"},
			{Name: "--binary-upgrade", Description: "For use by upgrade utilities only"},
			{Name: "--column-inserts", Description: "Dump data as INSERT commands with column names"},
			{Name: "--disable-dollar-quoting", Description: "Disable dollar quoting, use SQL standard quoting"},
			{Name: "--disable-triggers", Description: "Disable triggers during data-only restore"},
			{Name: "--enable-row-security", Description: "Enable row security (dump only content user has access to)"},
			{Name: "--exclude-table-data", Description: "Do NOT dump data for the specified table(s)"},
			{Name: "--extra-float-digits", Description: "Override default setting for extra_float_digits"},
			{Name: "--if-exists", Description: "Use IF EXISTS when dropping objects"},
			{Name: "--include-foreign-data", Description: "Include data of foreign tables on foreign servers matching PATTERN"},
			{Name: "--inserts", Description: "Dump data as INSERT commands, rather than COPY"},
			{Name: "--load-via-partition-root", Description: "Load partitions via the root table"},
			{Name: "--no-comments", Description: "Do not dump comments"},
			{Name: "--no-publications", Description: "Do not dump publications"},
			{Name: "--no-security-labels", Description: "Do not dump security label assignments"},
			{Name: "--no-subscriptions", Description: "Do not dump subscriptions"},
			{Name: "--no-synchronized-snapshots", Description: "Do not use synchronized snapshots in parallel jobs"},
		},
	})
}
