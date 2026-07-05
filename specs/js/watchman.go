package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "watchman",
		Description: "A file watching service",
		Subcommands: []core.Subcommand{
			{Name: "clock", Description: "Returns the current clock value for a watched root"},
			{Name: "path", Description: "The path to directory"},
			{Name: "find", Description: "Finds all files that match the optional list of patterns under the specified dir"},
			{Name: "get-config", Description: "Returns the .watchmanconfig for the root"},
			{Name: "get-sockname", Description: "Get socket path"},
			{Name: "list-capabilities", Description: "Returns the full list of supported capabilities offered by the watchman server"},
			{Name: "log", Description: "Generates a log line in the watchman log"},
			{Name: "level", Description: "The log level"},
			{Name: "log-level", Description: "Changes the log level of your connection to the watchman service"},
			{Name: "query", Description: "Executes a query against the specified root"},
			{Name: "shutdown-server", Description: "This causes your watchman service to exit with a normal status code"},
			{Name: "state-enter", Description: "This causes a watch to be marked as being in a particular named state"},
			{Name: "state-leave", Description: "This causes a watch to no longer be marked as being in a particular named state"},
			{Name: "trigger", Description: "This will create or replace a trigger"},
			{Name: "trigger-del", Description: "Deletes a named trigger from the list of registered triggers"},
			{Name: "trigger-list", Description: "Returns the set of registered triggers associated with a root directory"},
			{Name: "unsubscribe", Description: "Cancels a named subscription against the specified root"},
			{Name: "version", Description: "The version and build information for the currently running watchman service"},
			{Name: "watch-del", Description: "Removes a watch and any associated triggers"},
			{Name: "watch-del-all", Description: "Removes all watches and associated triggers"},
			{Name: "watch-list", Description: "Returns a list of watched dirs"},
			{Name: "watch-project", Description: "Requests that the project containing the requested dir is watched for changes"},
		},
		Options: []core.Option{
			{Name: "--help", Description: "Show help for watchman"},
			{Name: "--inetd", Description: "Spawning from an inetd style supervisor"},
			{Name: "-S", Description: "Don't use the site or system spawner"},
			{Name: "-v", Description: "Show version number for watchman"},
			{Name: "--named-pipe-path", Description: "Specify alternate named pipe path"},
			{Name: "-u", Description: "Specify alternate unix domain socket path"},
			{Name: "--unix-listener-path", Description: "Specify alternate unix domain socket path"},
			{Name: "-o", Description: "Specify the path to logfile"},
			{Name: "--logfile", Description: "Specify the path to logfile"},
			{Name: "--log-level", Description: "Set the log level"},
			{Name: "--pidfile", Description: "Specify path to pidfile"},
			{Name: "-p", Description: "Persist and wait for further responses"},
			{Name: "-n", Description: "Don't save state between invocations"},
			{Name: "--statefile", Description: "Specify path to file to hold watch and trigger state"},
			{Name: "-j", Description: "Instead of parsing CLI arguments, take a single json object from stdin"},
			{Name: "--output-encoding", Description: "CLI output encoding"},
			{Name: "--server-encoding", Description: "CLI<->server encoding"},
			{Name: "-f", Description: "Run the service in the foreground"},
			{Name: "--no-pretty", Description: "Don't pretty print JSON"},
			{Name: "--no-spawn", Description: "Don't try to start the service if it is not available"},
		},
	})
}
