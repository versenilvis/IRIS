package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "split",
		Description: "Use suffix_length letters to form the suffix of the file name",
		Options: []spec.Option{
			{Name: "-a", Description: "Use suffix_length letters to form the suffix of the file name"},
			{Name: "-b", Description: "N[K|k|M|m|G|g]"},
			{Name: "-d", Description: "Use a numeric suffix instead of a alphabetic suffix"},
			{Name: "-l", Description: "Create split files line_count lines in length"},
			{Name: "-p", Description: "The file to split"},
		},
	})
}
