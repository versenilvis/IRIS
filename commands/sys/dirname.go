package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dirname",
		Description: "Return directory portion of pathname",
	})
}
