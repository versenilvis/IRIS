package sys

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "traceroute",
		Description: "Print the route packets take to network host",
		Options: []core.Option{
			{Name: "--help", Description: "Show help for traceroute"},
			{Name: "-a", Description: "Turn on AS# lookups for each hop encountered"},
			{Name: "-A", Description: "Turn on AS# lookups and use the given server instead of the default"},
			{Name: "-d", Description: "Enable socket level debugging"},
			{Name: "-D", Description: "Set the initial time-to-live used in the first outgoing probe packet"},
			{Name: "-F", Description: "Set the `don't fragment` bit"},
			{Name: "-g", Description: "Specify a loose source route gateway (8 maximum)"},
			{Name: "-i", Description: "Use ICMP ECHO instead of UDP datagrams. (A synonym for `-P icmp`)"},
			{Name: "-M", Description: "Set the number of probes per ``ttl'' to nqueries (default is three probes)"},
			{Name: "-r", Description: "Print a summary of how many probes were not answered for each hop"},
			{Name: "-t", Description: "Set the time (in seconds) to wait for a response to a probe (default 5 sec.)"},
		},
	})
}
