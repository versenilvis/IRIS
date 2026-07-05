package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "grep",
		Description: "search text in files",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-R", Description: "recursive + follow symlinks"},
			{Name: "-i", Description: "ignore case"},
			{Name: "-n", Description: "show line numbers"},
			{Name: "-v", Description: "invert match"},
			{Name: "-l", Description: "show filenames only"},
			{Name: "-L", Description: "files without match"},
			{Name: "-c", Description: "count matching lines"},
			{Name: "-E", Description: "extended regex (egrep)"},
			{Name: "-P", Description: "perl-compatible regex"},
			{Name: "-F", Description: "fixed string (fgrep)"},
			{Name: "-w", Description: "whole word match"},
			{Name: "-x", Description: "whole line match"},
			{Name: "-A", Description: "lines after match"},
			{Name: "-B", Description: "lines before match"},
			{Name: "-C", Description: "lines around match"},
			{Name: "-m", Description: "max matches per file"},
			{Name: "-o", Description: "only matching part"},
			{Name: "-q", Description: "quiet mode"},
			{Name: "--include", Description: "include file pattern"},
			{Name: "--exclude", Description: "exclude file pattern"},
			{Name: "--exclude-dir", Description: "exclude directory"},
			{Name: "-H", Description: "always print filename"},
			{Name: "-h", Description: "never print filename"},
		},
	})

	// egrep and fgrep aliases
	spec.Register(&spec.Spec{
		Name:        "egrep",
		Description: "grep with extended regex",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-i", Description: "ignore case"},
			{Name: "-n", Description: "show line numbers"},
			{Name: "-v", Description: "invert match"},
			{Name: "-l", Description: "show filenames only"},
			{Name: "-w", Description: "whole word"},
			{Name: "-A", Description: "lines after"},
			{Name: "-B", Description: "lines before"},
			{Name: "-C", Description: "lines around"},
		},
	})
}
