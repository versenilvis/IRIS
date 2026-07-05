package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "curl",
		Description: "transfer data via URL",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-X", Description: "request method (GET/POST/PUT/DELETE)"},
			{Name: "-H", Description: "request header"},
			{Name: "-d", Description: "request body data"},
			{Name: "-o", Description: "output to file"},
			{Name: "-O", Description: "save with remote filename"},
			{Name: "-s", Description: "silent mode"},
			{Name: "-S", Description: "show error in silent"},
			{Name: "-L", Description: "follow redirects"},
			{Name: "-k", Description: "insecure (skip TLS verify)"},
			{Name: "-I", Description: "fetch headers only"},
			{Name: "-v", Description: "verbose"},
			{Name: "-u", Description: "user:password"},
			{Name: "-b", Description: "send cookie"},
			{Name: "-c", Description: "save cookie to file"},
			{Name: "-F", Description: "multipart form data"},
			{Name: "--json", Description: "send JSON body"},
			{Name: "--compressed", Description: "request compressed response"},
			{Name: "--max-time", Description: "max time in seconds"},
			{Name: "--retry", Description: "retry count"},
			{Name: "-A", Description: "custom user-agent"},
			{Name: "--proxy", Description: "use proxy"},
			{Name: "-T", Description: "upload file"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "wget",
		Description: "non-interactive downloader",
		Generator:   spec.FileGenerator(),
		Options: []spec.Option{
			{Name: "-O", Description: "output filename"},
			{Name: "-q", Description: "quiet"},
			{Name: "-v", Description: "verbose"},
			{Name: "-r", Description: "recursive"},
			{Name: "-P", Description: "output directory"},
			{Name: "-c", Description: "continue download"},
			{Name: "--no-check-certificate", Description: "skip TLS verify"},
			{Name: "-b", Description: "background"},
			{Name: "--limit-rate", Description: "limit rate"},
			{Name: "-U", Description: "custom user-agent"},
			{Name: "--spider", Description: "check URL only"},
			{Name: "-np", Description: "no parent directories"},
			{Name: "-nH", Description: "no host directories"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "nc",
		Description: "netcat - TCP/UDP tool",
		Options: []spec.Option{
			{Name: "-l", Description: "listen mode"},
			{Name: "-p", Description: "local port"},
			{Name: "-z", Description: "port scan mode"},
			{Name: "-v", Description: "verbose"},
			{Name: "-u", Description: "UDP mode"},
			{Name: "-n", Description: "no DNS"},
			{Name: "-w", Description: "timeout"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "ping",
		Description: "test network connectivity",
		Options: []spec.Option{
			{Name: "-c", Description: "packet count"},
			{Name: "-i", Description: "interval"},
			{Name: "-t", Description: "TTL"},
			{Name: "-s", Description: "packet size"},
			{Name: "-q", Description: "quiet output"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "dig",
		Description: "DNS lookup",
		Options: []spec.Option{
			{Name: "+short", Description: "short output"},
			{Name: "+noall", Description: "no output flags"},
			{Name: "+answer", Description: "answer section only"},
			{Name: "-t", Description: "record type"},
			{Name: "@", Description: "nameserver"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "nslookup",
		Description: "query DNS",
	})

	spec.Register(&spec.Spec{
		Name:        "ifconfig",
		Description: "configure network interface",
		Options: []spec.Option{
			{Name: "-a", Description: "all interfaces"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "ip",
		Description: "show/manage network",
		Subcommands: []spec.Subcommand{
			{Name: "addr", Description: "address management"},
			{Name: "link", Description: "network device management"},
			{Name: "route", Description: "routing table management"},
			{Name: "neigh", Description: "ARP table"},
			{Name: "rule", Description: "routing policy"},
		},
		Options: []spec.Option{
			{Name: "-4", Description: "IPv4 only"},
			{Name: "-6", Description: "IPv6 only"},
			{Name: "-br", Description: "brief output"},
			{Name: "-c", Description: "colorize"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "netstat",
		Description: "network statistics",
		Options: []spec.Option{
			{Name: "-t", Description: "TCP connections"},
			{Name: "-u", Description: "UDP connections"},
			{Name: "-l", Description: "listening only"},
			{Name: "-n", Description: "numeric addresses"},
			{Name: "-p", Description: "show PID"},
			{Name: "-r", Description: "routing table"},
			{Name: "-s", Description: "statistics"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "ss",
		Description: "socket statistics",
		Options: []spec.Option{
			{Name: "-t", Description: "TCP"},
			{Name: "-u", Description: "UDP"},
			{Name: "-l", Description: "listening"},
			{Name: "-n", Description: "numeric"},
			{Name: "-p", Description: "show process"},
			{Name: "-a", Description: "all"},
			{Name: "-4", Description: "IPv4"},
			{Name: "-6", Description: "IPv6"},
		},
	})
}
