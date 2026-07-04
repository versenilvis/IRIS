package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kdoctor",
		Description: "Report a version of KDoctor",
		Options: []core.Option{
			{Name: "--version", Description: "Report a version of KDoctor"},
			{Name: "--verbose", Description: "Report an extended information"},
			{Name: "--all", Description: "Run extra diagnostics"},
			{Name: "--team-ids", Description: "Report all available Apple dev team ids"},
			{Name: "--help", Description: "Usage info"},
		},
	})
}
