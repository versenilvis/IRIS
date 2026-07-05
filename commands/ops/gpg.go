package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "gpg",
		Description: "Encryption and signing tool",
		Options: []spec.Option{
			{Name: "--homedir", Description: "Set the name of the home directory"},
			{Name: "--options", Description: "Read options from file"},
			{Name: "-a", Description: "Create ASCII armored output"},
			{Name: "-o", Description: "Write output to file"},
			{Name: "-u", Description: "Use name as the user ID to sign"},
			{Name: "--default-key", Description: "Use name as default user ID for signatures"},
			{Name: "-r", Description: "Encrypt for user id name"},
			{Name: "--default-recipient", Description: "Use name as default recipient"},
			{Name: "--default-recipient-self", Description: "Use the default key as default recipient"},
			{Name: "--no-default-recipient", Description: "Reset --default-recipient and --default-recipient-self"},
			{Name: "--encrypt-to", Description: "Same as --recipient but this one is intended for in the options file"},
			{Name: "--no-encrypt-to", Description: "Disable the use of all --encrypt-to keys"},
			{Name: "-v", Description: "Give more information during processing"},
			{Name: "-q", Description: "Try to be as quiet as possible"},
			{Name: "-Z", Description: "Set compression level to n"},
			{Name: "-t", Description: "Use canonical text mode"},
			{Name: "-n", Description: "Don't make any changes"},
			{Name: "-i", Description: "Prompt before overwriting any files"},
			{Name: "--batch", Description: "Use batch mode"},
			{Name: "--no-tty", Description: "Make sure that the TTY is never used for any output"},
			{Name: "--no-batch", Description: "Disable batch mode"},
			{Name: "--yes", Description: "Skip key validation"},
			{Name: "--keyserver", Description: "Use name to lookup keys which are not yet in your keyring"},
			{Name: "--no-auto-key-retrieve", Description: "Disables the automatic retrieving of keys"},
			{Name: "--honor-http-proxy", Description: "Try to access the keyserver over the proxy"},
			{Name: "--keyring", Description: "Add file to the list of keyrings"},
			{Name: "--secret-keyring", Description: "Same as --keyring but for the secret keyrings"},
			{Name: "--charset", Description: "Set the name of the native character set"},
			{Name: "--utf8-strings", Description: "Assume that the arguments are already given as UTF8"},
			{Name: "--no-utf8-strings", Description: "Load an extension module"},
			{Name: "--debug", Description: "Set debugging flags"},
			{Name: "--debug-all", Description: "Set all useful debugging flags"},
			{Name: "--status-fd", Description: "Write special status strings to the file descriptor n"},
			{Name: "--logger-fd", Description: "Write log output to file descriptor n and not to stderr"},
			{Name: "--no-comment", Description: "Do not write comment packets"},
			{Name: "--comment", Description: "Use string as comment string in clear text signatures"},
			{Name: "--default-comment", Description: "Force to write the standard comment string"},
			{Name: "--no-version", Description: "Omit the version string in clear text signatures"},
			{Name: "--emit-version", Description: "Force to write the version string"},
			{Name: "-N", Description: "Put the name value pair into the signature as notation data"},
		},
	})
}
