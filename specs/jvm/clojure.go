package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "clojure",
		Description: "An alias to refer to its function or a qualified function",
		Options: []core.Option{
			{Name: "-A", Description: "Use concatenated aliases to modify classpath"},
			{Name: "-X", Description: "Invoke tool by name or via aliases ala -X"},
			{Name: "-M", Description: "Use concatenated aliases to modify classpath or supply main opts"},
			{Name: "-P", Description: "Prepare deps - download libs, cache classpath, but don't exec"},
			{Name: "-J", Description: "Pass opt through in java_opts"},
			{Name: "-Sdeps", Description: "Pass the deps data on the command line"},
			{Name: "-Spath", Description: "Compute classpath and echo to stdout only"},
			{Name: "-Scp", Description: "Use specified classpath instead of cached or computed one"},
			{Name: "-Sdescribe", Description: "Print environment and command parsing information as data"},
			{Name: "-Sforce", Description: "Ignore classpath cache and force recomputation"},
			{Name: "-Spom", Description: "Generate (or update) pom.xml with deps and paths"},
			{Name: "-Srepro", Description: "Ignore the ~/.clojure/deps.edn config file"},
			{Name: "-Sthreads", Description: "Set the number of threads to use when downloading dependencies"},
			{Name: "-Strace", Description: "Write a trace.edn file that traces deps expansion"},
			{Name: "-Stree", Description: "Print dependency tree"},
			{Name: "-Sverbose", Description: "Print all path locations"},
			{Name: "-version", Description: "Print the Clojure CLI version"},
			{Name: "-i", Description: "Load a file or resource"},
			{Name: "-e", Description: "Evaluate expressions in string; print non-nil values"},
			{Name: "--report", Description: "Report uncaught exceptions"},
			{Name: "-m", Description: "Call the -main function from a namespace with args"},
			{Name: "-r", Description: "Run a REPL"},
			{Name: "-h", Description: "Show help for clojure"},
		},
	})
}
