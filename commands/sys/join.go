package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "join",
		Description: "The join utility performs an",
		Options: []spec.Option{
			{Name: "-a", Description: "Replace empty output fields with string"},
			{Name: "-o", Description: "Join on the field'th field of file1"},
			{Name: "-2", Description: "Join on the field'th field of file2"},
			{Name: "-j", Description: "Join on the field'th field of both file1 and file2"},
		},
	})
}
