package view

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "vim",
		Description: "Vi IMproved, a programmer",
		Options: []core.Option{
			{Name: "-v", Description: "Vi mode (like 'vi')"},
			{Name: "-e", Description: "Ex mode (like 'ex')"},
			{Name: "-E", Description: "Improved Ex mode"},
			{Name: "-s", Description: "Enable silent mode (when in ex mode), or Read Normal mode commands from file"},
			{Name: "-d", Description: "Diff mode (like 'vimdiff')"},
			{Name: "-y", Description: "Easy mode (like 'evim', modeless)"},
			{Name: "-R", Description: "Readonly mode (like 'view')"},
			{Name: "-Z", Description: "Restricted mode (like 'rvim')"},
			{Name: "-m", Description: "Modifications (writing files) not allowed"},
			{Name: "-M", Description: "Modifications in text not allowed"},
			{Name: "-b", Description: "Binary mode"},
			{Name: "-l", Description: "Lisp mode"},
			{Name: "-C", Description: "Compatible with Vi: 'compatible'"},
			{Name: "-N", Description: "Not fully Vi compatible: 'nocompatible'"},
			{Name: "-V", Description: "Be verbose [level N] [log messages to fname]"},
			{Name: "-D", Description: "Debugging mode"},
			{Name: "-n", Description: "No swap file, use memory only"},
			{Name: "-r", Description: "Same as -r"},
			{Name: "-T", Description: "Set terminal type to <terminal>"},
			{Name: "--not-a-term", Description: "Skip warning for input/output not being a terminal"},
			{Name: "--ttyfail", Description: "Exit if input or output is not a terminal"},
			{Name: "-u", Description: "Use <vimrc> instead of any .vimrc"},
			{Name: "--noplugin", Description: "Don't load plugin scripts"},
			{Name: "-p", Description: "Open N tab pages (default: one for each file)"},
			{Name: "-o", Description: "Open N windows (default: one for each file)"},
			{Name: "-O", Description: "Like -o but split vertically"},
			{Name: "--cmd", Description: "Execute <command> before loading any vimrc file"},
			{Name: "-c", Description: "Execute <command> after loading the first file"},
			{Name: "-S", Description: "Source file <session> after loading the first file"},
			{Name: "-w", Description: "Append all typed commands to file <scriptout>"},
			{Name: "-W", Description: "Write all typed commands to file <scriptout>"},
			{Name: "-x", Description: "Edit encrypted files"},
			{Name: "--startuptime", Description: "Write startup timing messages to <file>"},
			{Name: "-i", Description: "Use <viminfo> instead of .viminfo"},
			{Name: "--clean", Description: "'nocompatible', Vim defaults, no plugins, no viminfo"},
			{Name: "-h", Description: "Print Help message and exit"},
			{Name: "--version", Description: "Print version information and exit"},
		},
	})
}
