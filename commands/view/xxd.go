package view

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xxd",
		Description: "Make a hexdump or do the reverse",
		Options: []spec.Option{
			{Name: "-help", Description: "Show help for xxd"},
			{Name: "-autoskip", Description: "Toggle autoskip: A single '*' replaces nul-lines.  Default off"},
			{Name: "-bits", Description: "Switch to bits (binary digits) dump, rather than hexdump"},
			{Name: "-cols", Description: "Format <cols> octets per line. Default 16"},
			{Name: "-capitalize", Description: "Capitalize variable names in C include file style, when using -i"},
			{Name: "-EBCDIC", Description: "Change the character encoding in the righthand column from ASCII to EBCDIC"},
			{Name: "-e", Description: "Switch to little-endian hexdump"},
			{Name: "-groupsize", Description: "Separate the output of every <bytes> bytes"},
			{Name: "-include", Description: "Output in C include file style"},
			{Name: "-len", Description: "Stop after writing <len> octets"},
			{Name: "-name", Description: "Override the variable name output when -i is used"},
			{Name: "-o", Description: "Add <offset> to the displayed file position"},
			{Name: "-postscript", Description: "Output in postscript continuous hexdump style"},
			{Name: "-revert", Description: "Reverse operation: convert (or patch) hexdump into binary"},
			{Name: "-seek", Description: "Use upper case hex letters. Default is lower case"},
			{Name: "-version", Description: "Show version string"},
		},
	})
}
