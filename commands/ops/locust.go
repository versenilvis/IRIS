package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "locust",
		Description: "Show program",
		Options: []spec.Option{
			{Name: "-v", Description: "Show program's version number and exit"},
			{Name: "-h", Description: "Show this help message and exit"},
			{Name: "-f", Description: "Show program's version number and exit"},
			{Name: "--config", Description: "Config file path"},
			{Name: "-H", Description: "Host to load test in the following format: http://10.21.32.33"},
			{Name: "-u", Description: "Stop after the specified amount of time, e.g. (300s, 20m, 3h, 1h30m, etc.)"},
			{Name: "-l", Description: "Show list of possible User classes and exit"},
			{Name: "--web-host", Description: "Host to bind the web interface to. Defaults to '*' (all interfaces)"},
			{Name: "-P", Description: "Port on which to run web host"},
			{Name: "--headless", Description: "Disable the web interface, and start the test immediately"},
			{Name: "--autostart", Description: "Starts the test immediately (without disabling the web UI)"},
			{Name: "--autoquit", Description: "Quits Locust entirely, X seconds after the run is finished"},
			{Name: "--web-auth", Description: "Turn on Basic Auth for the web interface. e.g. username:password"},
			{Name: "--tls-cert", Description: "Optional path to TLS certificate to use to serve over HTTPS"},
			{Name: "--tls-key", Description: "Optional path to TLS private key to use to serve over HTTPS"},
			{Name: "--class-picker", Description: "Set locust to run in distributed mode with this process as master"},
			{Name: "--master-bind-host", Description: "Interfaces (hostname, ip) that locust master should bind to"},
			{Name: "--master-bind-port", Description: "Port that locust master should bind to"},
			{Name: "--expect-workers", Description: "Set locust to run in distributed mode with this process as worker"},
			{Name: "--master-host", Description: "Show list of possible User classes and exit"},
			{Name: "--csv", Description: "Store current request stats to files in CSV format"},
			{Name: "--csv-full-history", Description: "Print stats in the console"},
			{Name: "--only-summary", Description: "Only print the summary stats"},
			{Name: "--reset-stats", Description: "Store HTML report to file path specified"},
			{Name: "--skip-log-setup", Description: "Choose between DEBUG/INFO/WARNING/ERROR/CRITICAL. Default is INFO"},
			{Name: "--logfile", Description: "Path to log file. If not set, log will go to stderr"},
		},
	})
}
