package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rubocop",
		Description: "Run only lint cops",
		Options: []spec.Option{
			{Name: "-l", Description: "Run only lint cops"},
			{Name: "-x", Description: "Run only layout cops, with autocorrect on"},
			{Name: "--safe", Description: "Run only safe cops"},
			{Name: "--except", Description: "Exclude the given cop(s)"},
			{Name: "--only", Description: "Run only the given cop(s)"},
			{Name: "--only-guide-cops", Description: "Run only cops for rules that link to a style guide"},
			{Name: "-F", Description: "Run without pending cops"},
			{Name: "--enable-pending-cops", Description: "Run with pending cops"},
			{Name: "--ignore-disable-comments", Description: "Run cops even when they are disabled locally by a `rubocop:disable` directive"},
			{Name: "--force-exclusion", Description: "Prevent from inheriting `AllCops/Exclude` from parent folders"},
			{Name: "--ignore-unrecognized-cops", Description: "Ignore unrecognized cops or departments in the config"},
			{Name: "--force-default-config", Description: "Use available CPUs to execute inspection in parallel default true"},
			{Name: "--no-parallel", Description: "Disable parallel inspection (default: false)"},
			{Name: "--fail-level", Description: "Minimum severity for exit with error code"},
			{Name: "-C", Description: "Do not use server even if it's available"},
			{Name: "--server", Description: "Restart server process"},
			{Name: "--start-server", Description: "Start server process"},
			{Name: "--stop-server", Description: "Stop server process"},
			{Name: "--server-status", Description: "Show server status"},
			{Name: "-f", Description: "Display cop names in offense messages.  Default is true"},
			{Name: "--no-display-cop-names", Description: "Disable displaying cop names in offense messages. Default false"},
			{Name: "-E", Description: "Display extra details in offense messages"},
			{Name: "-S", Description: "Display style guide URLs in offense messages"},
			{Name: "-o", Description: "Display elapsed time in seconds"},
			{Name: "--display-only-failed", Description: "Only output offense messages. Omit passing cops. Only valid for --format junit"},
			{Name: "--display-only-correctable", Description: "Only output correctable offense messages"},
			{Name: "-A", Description: "Autocorrect offenses (safe and unsafe)"},
			{Name: "--disable-uncorrectable", Description: "Generate a configuration file acting as a TODO list"},
			{Name: "--regenerate-todo", Description: "List all files RuboCop will inspect"},
			{Name: "--show-cops", Description: "Display url to documentation for the given cops, or base url by default"},
			{Name: "--init", Description: "Generate a .rubocop.yml file in the current directory"},
			{Name: "-c", Description: "Specify configuration file"},
			{Name: "-d", Description: "Display debug info"},
			{Name: "-r", Description: "Require Ruby file"},
			{Name: "--no-color", Description: "Disable color output"},
			{Name: "--color", Description: "Force color output"},
			{Name: "-v", Description: "Display version"},
			{Name: "-V", Description: "Display verbose version"},
		},
	})
}
