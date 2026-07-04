package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "jmeter",
		Description: "Apache JMeter - 100% Java Load Testing Tool",
		Options: []core.Option{
			{Name: "-v", Description: "Print the JMeter version information and exit"},
			{Name: "-h", Description: "Print usage information and exit"},
			{Name: "-p", Description: "The jmeter property file to use"},
			{Name: "-q", Description: "Additional JMeter property file(s)"},
			{Name: "-t", Description: "The file to log samples to"},
			{Name: "-i", Description: "JMeter logging configuration file"},
			{Name: "-j", Description: "JMeter run log file"},
			{Name: "-n", Description: "Run JMeter in nongui mode"},
			{Name: "-s", Description: "Run the JMeter server"},
			{Name: "-E", Description: "Set a proxy scheme to use for the proxy server"},
			{Name: "-H", Description: "Set a proxy server for JMeter to use"},
			{Name: "-P", Description: "Set proxy server port for JMeter to use"},
			{Name: "-N", Description: "Set nonproxy host list (e.g. *.apache.org|localhost)"},
			{Name: "-u", Description: "Set username for proxy server that JMeter is to use"},
			{Name: "-a", Description: "Set password for proxy server that JMeter is to use"},
			{Name: "-J", Description: "Define additional JMeter properties <argument>=<value>"},
			{Name: "-G", Description: "Define additional system properties <argument>=<value>"},
			{Name: "-S", Description: "Additional system property file(s)"},
			{Name: "-f", Description: "Force delete existing results files and web report folder"},
			{Name: "-L", Description: "[category=]level e.g. jorphan=INFO, jmeter.util=DEBUG or com.example.foo=WARN"},
			{Name: "-r", Description: "Start remote servers (as defined in remote_hosts)"},
			{Name: "-R", Description: "Start these remote servers (overrides remote_hosts)"},
			{Name: "-d", Description: "The JMeter home directory to use"},
			{Name: "-X", Description: "Exit the remote servers at end  of test (non-GUI)"},
			{Name: "-g", Description: "Generate report dashboard only, from a test results file"},
			{Name: "-e", Description: "Generate report dashboard after load test"},
			{Name: "-o", Description: "Output folder for report dashboard"},
		},
	})
}
