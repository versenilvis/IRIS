package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "nmap",
		Description: "Network exploration tool and security / port scanner",
		Options: []core.Option{
			{Name: "-iR", Description: "Choose random targets"},
			{Name: "--exclude", Description: "Excluide hosts/networks"},
			{Name: "--excludefile", Description: "Exclude list from file"},
			{Name: "-sT", Description: "TCP scan"},
			{Name: "-sA", Description: "TCP ACK scan"},
			{Name: "-sW", Description: "TCP Window scan"},
			{Name: "-sM", Description: "TCP Maimon scan"},
			{Name: "-sU", Description: "UDP scan"},
			{Name: "-sP", Description: "Ping scan"},
			{Name: "-sN", Description: "TCP Null scan"},
			{Name: "-sF", Description: "FIN scan"},
			{Name: "-sX", Description: "Xmas scan"},
			{Name: "-sO", Description: "IP protocol scan"},
			{Name: "-p", Description: "Scan specified ports"},
			{Name: "-v", Description: "Increase verbosity level"},
			{Name: "--osscan-limit", Description: "Limit OS detection to promising targets"},
			{Name: "-6", Description: "Enable IPV6 scanning"},
			{Name: "-V", Description: "Print version number"},
			{Name: "-privileged", Description: "Assume that user is fully privileges"},
			{Name: "-unprivileged", Description: "Assume that user lacks raw socket privileges"},
			{Name: "--help", Description: "Show help for nmap"},
		},
	})
}
