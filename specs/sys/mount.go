package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mount",
		Description: "Mount disks and manage subtrees",
		Options: []core.Option{
			{Name: "-h", Description: "Help for abc"},
			{Name: "-a", Description: "Mount all filesystems in fstab"},
			{Name: "-c", Description: "Don't canonicalize paths"},
			{Name: "-f", Description: "Dry run; skip the mount(2) syscall"},
			{Name: "-F", Description: "Fork off for each device (use with -a)"},
			{Name: "-T", Description: "Alternative file to /etc/fstab"},
			{Name: "-i", Description: "Don't call the mount.<type> helpers"},
			{Name: "-l", Description: "Show also filesystem labels"},
			{Name: "-m", Description: "Alias to '-o X-mount.mkdir"},
			{Name: "-n", Description: "Don't write to /etc/mtab"},
			{Name: "--options-mode", Description: "What to do with options loaded from fstab"},
			{Name: "--options-source", Description: "Mount options source"},
			{Name: "--options-source-force", Description: "Force use of options from fstab/mtab"},
			{Name: "-o", Description: "Comma-separated list of mount options"},
			{Name: "-O", Description: "Limit the set of filesystems (use with -a)"},
			{Name: "-r", Description: "Mount the filesystem read-only (same as -o ro)"},
			{Name: "-t", Description: "Limit the set of filesystem types"},
			{Name: "--source", Description: "Explicitly specifies source"},
			{Name: "--target", Description: "Explicitly specifies mountpoint"},
			{Name: "--target-prefix", Description: "Specifies path used for all mountpoints"},
			{Name: "-v", Description: "Say what is being done"},
			{Name: "-w", Description: "Mount the filesystem read-write (default)"},
			{Name: "-V", Description: "Display version"},
			{Name: "-B", Description: "Mount a subtree somewhere else (same as -o bind)"},
			{Name: "-M", Description: "Move a subtree to some other place"},
			{Name: "-R", Description: "Mount a subtree and all submounts somewhere else"},
			{Name: "--make-shared", Description: "Mark a subtree as shared"},
			{Name: "--make-slave", Description: "Mark a subtree as slave"},
			{Name: "--make-private", Description: "Mark a subtree as private"},
			{Name: "--make-unbindable", Description: "Mark a subtree as unbindable"},
			{Name: "--make-rshared", Description: "Recursively mark a whole subtree as shared"},
			{Name: "--make-rslave", Description: "Recursively mark a whole subtree as slave"},
			{Name: "--make-rprivate", Description: "Recursively mark a whole subtree as private"},
			{Name: "--make-runbindable", Description: "Recursively mark a whole subtree as unbindable"},
		},
	})
}
