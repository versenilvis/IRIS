package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "wscat",
		Description: "Communicate over websocket",
		Options: []spec.Option{
			{Name: "-c", Description: "Connect to a WebSocket server"},
			{Name: "-V", Description: "Output the version number"},
			{Name: "--auth", Description: "Add basic HTTP authentication header (--connect only)"},
			{Name: "--ca", Description: "Specify a Certificate Authority (--connect only)"},
			{Name: "--cert", Description: "Specify a Client SSL Certificate (--connect only)"},
			{Name: "--host", Description: "Optional host"},
			{Name: "--key", Description: "Specify a Client SSL Certificate's key (--connect only)"},
			{Name: "--max-redirects", Description: "Maximum number of redirects allowed (--connect only) (default: 10)"},
			{Name: "--no-color", Description: "Run without color"},
			{Name: "--passphrase", Description: "Connect via a proxy. Proxy must support CONNECT method"},
			{Name: "--slash", Description: "Set an HTTP header. Repeat to set multiple (--connect only) (default: [])"},
			{Name: "-L", Description: "Follow redirects (--connect only)"},
			{Name: "-l", Description: "Listen on port"},
			{Name: "-n", Description: "Do not check for unauthorized certificates"},
			{Name: "-o", Description: "Optional origin"},
			{Name: "-p", Description: "Optional protocol version"},
			{Name: "-P", Description: "Print a notification when a ping or pong is received"},
			{Name: "-s", Description: "Optional subprotocol (default: [])"},
			{Name: "-w", Description: "Wait given seconds after executing command"},
			{Name: "-x", Description: "Execute command after connecting"},
			{Name: "-h", Description: "Display help for command"},
		},
	})
}
