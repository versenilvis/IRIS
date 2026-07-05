package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "sha1sum",
		Description: "Print or check SHA1 (160-bit) checksums",
		Options: []spec.Option{
			{Name: "-b", Description: "Read in binary mode"},
			{Name: "-c", Description: "Read SHA1 sums from the FILEs and check them"},
			{Name: "--tag", Description: "Create a BSD-style checksum"},
			{Name: "-t", Description: "Read in text mode (default)"},
			{Name: "-z", Description: "End each output line with NUL, not newline, and disable file name escaping"},
			{Name: "--ignore-missing", Description: "Don't fail or report status for missing files"},
			{Name: "--quiet", Description: "Don't print OK for each successfully verified file"},
			{Name: "--status", Description: "Don't output anything, status code shows success"},
			{Name: "--strict", Description: "Exit non-zero for improperly formatted checksum lines"},
			{Name: "-w", Description: "Warn about improperly formatted checksum lines"},
			{Name: "--help", Description: "Output help message and exit"},
			{Name: "--version", Description: "Output version information and exit"},
		},
	})
}
