package text

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "pandoc",
		Description: "A universal document converter",
		Options: []spec.Option{
			{Name: "-f", Description: "Specify input format"},
			{Name: "-t", Description: "Specify output format"},
			{Name: "-o", Description: "Write output to FILE instead of stdout"},
			{Name: "--data-dir", Description: "Specify the user data directory to search for pandoc data files"},
			{Name: "-d", Description: "Specify a set of default option settings"},
			{Name: "--bash-completion", Description: "Generate a bash completion script"},
			{Name: "--verbose", Description: "Give verbose debugging output"},
			{Name: "--quiet", Description: "Suppress warning messages"},
			{Name: "--fail-if-warnings", Description: "Exit with error status if there are any warnings"},
			{Name: "--log", Description: "Write log messages in machine-readable JSON format to FILE"},
			{Name: "--list-input-formats", Description: "List supported input formats, one per line"},
			{Name: "--list-output-formats", Description: "List supported output formats, one per line"},
			{Name: "--list-extensions", Description: "List supported languages for syntax highlighting, one per line"},
			{Name: "--list-highlight-styles", Description: "List supported styles for syntax highlighting, one per line"},
			{Name: "-v", Description: "Print version"},
			{Name: "-h", Description: "Show usage message"},
			{Name: "--shift-heading-level-by", Description: "Shift heading levels by a positive or negative integer"},
			{Name: "--indented-code-classes", Description: "Set the metadata field KEY to the value VAL"},
			{Name: "--metadata-file", Description: "Read metadata from the supplied YAML (or JSON) file"},
			{Name: "-p", Description: "Preserve tabs instead of converting them to spaces"},
			{Name: "--tab-stop", Description: "Specify the number of spaces per tab"},
			{Name: "--track-changes", Description: "Processes all the insertions and deletions and ignores comments"},
			{Name: "--extract-media", Description: "Use the specified file as a custom template for the generated document"},
			{Name: "-V", Description: "Print a system default data file. Files in the user data directory are ignored"},
			{Name: "--eol", Description: "Windows"},
			{Name: "--dpi", Description: "Attempts to wrap lines to the column width specified by --columns (default 72)"},
			{Name: "--columns", Description: "Specify the number of section levels to include in the table of contents"},
			{Name: "--strip-comments", Description: "Specifies the coloring style to be used in highlighted source code"},
			{Name: "--print-highlight-style", Description: "Include contents of FILE, verbatim, at the end of the header"},
			{Name: "-B", Description: "Include contents of FILE, verbatim, at the beginning of the document body"},
			{Name: "-A", Description: "Include contents of FILE, verbatim, at the end of the document body"},
			{Name: "--resource-path", Description: "Disable the certificate verification to allow access to unsecure HTTP resources"},
			{Name: "--self-contained", Description: "Deprecated synonym for --markdown-headings=atx"},
			{Name: "--top-level-division", Description: "Offset for section headings in HTML output (ignored in other output formats)"},
			{Name: "--listings", Description: "Make list items in slide shows display incrementally"},
			{Name: "--slide-level", Description: "Number"},
			{Name: "--section-divs", Description: "Leaves mailto: links as they are"},
			{Name: "--id-prefix", Description: "Use the specified file as a style reference in producing a docx or ODT file"},
			{Name: "--epub-cover-image", Description: "Use the specified engine when producing PDF output"},
			{Name: "--pdf-engine-opt", Description: "Use the given string as a command-line argument to the pdf-engine"},
		},
	})
}
