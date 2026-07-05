package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "serve",
		Description: "Static file serving and directory listing",
		Options: []spec.Option{
			{Name: "-h", Description: "Shows help message"},
			{Name: "-v", Description: "Displays the current version of serve"},
			{Name: "-l", Description: "E.g. serve --listen 1234"},
			{Name: "-p", Description: "Specify custom port"},
			{Name: "-d", Description: "Show debugging information"},
			{Name: "-s", Description: "Rewrite all not-found requests to `index.html`"},
			{Name: "-c", Description: "Specify custom path to `serve.json`"},
			{Name: "-C", Description: "Enable CORS, sets `Access-Control-Allow-Origin` to `*`"},
			{Name: "-n", Description: "Do not copy the local address to the clipboard"},
			{Name: "-u", Description: "Do not compress files"},
			{Name: "--no-etag", Description: "Send `Last-Modified` header instead of `ETag`"},
			{Name: "-S", Description: "Resolve symlinks instead of showing 404 errors"},
			{Name: "--ssl-cert", Description: "Optional path to an SSL/TLS certificate to serve with HTTPS"},
			{Name: "--ssl-key", Description: "Optional path to the SSL/TLS certificate's private key"},
			{Name: "--ssl-pass", Description: "Optional path to the SSL/TLS certificate's passphrase"},
			{Name: "--no-port-switching", Description: "Do not open a port other than the one specified when it's taken"},
		},
	})
}
