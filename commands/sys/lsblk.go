package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "lsblk",
		Description: "List block devices",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for lsblk"},
			{Name: "--version", Description: "Show version for lsblk"},
			{Name: "--all", Description: "Also list empty devices and RAM disk devices"},
			{Name: "--bytes", Description: "Print the SIZE column in bytes"},
			{Name: "--discard", Description: "Do not print holder devices or slaves"},
			{Name: "--dedup", Description: "Use column as a de-duplication key to de-duplicate output tree"},
			{Name: "--exclude", Description: "Output info about filesystems"},
			{Name: "--include", Description: "Include devices specified by the comma-separated list of major device numbers"},
			{Name: "--ascii", Description: "Use ASCII characters for tree formatting"},
			{Name: "--json", Description: "Use JSON output format"},
			{Name: "--list", Description: "Produce output in the form of a list"},
			{Name: "--merge", Description: "Output info about device owner, group and mode"},
			{Name: "--noheadings", Description: "Do not print a header line"},
			{Name: "--output", Description: "Specify which output columns to print"},
			{Name: "--output-all", Description: "Output all available columns"},
			{Name: "--pairs", Description: "Produce output in the form of key-value pairs"},
			{Name: "--raw", Description: "Produce output in raw format"},
			{Name: "--scsi", Description: "Output info about SCSI devices only"},
			{Name: "--inverse", Description: "Print dependencies in inverse order"},
			{Name: "--tree", Description: "Force tree-like output format"},
			{Name: "--topology", Description: "Output info about block-device topology"},
			{Name: "--width", Description: "Specifies output width as a number of characters"},
			{Name: "--sort", Description: "Sort output lines by column"},
			{Name: "--zoned", Description: "Print the zone model for each device"},
			{Name: "--sysroot", Description: "Device to list"},
		},
	})
}
