package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "hyperfine",
		Description: "A command-line benchmarking tool",
		Options: []spec.Option{
			{Name: "--warmup", Description: "Perform warmupruns (number) before the actual benchmarking starts"},
			{Name: "--min-runs", Description: "Perform at least NUM runs for each command"},
			{Name: "--max-runs", Description: "Perform at most NUM runs for each command. Default: no limit"},
			{Name: "--runs", Description: "Perform exactly NUM runs for each command"},
			{Name: "--setup", Description: "Execute cmd once before each set of timing runs"},
			{Name: "--prepare", Description: "Perform benchmark runs for each value in the range min..max"},
			{Name: "--parameter-step-size", Description: "Perform benchmark runs for each value in the comma-separated list of values"},
			{Name: "--style", Description: "Set output style type"},
			{Name: "--shell", Description: "Set the shell to use for executing benchmarked commands"},
			{Name: "--ignore-failure", Description: "Ignore non-zero exit codes of the benchmarked commands"},
			{Name: "--time-unit", Description: "Set the time unit to use for the benchmark results"},
			{Name: "--export-asciidoc", Description: "Export the timing summary statistics as an AsciiDoc table to the given file"},
			{Name: "--export-csv", Description: "Export the timing summary statistics as CSV to the given file"},
			{Name: "--export-json", Description: "Export the timing summary statistics as a Markdown table to the given file"},
			{Name: "--show-output", Description: "Print the stdout and stderr of the benchmark instead of suppressing it"},
			{Name: "--command-name", Description: "Identify a command with the given name"},
			{Name: "--help", Description: "Prints help message"},
			{Name: "--version", Description: "Shows version information"},
		},
	})
}
