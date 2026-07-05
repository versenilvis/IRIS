package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pkg-config",
		Description: "Return metainformation about installed libraries",
		Options: []spec.Option{
			{Name: "--mod-version", Description: "Display the version information of the libraries specified on the command line"},
			{Name: "--version", Description: "Display the version of pkg-config and terminates"},
			{Name: "--atleast-pkgconfig-version", Description: "Require at least the given version of pkg-config"},
			{Name: "--help", Description: "Displays a help message and terminates"},
			{Name: "--print-errors", Description: "Print short error messages"},
			{Name: "--silence-errors", Description: "If printing errors, print them to stdout rather than the default stderr"},
			{Name: "--debug", Description: "Print debugging information"},
			{Name: "--cflags", Description: "Print link flags required to compile the packages on the command line"},
			{Name: "--libs-only-L", Description: "Return the value of a variable defined in a package's .pc file"},
			{Name: "--define-variable", Description: "Set a global value for a variable, overriding the value in any .pc files"},
			{Name: "--print-variables", Description: "Return a list of all variables defined in the package"},
			{Name: "--uninstalled", Description: "Test whether the packages on the command line exist"},
			{Name: "--atleast-version", Description: "Check the syntax of a package's .pc file for validity"},
			{Name: "--msvc-syntax", Description: "Use the installed location of the .pc file to determine the prefix"},
			{Name: "--dont-define-prefix", Description: "Use the specified prefix variable value defined in the .pc file as the prefix"},
			{Name: "--prefix-variable", Description: "Output libraries suitable for static linking"},
			{Name: "--list-all", Description: "List all modules found in the pkg-config path"},
			{Name: "--print-provides", Description: "List all modules the given packages provides"},
			{Name: "--print-requires", Description: "List all modules the given packages requires"},
			{Name: "--print-requires-private", Description: "List all modules the given packages requires for static linking"},
		},
	})
}
