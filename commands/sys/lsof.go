package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "lsof",
		Description: "List open files",
		Options: []spec.Option{
			{Name: "-?", Description: "Help"},
			{Name: "-a", Description: "Apply AND to the selections (defaults to OR)"},
			{Name: "-b", Description: "Avoid kernel blocks"},
			{Name: "-c", Description: "Select the listing of files for processes executing a command"},
			{Name: "-d", Description: "Search tree for all open instances/files/directories of directory. *SLOW?*"},
			{Name: "-f", Description: "Inhibit path name arguments to be interpreted"},
			{Name: "-F", Description: "Select fields to output"},
			{Name: "-g", Description: "Exclude or select by process group IDs (PGID)"},
			{Name: "-i", Description: "Selects files by [46][protocol][@hostname|hostaddr][:service|port]"},
			{Name: "-l", Description: "Inhibit conversion of user IDs to login names"},
			{Name: "-L", Description: "Disable listing of file link counts"},
			{Name: "-M", Description: "Disable portMap registration"},
			{Name: "-n", Description: "No host names"},
			{Name: "-N", Description: "Select NFS files"},
			{Name: "-o", Description: "List file offset"},
			{Name: "-O", Description: "No overhead *RISKY*"},
			{Name: "-p", Description: "Exclude or select process identification numbers (PIDs)"},
			{Name: "-P", Description: "No port names"},
			{Name: "-r", Description: "Repeat every t seconds (15) forever"},
			{Name: "-R", Description: "List parent PID"},
			{Name: "-s", Description: "List file size or exclude/select protocol"},
			{Name: "-S", Description: "Stat timeout in seconds (lstat/readlink/stat)"},
			{Name: "-T", Description: "Disable TCP/TPI info"},
			{Name: "-t", Description: "Specify terse listing"},
			{Name: "-u", Description: "Exclude/select login|UID set"},
			{Name: "-U", Description: "Select Unix socket"},
			{Name: "-v", Description: "List version info"},
			{Name: "-V", Description: "Verbose search"},
			{Name: "-w", Description: "Disable warnings"},
			{Name: "-x", Description: "Cross over +d|+D File systems or symbolic links"},
			{Name: "-X", Description: "File descriptor table only"},
		},
	})
}
