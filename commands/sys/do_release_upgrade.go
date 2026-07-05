package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "do-release-upgrade",
		Description: "Upgrade Ubuntu to latest release",
		Options: []spec.Option{
			{Name: "-h", Description: "Show help message and exit"},
			{Name: "-d", Description: "If using the latest supported release, upgrade to the development release"},
			{Name: "-p", Description: "Try upgrading to the latest release using the upgrader from Ubuntu-proposed"},
			{Name: "-m", Description: "Run in a special upgrade mode"},
			{Name: "--mode", Description: "Run in a special upgrade mode"},
			{Name: "-f", Description: "Run the specified frontend"},
			{Name: "--frontend", Description: "Run the specified frontend"},
		},
	})
}
