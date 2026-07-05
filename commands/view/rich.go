package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rich",
		Description: "Defined by terminal, appearance may differ",
		Options: []spec.Option{
			{Name: "-p", Description: "Print console markup. See https://rich.readthedocs.io/en/latest/markup.html"},
			{Name: "-u", Description: "Display a horizontal rule"},
			{Name: "-j", Description: "Display as JSON"},
			{Name: "-m", Description: "Display as markdown"},
			{Name: "--rst", Description: "Display restructured text"},
			{Name: "--csv", Description: "Display CSV as a table"},
			{Name: "--ipynb", Description: "Display Jupyter notebook"},
			{Name: "--syntax", Description: "Syntax highlighting"},
			{Name: "--inspect", Description: "Inspect a python object"},
			{Name: "-h", Description: "Display first LINES of the file"},
			{Name: "-t", Description: "Display last LINES of the file"},
			{Name: "-l", Description: "Align to left"},
			{Name: "-r", Description: "Align to right"},
			{Name: "-c", Description: "Align to center"},
			{Name: "-L", Description: "Justify text to left"},
			{Name: "-R", Description: "Justify text to right"},
			{Name: "-C", Description: "Justify text to center"},
			{Name: "-F", Description: "Justify text to both left and right edges"},
			{Name: "--soft", Description: "Enable soft wrapping of text"},
			{Name: "-e", Description: "Expand to full width"},
			{Name: "-w", Description: "Fit output to SIZE characters"},
			{Name: "-W", Description: "Set maximum width to SIZE characters"},
			{Name: "-s", Description: "Set text style to STYLE"},
			{Name: "--rule-style", Description: "Set rule style to STYLE"},
			{Name: "--rule-char", Description: "Use CHARACTER to generate a line with --rule"},
			{Name: "-d", Description: "Padding around output. 1, 2 or 4 comma separated integers, e.g. 2,4"},
			{Name: "-a", Description: "Set panel type to BOX"},
			{Name: "-S", Description: "Set the panel style to STYLE"},
			{Name: "--theme", Description: "Set syntax theme to THEME. See https://pygments.org/styles/"},
			{Name: "-n", Description: "Enable line number in syntax"},
			{Name: "-g", Description: "Enable indentation guides in syntax highlighting"},
			{Name: "-y", Description: "Render hyperlinks in markdown"},
			{Name: "--no-wrap", Description: "Don't word wrap syntax highlighted files"},
			{Name: "--title", Description: "Set panel title to TEXT"},
			{Name: "--caption", Description: "Set panel caption to TEXT"},
			{Name: "--force-terminal", Description: "Force terminal output when not writing to a terminal"},
			{Name: "-o", Description: "Write HTML to PATH"},
			{Name: "--export-svg", Description: "Write SVG to PATH"},
			{Name: "--pager", Description: "Display in an interactive pager"},
			{Name: "-v", Description: "Print version and exit"},
		},
	})
}
