package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "lvim",
		Description: "Hyperextensible Vim-based text editor",
		Options: []core.Option{
			{Name: "--cmd", Description: "Execute <cmd> before any config"},
			{Name: "-c", Description: "Execute <cmd> after config and first file"},
			{Name: "-b", Description: "Binary mode"},
			{Name: "-d", Description: "Diff mode"},
			{Name: "-e", Description: "Ex mode"},
			{Name: "-es", Description: "Silent (batch) mode"},
			{Name: "-h", Description: "Print this help message"},
			{Name: "-i", Description: "Use this shada file"},
			{Name: "-m", Description: "Modifications (writing files) not allowed"},
			{Name: "-M", Description: "Modifications in text not allowed"},
			{Name: "-n", Description: "No swap file, use memory only"},
			{Name: "-o", Description: "Open N windows (default: one per file)"},
			{Name: "-O", Description: "Open N vertical windows (default: one per file)"},
			{Name: "-p", Description: "Open N tab pages (default: one per file)"},
			{Name: "-L", Description: "List swap files"},
			{Name: "-r", Description: "Recover edit state for this file"},
			{Name: "-R", Description: "Read-only mode"},
			{Name: "-S", Description: "Source <session> after loading the first file"},
			{Name: "-s", Description: "Read Normal mode commands from <scriptin>"},
			{Name: "-u", Description: "Use this config file"},
			{Name: "-v", Description: "Print version information"},
			{Name: "--api-info", Description: "Write msgpack-encoded API metadata to stdout"},
			{Name: "--embed", Description: "Use stdin/stdout as a msgpack-rpc channel"},
			{Name: "--headless", Description: "Don't start a user interface"},
			{Name: "--listen", Description: "Serve RPC API from this address"},
			{Name: "--noplugin", Description: "Don't load plugins"},
			{Name: "--remote", Description: "Execute commands remotely on a server"},
			{Name: "--server", Description: "Specify RPC server to send commands to"},
			{Name: "--startuptime", Description: "Write startup timing messages to <file>"},
		},
	})
}
