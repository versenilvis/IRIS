package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "croc",
		Description: "Send file(s), or folder",
		Subcommands: []spec.Subcommand{
			{Name: "send", Description: "Send file(s), or folder"},
			{Name: "relay", Description: "Start your own relay"},
			{Name: "help", Description: "Shows a list of commands or help for one command"},
		},
		Options: []spec.Option{
			{Name: "--code", Description: "Codephrase used to connect to relay"},
			{Name: "--hash", Description: "Hash algorithm"},
			{Name: "--text", Description: "Send some text"},
			{Name: "--no-local", Description: "Disable local relay when sending"},
			{Name: "--no-multi", Description: "Disable multiplexing"},
			{Name: "--ports", Description: "Ports of the local relay"},
			{Name: "--host", Description: "Host of the relay"},
			{Name: "--help", Description: "Show help for croc"},
			{Name: "--internal-dns", Description: "Use a built-in DNS stub resolver rather than the host operating system"},
			{Name: "--remember", Description: "Save these settings to reuse next time"},
			{Name: "--debug", Description: "Toggle debug mode"},
			{Name: "--yes", Description: "Automatically agree to all prompts"},
			{Name: "--stdout", Description: "Redirect file to stdout"},
			{Name: "--no-compress", Description: "Disable compression"},
			{Name: "--ask", Description: "Make sure sender and recipient are prompted"},
			{Name: "--local", Description: "Force to use only local connections"},
			{Name: "--ignore-stdin", Description: "Ignore piped stdin"},
			{Name: "--overwrite", Description: "Do not prompt to overwrite"},
			{Name: "--curve", Description: "Choose an encryption curve"},
			{Name: "--ip", Description: "Set sender ip if known"},
			{Name: "--relay", Description: "Address of the relay"},
			{Name: "--relay6", Description: "Ipv6 address of the relay"},
			{Name: "--out", Description: "Specify an output folder to receive the file"},
			{Name: "--pass", Description: "Add a socks5 proxy"},
			{Name: "--throttleUpload", Description: "Throttle the upload speed e.g. 500k"},
			{Name: "--version", Description: "Print the version"},
		},
	})
}
