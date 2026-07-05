package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "whois",
		Description: "Query a database for information about a domain registrant",
		Options: []spec.Option{
			{Name: "-a", Description: "Use the American Registry for Internet Numbers (ARIN) database"},
			{Name: "-A", Description: "Use the Asia/Pacific Network Information Center (APNIC) database"},
			{Name: "-b", Description: "Use the Network Abuse Clearinghouse database"},
			{Name: "-c", Description: "Equivalent to '-h TLD.whois-servers.net', where 'TLD' is this option's argument"},
			{Name: "-f", Description: "Use the African Network Information Centre (AfriNIC) database"},
			{Name: "-g", Description: "Use the US non-military federal government database"},
			{Name: "-h", Description: "Use the specified host instead of the default (host name or IP)"},
			{Name: "-i", Description: "Use the traditional Network Information Center (InterNIC) database"},
			{Name: "-I", Description: "Use the Internet Assigned Numbers Authority (IANA) database"},
			{Name: "-k", Description: "Use the National Internet Development Agency of Korea (KRNIC) database"},
			{Name: "-l", Description: "Use the Route Arbiter Database (RADB) database"},
			{Name: "-p", Description: "Connect to the whois server on the given port"},
			{Name: "-P", Description: "Use the PeeringDB database of AS numbers"},
			{Name: "-Q", Description: "Do a quick lookup (don't follow referrals)"},
			{Name: "-r", Description: "Use the Réseaux IP Européens (RIPE) database"},
			{Name: "-R", Description: "Do a recursive lookup"},
			{Name: "-S", Description: "Print the output verbatim"},
		},
	})
}
