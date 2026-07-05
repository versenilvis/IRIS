package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pass",
		Description: "Pass - stores, retrieves, generates, and synchronizes passwords securely",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Initialize new password storage and use gpg-id for encryption"},
			{Name: "gpg-id", Description: "The gpg-id you want to use to encrypt your password store"},
			{Name: "insert", Description: "Insert a new password into the password store called pass-name"},
			{Name: "pass-name", Description: "The password name"},
			{Name: "git", Description: "Password store git functions"},
			{Name: "version", Description: "Show version information"},
			{Name: "help", Description: "Show usage message"},
			{Name: "cp", Description: "Copies the password or directory named old-path to new-path"},
			{Name: "old-path", Description: "The old password name or directory"},
			{Name: "new-path", Description: "The new password name or directory"},
			{Name: "mv", Description: "Renames the password or directory named old-path to new-path"},
			{Name: "rm", Description: "Remove the password named pass-name from the password store"},
			{Name: "generate", Description: "Generate a new password of length pass-length and insert into pass-name"},
			{Name: "pass-length", Description: "The length of the password"},
			{Name: "ls", Description: "List names of passwords inside the tree at subfolder by using the tree"},
			{Name: "password sub-directory", Description: "The password sub directory you want to list"},
			{Name: "find", Description: "List names of passwords inside the tree that match pass-names"},
			{Name: "show", Description: "Decrypt and print a password"},
		},
		Options: []core.Option{
			{Name: "--clip", Description: "Copy the password to the clipboard"},
			{Name: "--qrcode", Description: "Display a QRcode of the password"},
			{Name: "--help", Description: "Show help for pass"},
			{Name: "--path", Description: "Insert a new password into the password store called pass-name"},
			{Name: "--echo", Description: "Don't prompt before overwriting an existing password"},
			{Name: "--force", Description: "Do not interactively prompt before moving"},
			{Name: "--recursive", Description: "Delete pass-name recursively if it is a directory"},
			{Name: "--no-symbols", Description: "Do not use any non-alphanumeric characters in the generated password"},
			{Name: "--in-place", Description: "Overwrite the existing password"},
		},
	})
}
