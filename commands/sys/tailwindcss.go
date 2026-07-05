package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tailwindcss",
		Description: "Display usage information",
		Subcommands: []spec.Subcommand{
			{Name: "init", Description: "Creates Tailwind config file. Default: tailwind.config.js"},
			{Name: "build", Description: "Build CSS file"},
		},
		Options: []spec.Option{
			{Name: "--help", Description: "Display usage information"},
			{Name: "-i", Description: "Specify input file"},
			{Name: "-o", Description: "Specify output file"},
			{Name: "-c", Description: "Specify config file to use"},
			{Name: "--postcss", Description: "Load custom PostCSS configuration"},
			{Name: "--purge", Description: "Content paths to use for removing unused classes"},
			{Name: "--watch", Description: "Watch for changes and rebuild as needed"},
			{Name: "--minify", Description: "Minify the output"},
			{Name: "--no-autoprefixer", Description: "Disable autoprefixer"},
			{Name: "-p", Description: "Initialize a 'postcss.config.js' file"},
			{Name: "-f", Description: "Initialize a full 'tailwind.config.js' file"},
		},
	})
}
