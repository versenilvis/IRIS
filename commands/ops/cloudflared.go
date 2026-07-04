package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "cloudflared",
		Description: "Specify the hostname of your application",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "Install cloudflared as an user launch agent"},
			{Name: "uninstall", Description: "Uninstall the cloudflared launch agent"},
			{Name: "help", Description: "Shows a list of commands or help for one command"},
		},
		Options: []core.Option{
			{Name: "--hostname", Description: "Specify the hostname of your application"},
			{Name: "--destination", Description: "Specify the destination address of your SSH server"},
			{Name: "--url", Description: "Specify the host:port to forward data to Cloudflare edge"},
			{Name: "--service-token-id", Description: "Specify an Access service token ID you wish to use. [$TUNNEL_SERVICE_TOKEN_ID]"},
			{Name: "--service-token-secret", Description: "Access service token secret"},
			{Name: "--log-directory", Description: "Save application log to this directory for reporting issues"},
			{Name: "--log-level", Description: "Application logging level {debug, info, warn, error, fatal}"},
			{Name: "--beta", Description: "Specify if you wish to update to the latest beta version (default: false)"},
			{Name: "--version", Description: "Specify a version you wish to upgrade or downgrade to"},
			{Name: "--metrics", Description: "Manages the cloudflared launch agent"},
			{Name: "--allow-request", Description: "The token subcommand produces a JWT which can be used to authenticate requests"},
			{Name: "--app", Description: "Url of access application"},
			{Name: "--short-lived-cert", Description: "Specify if you wish to generate short lived certs. (default: false)"},
			{Name: "--config", Description: "Config file"},
			{Name: "--origincert", Description: "Certificate generated"},
			{Name: "--autoupdate-freq", Description: "Autoupdate frequency. Default is 24h0m0s. (default: 24h0m0s)"},
			{Name: "--no-autoupdate", Description: "Adress"},
			{Name: "--pidfile", Description: "Application PID"},
			{Name: "--loglevel", Description: "Logging level"},
			{Name: "--transport-loglevel", Description: "Logging level"},
			{Name: "--logfile", Description: "Save application log to this file for reporting issues. [$TUNNEL_LOGFILE]"},
			{Name: "--trace-output", Description: "Output file"},
			{Name: "--output", Description: "Render output using given FORMAT. Valid options are 'json' or 'yaml'"},
			{Name: "--credentials-file", Description: "Filepath at which to read/write the tunnel credentials [$TUNNEL_CRED_FILE]"},
			{Name: "--secret", Description: "Base 64 encoded secret"},
			{Name: "--overwrite-dns", Description: "Tunnel"},
			{Name: "--vnet", Description: "The ID or name of the virtual network to which the route is associated to"},
			{Name: "--filter-tunnel-id", Description: "Show only routes with the given tunnel ID"},
			{Name: "--filter-network-is-subset-of", Description: "Show only routes whose network is a subset of the given network"},
			{Name: "--filter-comment-is", Description: "Show only routes with this comment"},
			{Name: "--filter-vnet-id", Description: "Show only routes that are attached to the given virtual network ID"},
			{Name: "--default", Description: "Lists the virtual networks based on the given filter flags"},
			{Name: "--id", Description: "List virtual networks with the given ID"},
			{Name: "--name", Description: "List virtual networks with the given NAME"},
			{Name: "--is-default", Description: "Render output using given FORMAT. Valid options are 'json' or 'yaml'"},
			{Name: "--comment", Description: "A new comment describing the purpose of the virtual network"},
			{Name: "--force", Description: "Filepath at which to read/write the tunnel credentials [$TUNNEL_CRED_FILE]"},
			{Name: "--credentials-contents", Description: "Contents of the tunnel credentials JSON"},
			{Name: "--features", Description: "Features"},
			{Name: "--token", Description: "Tunnel token"},
		},
	})
}
