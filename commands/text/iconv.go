package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "iconv",
		Description: "Character set conversion",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for iconv"},
			{Name: "--version", Description: "Output version information and exit"},
			{Name: "-f", Description: "Specifies the encoding of the input"},
			{Name: "-t", Description: "Specifies the encoding of the output"},
			{Name: "-c", Description: "Discard unconvertible characters"},
			{Name: "-l", Description: "List the supported encodings"},
			{Name: "--unicode-subst", Description: "Substitution for unconvertible Unicode characters"},
			{Name: "--byte-subst", Description: "Substitution for unconvertible bytes"},
			{Name: "--widechar-subst", Description: "Substitution for unconvertible wide characters"},
		},
	})
}
