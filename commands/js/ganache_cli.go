package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "ganache-cli",
		Description: "Fast Ethereum RPC client",
		Options: []core.Option{
			{Name: "-a", Description: "Specify the number of accounts to generate at startup"},
			{Name: "-e", Description: "Amount of ether to assign each test account. Default is 100"},
			{Name: "-b", Description: "Generate deterministic addresses based on a pre-defined mnemonic"},
			{Name: "-n", Description: "Lock available accounts by default (good for third party transaction signing)"},
			{Name: "-m", Description: "Port number to listen on. Defaults to 8545"},
			{Name: "-h", Description: "Use arbitrary data to generate the HD wallet mnemonic to be used"},
			{Name: "-g", Description: "The price of gas in wei (defaults to 20000000000)"},
			{Name: "-l", Description: "The block gas limit (defaults to 0x6691b7)"},
			{Name: "--callGasLimit", Description: "Sets the transaction gas limit for eth_call and eth_estimateGas calls"},
			{Name: "-k", Description: "Allows users to specify which hardfork should be used"},
			{Name: "-f", Description: "Output VM opcodes for debugging"},
			{Name: "--mem", Description: "Output ganache-cli memory usage statistics. This replaces normal output"},
			{Name: "-q", Description: "Run ganache-cli without any logs"},
			{Name: "-?", Description: "Display help information"},
			{Name: "--version", Description: "Display the version of ganache-cli"},
			{Name: "--account_keys_path", Description: "Specifies a file to save accounts and private keys to, for testing"},
			{Name: "--noVMErrorsOnRPCResponse", Description: "Do not transmit transaction failures as RPC errors"},
			{Name: "--allowUnlimitedContractSize", Description: "Allows unlimited contract sizes while debugging"},
			{Name: "--keepAliveTimeout", Description: "Sets the HTTP server's keepAliveTimeout in milliseconds"},
			{Name: "-t", Description: "Date (ISO 8601) that the first block should start"},
		},
	})
}
