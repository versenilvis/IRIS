package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "mongosh",
		Description: "Default Connection String; Equivalent to running mongosh without any commands",
		Options: []core.Option{
			{Name: "-v", Description: "View your current mongosh version"},
			{Name: "--shell", Description: "Returns information on the options and use of the MongoDB Shell"},
			{Name: "--authenticationDatabase", Description: "This option is available only in MongoDB Enterprise"},
			{Name: "--gssapiServiceName", Description: "Enables retryable writes as the default for sessions in the MongoDB Shell"},
			{Name: "--authenticationMechanism", Description: "MongoDB TLS SSL certificate authentication"},
			{Name: "--nodb", Description: "Prevents the shell from connecting to any database instances"},
			{Name: "--norc", Description: "Prevents the shell from sourcing and evaluating ~/.mongoshrc.js on startup"},
			{Name: "--quiet", Description: "Default-port"},
			{Name: "--tls", Description: "Enables connection to a mongod or mongos that has TLS SSL support enabled"},
			{Name: "--tlsAllowInvalidHostnames", Description: "Disables the specified TLS protocols"},
			{Name: "--tlsCAFile", Description: "Enables connection to a mongod or mongos that has TLS SSL support enabled"},
		},
	})
}
