package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "streamlit",
		Description: "Streamlit",
		Subcommands: []core.Subcommand{
			{Name: "activate", Description: "Activate Streamlit by entering your email"},
			{Name: "cache", Description: "Manage the Streamlit cache"},
			{Name: "clear", Description: "Clear st.cache, st.memo, and st.singleton caches"},
			{Name: "config", Description: "Manage Streamlit's config settings"},
			{Name: "show", Description: "Show all of Streamlit's config settings"},
			{Name: "docs", Description: "Show help in browser"},
			{Name: "hello", Description: "Runs the Hello World script"},
			{Name: "help", Description: "Print the help message"},
			{Name: "run", Description: "Run a Python script, piping stderr to Streamlit"},
			{Name: "file", Description: "The Python script to run"},
			{Name: "version", Description: "Print Streamlit's version number"},
		},
		Options: []core.Option{
			{Name: "--log_level", Description: "Set the log level"},
			{Name: "--help", Description: "Show a help message and exit"},
			{Name: "--version", Description: "Show the version and exit"},
		},
	})
}
