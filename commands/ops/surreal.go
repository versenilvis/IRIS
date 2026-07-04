package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "surreal",
		Description: "Database authentication password to use when connecting [default: root]",
		Subcommands: []core.Subcommand{
			{Name: "help", Description: "Print this message or the help of the given subcommand(s)"},
			{Name: "start", Description: "Start the database server"},
			{Name: "path", Description: "Database path used for storing data [env: DB_PATH=] [default: memory]"},
			{Name: "key", Description: "Encryption key to use for on-disk encryption [env: KEY=]"},
			{Name: "kvs-ca", Description: "Path to the CA file used when connecting to the remote KV store [env: KVS_CA=]"},
			{Name: "pass", Description: "The master password for the database [env: PASS=]"},
			{Name: "strict", Description: "Whether strict mode is enabled on this database instance [env: STRICT=]"},
			{Name: "user", Description: "The master username for the database [env: USER=] [default: root]"},
			{Name: "web-crt", Description: "Path to the certificate file for encrypted client connections [env: WEB_CRT=]"},
			{Name: "web-key", Description: "Path to the private key file for encrypted client connections [env: WEB_KEY=]"},
			{Name: "backup", Description: "Backup data to or from an existing database"},
			{Name: "from", Description: "Path to the remote database or file from which to export"},
			{Name: "into", Description: "Path to the remote database or file into which to import"},
			{Name: "import", Description: "Import a SurrealQL script into an existing database"},
			{Name: "export", Description: "Export an existing database as a SurrealQL script"},
			{Name: "version", Description: "Output the command-line tool version information"},
			{Name: "sql", Description: "Start an SQL REPL in your terminal with pipe support"},
		},
		Options: []core.Option{
			{Name: "--pass", Description: "Database authentication password to use when connecting [default: root]"},
			{Name: "--user", Description: "Database authentication username to use when connecting [default: root]"},
			{Name: "--conn", Description: "Remote database server url to connect to [default: https://cloud.surrealdb.com]"},
			{Name: "--ns", Description: "Print this message or the help of the given subcommand(s)"},
			{Name: "--addr", Description: "Encryption key to use for on-disk encryption [env: KEY=]"},
			{Name: "--kvs-ca", Description: "Path to the CA file used when connecting to the remote KV store [env: KVS_CA=]"},
			{Name: "--kvs-crt", Description: "The master password for the database [env: PASS=]"},
			{Name: "--strict", Description: "Whether strict mode is enabled on this database instance [env: STRICT=]"},
			{Name: "--web-crt", Description: "Path to the certificate file for encrypted client connections [env: WEB_CRT=]"},
			{Name: "--web-key", Description: "Path to the private key file for encrypted client connections [env: WEB_KEY=]"},
			{Name: "--pretty", Description: "Whether database responses should be pretty printed"},
			{Name: "--help", Description: "Print help information"},
		},
	})
}
