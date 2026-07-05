package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "wrk",
		Description: "Wrk - a HTTP benchmarking tool",
		Options: []spec.Option{
			{Name: "-c", Description: "Connections to keep open"},
			{Name: "-d", Description: "Duration of test"},
			{Name: "-t", Description: "Number of threads"},
			{Name: "-s", Description: "Load Lua script file"},
			{Name: "-H", Description: "Add header to request"},
			{Name: "--latency", Description: "Print latency statistics"},
			{Name: "--timeout", Description: "Socket/request timeout"},
			{Name: "-v", Description: "Print version details"},
			{Name: "-h", Description: "Output usage information"},
		},
	})
}
