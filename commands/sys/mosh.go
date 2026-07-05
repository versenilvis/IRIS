package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "mosh",
		Description: "Address of remote machine to log into",
		Options: []spec.Option{
			{Name: "--help", Description: "Show help for mosh"},
			{Name: "--client", Description: "Local echo options"},
			{Name: "-4", Description: "Use IPv4 only"},
			{Name: "-6", Description: "Use IPv6 only"},
			{Name: "--family", Description: "Network Type"},
			{Name: "--port", Description: "Server-side UDP port or range, (No effect on server-side SSH port)"},
			{Name: "--bind-server", Description: "Do not allocate a pseudo tty on ssh connection"},
			{Name: "--no-init", Description: "Do not send terminal initialization string"},
			{Name: "--local", Description: "Run mosh-server locally without using ssh"},
			{Name: "--experimental-remote-ip", Description: "Select the method for discovering the remote IP address to use for mosh"},
			{Name: "--version", Description: "Version and copyright information"},
		},
	})
}
