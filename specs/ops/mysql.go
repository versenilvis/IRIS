package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mysql",
		Description: "Mysql is a terminal-based front-end to MySQL",
		Options: []core.Option{
			{Name: "--auto-rehash", Description: "Enable automatic rehashing"},
			{Name: "--auto-vertical-output", Description: "Enable automatic vertical result set display"},
			{Name: "--batch", Description: "Do not use history file"},
			{Name: "--binary-as-hex", Description: "Display binary values in hexadecimal notation"},
			{Name: "--binary-mode", Description: "Disable \\\\r\\\\n - to - \\\\n translation and treatment of \\\\0 as end-of-query"},
			{Name: "--bind-address", Description: "Use specified network interface to connect to MySQL Server"},
			{Name: "--character-sets-dir", Description: "Directory where character sets are installed"},
			{Name: "--column-names", Description: "Write column names in results"},
			{Name: "--column-type-info", Description: "Display result set metadata"},
			{Name: "--comments", Description: "Whether to retain or strip comments in statements sent to the server"},
			{Name: "--compress", Description: "Compress all information sent between client and server"},
			{Name: "--compression-algorithms", Description: "Permitted compression algorithms for connections to server"},
			{Name: "--connect-expired-password", Description: "Indicate to server that client can handle expired-password sandbox mode"},
			{Name: "--connect-timeout", Description: "Number of seconds before connection timeout"},
			{Name: "-D", Description: "The database to use"},
			{Name: "--debug", Description: "Write debugging log; supported only if MySQL was built with debugging support"},
			{Name: "--debug-check", Description: "Print debugging information when program exits"},
			{Name: "-T", Description: "Print debugging information, memory, and CPU statistics when program exits"},
			{Name: "--default-auth", Description: "Authentication plugin to use"},
			{Name: "--default-character-set", Description: "Specify default character set"},
			{Name: "--defaults-extra-file", Description: "Read named option file in addition to usual option files"},
			{Name: "--defaults-file", Description: "Read only named option file"},
			{Name: "--defaults-group-suffix", Description: "Option group suffix value"},
			{Name: "--delimiter", Description: "Set the statement delimiter"},
			{Name: "--disable-named-commands", Description: "Use DNS SRV lookup for host information"},
			{Name: "--enable-cleartext-plugin", Description: "Enable cleartext authentication plugin"},
			{Name: "-e", Description: "Execute the statement and quit"},
			{Name: "-f", Description: "Continue even if an SQL error occurs"},
			{Name: "--get-server-public-key", Description: "Request RSA public key from server"},
			{Name: "--help", Description: "Display help message and exit"},
			{Name: "--histignore", Description: "Patterns specifying which statements to ignore for logging"},
			{Name: "-h", Description: "Host on which MySQL server is located"},
			{Name: "-H", Description: "Produce HTML output"},
			{Name: "-i", Description: "Ignore spaces after function names"},
			{Name: "--init-command", Description: "SQL statement to execute after connecting"},
			{Name: "--line-numbers", Description: "Write line numbers for errors"},
			{Name: "--load-data-local-dir", Description: "Directory for files named in LOAD DATA LOCAL statements"},
			{Name: "--local-infile", Description: "Enable or disable for LOCAL capability for LOAD DATA"},
			{Name: "--login-path", Description: "Read login path options from .mylogin.cnf"},
			{Name: "--max-allowed-packet", Description: "Maximum packet length to send to or receive from server"},
		},
	})
}
