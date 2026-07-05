package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rscript",
		Description: "Scripting Front-End for R",
		Options: []spec.Option{
			{Name: "-e", Description: "R expression to run"},
			{Name: "--help", Description: "Print usage and exit"},
			{Name: "--version", Description: "Print version and exit"},
			{Name: "--verbose", Description: "Print information on progress"},
			{Name: "--no-echo", Description: "Run as quietly as possible"},
			{Name: "--no-restore", Description: "Don't restore anything"},
			{Name: "--save", Description: "Do save workspace at the end of the session"},
			{Name: "--no-environ", Description: "Don't read the site and user environment files"},
			{Name: "--no-site-file", Description: "Don't read the site-wide Rprofile"},
			{Name: "--no-init-file", Description: "Don't read the user R profile"},
			{Name: "--restore", Description: "Do restore previously saved objects at startup"},
			{Name: "--vanilla", Description: "Combine --no-save, --no-restore, --no-site-file --no-init-file and --no-environ"},
			{Name: "--default-packages", Description: "Comma separated list of default packages"},
		},
	})
}
