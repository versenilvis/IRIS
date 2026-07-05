package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "magento",
		Description: "Open-source E-commerce",
	})
}
