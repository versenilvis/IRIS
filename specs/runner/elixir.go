package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "elixir",
		Description: "Elixir Language",
		Options: []core.Option{
			{Name: "-e", Description: "Evaluates the given command"},
			{Name: "-h", Description: "Prints this message and exits"},
			{Name: "-r", Description: "Requires the given files/patterns"},
			{Name: "-S", Description: "Finds and executes the given script in $PATH"},
			{Name: "-pr", Description: "Requires the given files/patterns in parallel"},
			{Name: "-pa", Description: "Prepends the given path to Erlang code path"},
			{Name: "-pz", Description: "Appends the given path to Erlang code path"},
			{Name: "-v", Description: "Prints Elixir version and exits"},
			{Name: "--app", Description: "Starts the given app and its dependencies"},
			{Name: "--erl", Description: "Switches to be passed down to Erlang"},
			{Name: "--logger-otp-reports", Description: "Enables or disables OTP reporting"},
			{Name: "--logger-sasl-reports", Description: "Enables or disables SASL reporting"},
			{Name: "--no-halt", Description: "Does not halt the Erlang VM after execution"},
			{Name: "--werl", Description: "Uses Erlang's Windows shell GUI (Windows only)"},
			{Name: "--cookie", Description: "Sets a cookie for this distributed node"},
			{Name: "--hidden", Description: "Makes a hidden node"},
			{Name: "--name", Description: "Makes and assigns a name to the distributed node"},
			{Name: "--rpc-eval", Description: "Evaluates the given command on the given remote node"},
			{Name: "--sname", Description: "Makes and assigns a short name to the distributed node"},
			{Name: "--boot", Description: "Uses the given FILE.boot to start the system"},
			{Name: "--boot-var", Description: "Makes $VAR available as VALUE to FILE.boot"},
			{Name: "--erl-config", Description: "Loads configuration in FILE.config written in Erlang"},
			{Name: "--pipe-to", Description: "Starts the Erlang VM as a named PIPEDIR and LOGDIR"},
			{Name: "--vm-args", Description: "Passes the contents in file as arguments to the VM"},
		},
	})
}
