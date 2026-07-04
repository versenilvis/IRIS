package pkginstaller

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pkgutil",
		Description: "Query and manipulate for macOS Installer packages and receipts",
		Subcommands: []core.Subcommand{
			{Name: "package-id", Description: "The package ID to list the files of"},
			{Name: "path", Description: "The path to update ACLs on"},
			{Name: "group-id", Description: "The group ID to list the packages of"},
			{Name: "pkg-path", Description: "The path to the flat package to expand"},
			{Name: "dir-path", Description: "The path to the directory to expand the package into"},
		},
		Options: []core.Option{
			{Name: "--packages", Description: "List all installed package IDs on the specified --volume"},
			{Name: "--pkgs-plist", Description: "List all of the files installed under the package-id"},
			{Name: "--export-plist", Description: "The package ID to export the plist of"},
			{Name: "--verify", Description: "Run repair_packages(8) to verify the specified package-id"},
			{Name: "--repair", Description: "Run repair_packages(8) to repair the specified package-id"},
			{Name: "--pkg-info", Description: "Print extended information about the specified package-id"},
			{Name: "--pkg-info-plist", Description: "The package ID to print the info of"},
			{Name: "--forget", Description: "Discard all receipt data about package-id, but don't touch installed files"},
			{Name: "--learn", Description: "Update the ACLs of the given path in the receipt identified by --edit-pkg"},
			{Name: "--pkg-groups", Description: "List all of the package groups this package-id is a member of"},
			{Name: "--groups", Description: "List all of the package groups ont he specified --volume"},
			{Name: "--group-plist", Description: "List all of the packages that are members of this group-id"},
			{Name: "--file-info", Description: "Show the metadata known about path"},
			{Name: "--file-info-plist", Description: "Show the metadata known about path in Mac OS X plist(5) format"},
			{Name: "--expand", Description: "Expand the flat package at pkg-path into a new directory specified by dir-path"},
			{Name: "--flatten", Description: "Flatten the dir-path into a new flat package created at pkg-path"},
			{Name: "--bom", Description: "The path to the flat package to extract the BOM from"},
			{Name: "--payload-files", Description: "List the files archived within the uninstalled flat package(s) at path"},
			{Name: "--check-signature", Description: "Check the validity and trust of the signature on the package at pkg-path"},
			{Name: "-h", Description: "A brief summary of commands and usage"},
			{Name: "-f", Description: "Skip confirmation before a potentially destructive or ambiguous action"},
			{Name: "-v", Description: "Output in a human-readable format"},
			{Name: "--volume", Description: "Perform all operations on specified volume"},
			{Name: "--edit-pkg", Description: "Specifies an existing receipt to be modified in-place by --learn"},
			{Name: "--only-files", Description: "List only files (not directories) in --files listing"},
			{Name: "--only-dirs", Description: "List only directories (not files) in --files listing"},
			{Name: "--regexp", Description: "Use regex to match package-id arguments, if an exact match isn't found"},
			{Name: "--pkgs", Description: "Regular expression"},
		},
	})
}
