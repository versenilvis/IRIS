package fs

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "trash",
		Description: "Trash, move files/folders to the trash",
		Options: []core.Option{
			{Name: "-v", Description: "Print verbose output while moving items"},
			{Name: "-F", Description: "Use the Finder API to move items to the trash"},
			{Name: "-l", Description: "List items in the trash"},
			{Name: "-e", Description: "Empty the trash"},
			{Name: "-s", Description: "Skips the confirmation prompt for -e and -s. CAUTION: permanently instantly"},
		},
	})
}
