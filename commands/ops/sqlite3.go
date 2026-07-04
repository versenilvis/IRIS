package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "sqlite3",
		Description: "A command line interface for SQLite version 3",
		Options: []core.Option{
			{Name: "-append", Description: "Append the database to the end of the file"},
			{Name: "-ascii", Description: "Set output mode to 'ascii'"},
			{Name: "-bail", Description: "Stop after hitting an error"},
			{Name: "-batch", Description: "Force batch I/O"},
			{Name: "-column", Description: "Set output mode to 'column'"},
			{Name: "-cmd", Description: "Set output mode to 'csv'"},
			{Name: "-echo", Description: "Print commands before execution"},
			{Name: "-init", Description: "Read/process named file"},
			{Name: "-header", Description: "Turn headers on"},
			{Name: "-noheader", Description: "Turn headers off"},
			{Name: "-help", Description: "Show help message"},
			{Name: "-html", Description: "Set output mode to HTML"},
			{Name: "-interactive", Description: "Force interactive I/O"},
			{Name: "-line", Description: "Set output mode to 'line'"},
			{Name: "-list", Description: "Set output mode to 'list'"},
			{Name: "-lookaside", Description: "Use N entries of SZ bytes for lookaside memory"},
			{Name: "-memtrace", Description: "Trace all memory allocations and deallocations"},
			{Name: "-mmap", Description: "Default mmap size set to N"},
			{Name: "-newline", Description: "Set output row separator"},
			{Name: "-nofollow", Description: "Refuse to open symbolic links to database files"},
			{Name: "-nullvalue", Description: "Set text string for NULL values"},
			{Name: "-pagecache", Description: "Use N slots of SZ bytes each for page cache memory"},
			{Name: "-quote", Description: "Set output mode to 'quote'"},
			{Name: "-readonly", Description: "Open the database read-only"},
			{Name: "-separator", Description: "Set output column separator"},
			{Name: "-stats", Description: "Print memory stats before each finalize"},
			{Name: "-version", Description: "Show SQLite version"},
			{Name: "-vfs", Description: "Use NAME as the default VFS"},
		},
	})
}
