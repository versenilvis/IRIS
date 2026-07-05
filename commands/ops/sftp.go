package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "sftp",
		Description: "OpenSSH secure file transfer",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for sftp"},
			{Name: "-4", Description: "Forces scp to use IPv4 addresses only"},
			{Name: "-6", Description: "Forces scp to use IPv6 addresses only"},
			{Name: "-A", Description: "The buffer size"},
			{Name: "-b", Description: "The batch file"},
			{Name: "-C", Description: "Compression enable. Passes the -C flag to ssh(1) to enable compression"},
			{Name: "-c", Description: "The selected cipher specification"},
			{Name: "-D", Description: "Path to the SFTP server"},
			{Name: "-F", Description: "The selected ssh config"},
			{Name: "-f", Description: "Requests that files be flushed to disk immediately after transfer"},
			{Name: "-i", Description: "Specified identity file"},
			{Name: "-J", Description: "Scp destination"},
			{Name: "-l", Description: "Limits the used bandwidth, specified in Kbit/s"},
			{Name: "-N", Description: "Disables quiet mode, e.g. to override the implicit quiet mode set by the -b flag"},
			{Name: "-o", Description: "Preserves modification times, access times, and modes from the original file"},
			{Name: "-q", Description: "The number of requests"},
			{Name: "-r", Description: "Path to the SFTP server"},
		},
	})
}
