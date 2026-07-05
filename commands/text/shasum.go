package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "shasum",
		Description: "Print or Check SHA Checksums",
		Options: []spec.Option{
			{Name: "-a", Description: "Select SHA algorithm"},
			{Name: "-b", Description: "Read in binary mode"},
			{Name: "-c", Description: "Read SHA sums from the FILEs and check them"},
			{Name: "--tag", Description: "Create a BSD-style checksum"},
			{Name: "-t", Description: "Read in text mode (default)"},
			{Name: "-U", Description: "Read in Universal Newlines mode - produces same digest on Windows/Unix/Mac"},
			{Name: "-0", Description: "Read in BITS mode - ASCII '0' as 0-bit, ASCII '1' as 1-bit, others ignored"},
			{Name: "--ignore-missing", Description: "Don't fail or report status for missing files"},
			{Name: "-q", Description: "Don't print OK for each successfully verified file"},
			{Name: "-s", Description: "Don't output anything, status code shows success"},
			{Name: "--strict", Description: "Exit non-zero for improperly formatted checksum lines"},
			{Name: "-w", Description: "Warn about improperly formatted checksum lines"},
			{Name: "-h", Description: "Display help and exit"},
			{Name: "-v", Description: "Output version information and exit"},
		},
	})
}
