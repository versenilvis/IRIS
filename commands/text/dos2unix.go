package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dos2unix",
		Description: "DOS to Unix file format converter",
		Options: []spec.Option{
			{Name: "-h", Description: "Show help for dos2unix"},
			{Name: "--allow-chown", Description: "Allow file ownership change in old file mode"},
			{Name: "-ascii", Description: "Convert only line breaks. This is the default conversion mode"},
			{Name: "-iso", Description: "Conversion between DOS and ISO-8859-1 character set"},
			{Name: "-1252", Description: "Use Windows code page 1252 (Western European)"},
			{Name: "-437", Description: "Use DOS code page 850 (Western European)"},
			{Name: "-860", Description: "Use DOS code page 860 (Portuguese)"},
			{Name: "-863", Description: "Use DOS code page 863 (French Canadian)"},
			{Name: "-865", Description: "Use DOS code page 865 (Nordic)"},
			{Name: "-7", Description: "Convert 8 bit characters to 7 bit space"},
			{Name: "-b", Description: "Set conversion mode"},
			{Name: "-D", Description: "Set encoding of displayed text"},
			{Name: "-f", Description: "Force conversion of binary files"},
			{Name: "-gb", Description: "Display file information. No conversion is done"},
			{Name: "--info", Description: "Display file information. No conversion is done"},
			{Name: "-k", Description: "Display program's license"},
			{Name: "-l", Description: "Add additional newline"},
			{Name: "-m", Description: "Don't allow file ownership change in old file mode (default)"},
			{Name: "-o", Description: "Old file mode. Convert file FILE and overwrite output to it"},
			{Name: "-q", Description: "Remove Byte Order Mark (BOM). Do not write a BOM in the output file"},
			{Name: "-s", Description: "Skip binary files (default)"},
			{Name: "-u", Description: "Keep UTF-16 encoding"},
			{Name: "-ul", Description: "Assume that the input format is UTF-16LE"},
			{Name: "-ub", Description: "Assume that the input format is UTF-16BE"},
			{Name: "-v", Description: "Verbose operation"},
			{Name: "-F", Description: "Follow symbolic links and convert the targets"},
			{Name: "-R", Description: "Replace symbolic links with converted files"},
			{Name: "-S", Description: "Keep symbolic links and targets unchanged (default)"},
			{Name: "-V", Description: "Display version number"},
		},
	})
}
