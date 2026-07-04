package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "dirname",
		Description: "Return directory portion of pathname",
	})
}
