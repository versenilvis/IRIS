package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "asciinema",
		Description: "Terminal session recorder",
		Subcommands: []spec.Subcommand{
			{Name: "rec", Description: "Start a recording"},
			{Name: "play", Description: "Replay recorded asciicast in a terminal"},
			{Name: "seconds", Description: "Can be fractional"},
			{Name: "factor", Description: "Can be fractional"},
			{Name: "cat", Description: "Print full output of recorded asciicast to a terminal"},
			{Name: "upload", Description: "Upload recorded asciicast to asciinema.org site"},
			{Name: "auth", Description: "Link and manage your install ID with your asciinema.org user account"},
		},
		Options: []spec.Option{
			{Name: "--version", Description: "Output version information and exit"},
			{Name: "-h", Description: "Output help message and exit"},
			{Name: "--stdin", Description: "Enable stdin (keyboard) recording"},
			{Name: "--append", Description: "Append to existing recording"},
			{Name: "--raw", Description: "Save raw output, without timing or other metadata"},
			{Name: "--overwrite", Description: "Overwrite the recording if it already exists"},
			{Name: "-c", Description: "Specify command to record, defaults to $SHELL"},
			{Name: "-e", Description: "List of environment variables to capture"},
			{Name: "-t", Description: "Specify the title of the asciicast"},
			{Name: "-i", Description: "Limit recorded terminal inactivity to max amount of seconds"},
			{Name: "--cols", Description: "Override terminal columns for recorded process"},
			{Name: "--rows", Description: "Override terminal rows for recorded process"},
			{Name: "-y", Description: "Answer “yes” to all prompts (e.g. upload confirmation)"},
			{Name: "-q", Description: "Be quiet, suppress all notices/warnings (implies -y)"},
			{Name: "-s", Description: "Playback speed"},
		},
	})
}
