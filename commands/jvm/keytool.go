package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "keytool",
		Description: "Show help message",
		Options: []core.Option{
			{Name: "-h", Description: "Show help message"},
			{Name: "-v", Description: "Verbose output"},
			{Name: "-alias", Description: "Alias name of the entry to process"},
			{Name: "-keystore", Description: "Keystore name"},
			{Name: "-storepass", Description: "Keystore password"},
			{Name: "-storetype", Description: "Keystore type"},
			{Name: "-providername", Description: "Provider name"},
			{Name: "-addprovider", Description: "Add security provider by name (e.g. SunPKCS11)"},
			{Name: "-providerclass", Description: "Add security provider by fully-qualified class name"},
			{Name: "-providerarg", Description: "Configure argument for -addprovider or -providerclass"},
			{Name: "-providerpath", Description: "Provider classpath"},
			{Name: "-protected", Description: "Password through protected mechanism"},
			{Name: "-conf", Description: "Specify pre-configured options file"},
			{Name: "-certreq", Description: "Generates a certificate request"},
			{Name: "-sigalg", Description: "Signature algorithm name"},
			{Name: "-file", Description: "Output file name"},
			{Name: "-keypass", Description: "Key password"},
			{Name: "-dname", Description: "Distinguished name"},
			{Name: "-ext", Description: "X.509 extension"},
			{Name: "-changealias", Description: "Changes an entry's alias"},
			{Name: "-destalias", Description: "Destination alias"},
			{Name: "-cacerts", Description: "Access the cacerts keystore"},
			{Name: "-delete", Description: "Deletes an entry"},
			{Name: "-exportcert", Description: "Exports certificate"},
			{Name: "-rfc", Description: "Output in RFC style"},
			{Name: "-genkeypair", Description: "Generate a key pair"},
			{Name: "-keyalg", Description: "Key algorithm name"},
			{Name: "-keysize", Description: "Key bit size"},
			{Name: "-groupname", Description: "Group name. For example, an Elliptic Curve name"},
			{Name: "-startdate", Description: "Certificate validity start date/time"},
			{Name: "-validity", Description: "Validity number of days"},
			{Name: "-genseckey", Description: "Generates a secret key"},
			{Name: "-gencert", Description: "Generates certificate from a certificate request"},
			{Name: "-infile", Description: "Input file name"},
			{Name: "-outfile", Description: "Output file name"},
			{Name: "-importcert", Description: "Imports a certificate or a certificate chain"},
			{Name: "-noprompt", Description: "Do not prompt"},
			{Name: "-trustcacerts", Description: "Trust certificates from cacerts"},
			{Name: "-importpass", Description: "Imports a password"},
			{Name: "-importkeystore", Description: "Imports one or all entries from another keystore"},
		},
	})
}
