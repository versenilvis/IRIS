package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "dog",
		Description: "Human-readable host names, nameservers, types, or classes",
		Subcommands: []spec.Subcommand{
			{Name: "A", Description: "Query Domain A Record"},
			{Name: "MX", Description: "Query Domain MX Record"},
			{Name: "CNAME", Description: "Query Domain CNAME Record"},
			{Name: "TXT", Description: "Query Domain TXT Record"},
			{Name: "NS", Description: "Query NS Record"},
			{Name: "SOA", Description: "Query SOA Record"},
			{Name: "TTL", Description: "Query TTL Record"},
			{Name: "ANY +noall +answer", Description: "Query ALL DNS Records"},
		},
		Options: []spec.Option{
			{Name: "-q", Description: "Host name or IP address to query"},
			{Name: "-t", Description: "Type of the DNS record being queried (A, MX, NS...)"},
			{Name: "-n", Description: "Address of the nameserver to send packets to"},
			{Name: "-class", Description: "Network class of the DNS record being queried (IN, CH, HS)"},
			{Name: "--edns", Description: "Whether to OPT in to EDNS (disable, hide, show)"},
			{Name: "--txid", Description: "Set the transaction ID to a specific value"},
			{Name: "-Z", Description: "Set uncommon protocol-level tweaks"},
			{Name: "-U", Description: "Use the DNS protocol over UDP"},
			{Name: "-T", Description: "Use the DNS protocol over TCP"},
			{Name: "-S", Description: "Use the DNS-over-TLS protocol"},
			{Name: "-H", Description: "Use the DNS-over-HTTPS protocol"},
			{Name: "-1", Description: "Short mode: display nothing but the first result"},
			{Name: "-J", Description: "Display the output as JSON"},
			{Name: "--color", Description: "When to colourise the output (always, automatic, never)"},
			{Name: "--seconds", Description: "Do not format durations, display them as seconds"},
			{Name: "--time", Description: "Print how long the response took to arrive"},
		},
	})
}
