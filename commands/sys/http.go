package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "http",
		Description: "HTTPie: command-line HTTP client for the API era",
		Options: []spec.Option{
			{Name: "--json", Description: "Disables all sorting while formatting output"},
			{Name: "--sorted", Description: "Re-enables all sorting options while formatting output"},
			{Name: "--format-options", Description: "Controls output formatting"},
			{Name: "--print", Description: "Request headers"},
			{Name: "--headers", Description: "Print only the response headers. Shortcut for --print=h"},
			{Name: "--body", Description: "Print only the response body. Shortcut for --print=b"},
			{Name: "--verbose", Description: "Request headers"},
			{Name: "--stream", Description: "Always stream the response body by line, i.e., behave like `tail -f'"},
			{Name: "--output", Description: "Create or read a session without updating it form the request/response exchange"},
			{Name: "--auth", Description: "Basic HTTP auth"},
			{Name: "--ignore-netrc", Description: "Ignore credentials from .netrc"},
			{Name: "--offline", Description: "Build the request and print it but don’t actually send it"},
			{Name: "--proxy", Description: "Follow 30x Location redirects"},
			{Name: "--max-redirects", Description: "By default, requests have a limit of 30 redirects (works with --follow)"},
			{Name: "--max-headers", Description: "Bypass dot segment (/../ or /./) URL squashing"},
			{Name: "--chunked", Description: "Enable streaming via chunked transfer encoding"},
			{Name: "--verify", Description: "A string in the OpenSSL cipher list format"},
			{Name: "--cert", Description: "Do not attempt to read stdin"},
			{Name: "--help", Description: "Show the help message and exit"},
			{Name: "--version", Description: "Show version and exit"},
			{Name: "--traceback", Description: "Prints the exception traceback should one occur"},
			{Name: "--default-scheme", Description: "The default scheme to use if not specified in the URL"},
		},
	})
}
