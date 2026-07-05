package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "crontab",
		Description: "Maintain crontab file for individual users",
		Options: []spec.Option{
			{Name: "-e", Description: "Edit the current crontab"},
			{Name: "-l", Description: "Display the current crontab"},
			{Name: "-r", Description: "Remove the current crontab"},
			{Name: "-u", Description: "Specify the name of the user whose crontab is to be tweaked"},
		},
	})
}
