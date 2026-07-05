package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "newman",
		Description: "Newman is a command-line collection runner for Postman",
		Subcommands: []spec.Subcommand{
			{Name: "run", Description: "Initiate a Postman Collection run from a given URL or path"},
			{Name: "help", Description: "Display help for command"},
		},
		Options: []spec.Option{
			{Name: "-e", Description: "Specify a URL or path to a Postman Environment"},
			{Name: "-g", Description: "Specify a URL or path to a file containing Postman Globals"},
			{Name: "-r", Description: "Specify the reporters to use for this run"},
			{Name: "-n", Description: "Define the number of iterations to run"},
			{Name: "-d", Description: "Specify a data file to use for iterations (either JSON or CSV)"},
			{Name: "--folder", Description: "Exports the final environment to a file after completing the run"},
			{Name: "--export-globals", Description: "Exports the final globals to a file after completing the run"},
			{Name: "--export-collection", Description: "Exports the executed collection to a file after completing the run"},
			{Name: "--postman-api-key", Description: "API Key used to load the resources from the Postman API"},
			{Name: "--bail", Description: "Prevents Newman from automatically following 3XX redirect responses"},
			{Name: "-x", Description: "Specify whether or not to override the default exit code for the current run"},
			{Name: "--silent", Description: "Prevents Newman from showing output to CLI"},
			{Name: "--disable-unicode", Description: "Forces Unicode compliant symbols to be replaced by their plain text equivalents"},
			{Name: "--color", Description: "Enable/Disable colored output (auto|on|off)"},
			{Name: "--delay-request", Description: "Specify the extent of delay between requests (milliseconds)"},
			{Name: "--timeout", Description: "Specify a timeout for collection run (milliseconds)"},
			{Name: "--timeout-request", Description: "Specify a timeout for requests (milliseconds)"},
			{Name: "--timeout-script", Description: "Specify a timeout for scripts (milliseconds)"},
			{Name: "--working-dir", Description: "Specify the path to the working directory"},
			{Name: "--no-insecure-file-read", Description: "Prevents reading the files situated outside of the working directory"},
			{Name: "-k", Description: "Disables SSL validations"},
			{Name: "--ssl-client-cert-list", Description: "Specify the path to a client certificates configurations (JSON)"},
			{Name: "--ssl-client-cert", Description: "Specify the path to a client certificate (PEM)"},
			{Name: "--ssl-client-key", Description: "Specify the path to a client certificate private key"},
			{Name: "--ssl-client-passphrase", Description: "Specify the client certificate passphrase (for protected key)"},
			{Name: "--ssl-extra-ca-certs", Description: "Specify additionally trusted CA certificates (PEM)"},
			{Name: "--cookie-jar", Description: "Specify the path to a custom cookie jar (serialized tough-cookie JSON)"},
			{Name: "--export-cookie-jar", Description: "Exports the cookie jar to a file after completing the run"},
			{Name: "--verbose", Description: "Show detailed information of collection run and each request sent"},
			{Name: "-h", Description: "Display help for command"},
			{Name: "-v", Description: "Output the version number"},
		},
	})
}
