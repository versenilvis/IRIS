package ops

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "ssh-keygen",
		Description: "Generates, manages and converts authentication keys for ssh",
		Options: []spec.Option{
			{Name: "-A", Description: "When saving a private key, this option specifies the number of KDF"},
			{Name: "-B", Description: "Show the bubblebabble digest of specified private or public key file"},
			{Name: "-b", Description: "Specifies the number of bits in the key to create"},
			{Name: "-C", Description: "Provides a new comment"},
			{Name: "-c", Description: "Requests changing the comment in the private and public key files"},
			{Name: "-D", Description: "Download the public keys provided by the PKCS#11"},
			{Name: "-E", Description: "Specifies the hash algorithm used"},
			{Name: "-e", Description: "Read a OpenSSH key file and print to stdout"},
			{Name: "-F", Description: "Search for the specified hostname (with optional port number)"},
			{Name: "-f", Description: "Specifies the filename of the key file"},
			{Name: "-g", Description: "Use generic DNS format when printing fingerprint resource records"},
			{Name: "-H", Description: "Hash a known_hosts file"},
			{Name: "-h", Description: "Create a host certificate instead of a user"},
			{Name: "-I", Description: "Specify the key identity when signing a public key"},
			{Name: "-i", Description: "Read an unencrypted private (or public) key file"},
			{Name: "-K", Description: "Download resident keys from a FIDO	authenticator"},
			{Name: "-k", Description: "Generate a	KRL file"},
			{Name: "-L", Description: "Generate a	KRL file"},
			{Name: "-l", Description: "Show fingerprint of specified public key file"},
			{Name: "-M", Description: "Use for Moduli generation"},
			{Name: "-m", Description: "Specify a key format for key generation"},
			{Name: "-N", Description: "Provides the new passphrase"},
			{Name: "-n", Description: "Principals to be included in a certificate when signing a key"},
			{Name: "-O", Description: "Specify a key/value option"},
			{Name: "-P", Description: "Provides the (old) passphrase"},
			{Name: "-p", Description: "Test whether keys have been revoked in a KRL"},
			{Name: "-q", Description: "Silence ssh-keygen"},
			{Name: "-R", Description: "Removes all keys belonging to hostname"},
			{Name: "-r", Description: "Hostname for the specified public key file"},
			{Name: "-s", Description: "Certify (sign) a public key using the specified CA	key"},
			{Name: "-t", Description: "Specifies the type of key to create"},
			{Name: "-U", Description: "Update a KRL"},
			{Name: "-V", Description: "Specify a validity interval when signing a certificate"},
			{Name: "-v", Description: "Verbose mode"},
			{Name: "-w", Description: "Path to library to be used when creating FIDO authenticator-hosted keys"},
			{Name: "-Y", Description: "Find the principal(s) associated with the public key of a signature"},
			{Name: "-y", Description: "Read a private OpenSSH format file and print an OpenSSH public key to stdout"},
			{Name: "-Z", Description: "Cipher to use for encryption"},
			{Name: "-z", Description: "Serial number to distinguish this certificate"},
		},
	})
}
