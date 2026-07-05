package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "magento",
		Description: "Open-source E-commerce",
	})
}
