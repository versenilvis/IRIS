package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "xc",
		Description: "List tasks from an xc-compatible markdown file",
		Options: []core.Option{
			{Name: "-f", Description: "Print the markdown code of a task rather than running it"},
			{Name: "-H", Description: "List task names in a short format"},
			{Name: "-h", Description: "Print this help text"},
			{Name: "-complete", Description: "Install shell completion for xc"},
			{Name: "-uncomplete", Description: "Uninstall shell completion for xc"},
		},
	})
}
