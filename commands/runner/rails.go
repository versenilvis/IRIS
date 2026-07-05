package runner

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "rails",
		Description: "Create a new rails application",
		Options: []spec.Option{
			{Name: "-skip-namespace", Description: "Skip namespace (affects only isolated applications)"},
			{Name: "-r", Description: "Path to the Ruby binary of your choice"},
			{Name: "-m", Description: "Path to some application template (can be a filesystem path or URL)"},
			{Name: "-d", Description: "Preconfigure for selected database - defaults to sqlite3"},
			{Name: "--skip-gemfile", Description: "Don't create a Gemfile"},
			{Name: "-G", Description: "Skip .gitignore file"},
			{Name: "--skip-keeps", Description: "Skip source control .keep files"},
			{Name: "-M", Description: "Skip Action Mailer files"},
			{Name: "--skip-action-mailbox", Description: "Skip Action Mailbox gem"},
			{Name: "--skip-action-text", Description: "Skip Action Text gem"},
			{Name: "-O", Description: "Skip Active Record files"},
			{Name: "--skip-active-storage", Description: "Skip Active Storage files"},
			{Name: "-P", Description: "Skip Puma related files"},
			{Name: "-C", Description: "Skip Action Cable files"},
			{Name: "-S", Description: "Skip Sprockets files"},
			{Name: "--skip-spring", Description: "Don't install Spring application preloader"},
			{Name: "--skip-listen", Description: "Don't generate configuration that depends on the listen gem"},
			{Name: "-J", Description: "Skip JavaScript files"},
			{Name: "--skip-turbolinks", Description: "Skip turbolinks gem"},
			{Name: "-T", Description: "Skip test files"},
			{Name: "--skip-system-test", Description: "Skip system test files"},
			{Name: "--skip-bootsnap", Description: "Skip bootsnap gem"},
			{Name: "--dev", Description: "Setup the application with Gemfile pointing to your Rails checkout"},
			{Name: "--edge", Description: "Setup the application with Gemfile pointing to Rails repository"},
			{Name: "--rc", Description: "Path to file containing extra configuration options for rails command"},
			{Name: "--no-rc", Description: "Skip loading of extra configuration options from .railsrc file"},
			{Name: "--api", Description: "Preconfigure smaller stack for API only apps"},
			{Name: "-B", Description: "Don't run bundle install"},
			{Name: "--webpacker", Description: "Preconfigure Webpack with a particular framework"},
			{Name: "--f", Description: "Overwrite files that already exist"},
			{Name: "--p", Description: "Run but do not make any changes"},
			{Name: "--q", Description: "Suppress status output"},
			{Name: "--s", Description: "Skip files that already exist"},
			{Name: "--h", Description: "Show this help message and quit"},
			{Name: "--v", Description: "Show Rails version number and quit"},
			{Name: "--backtrace", Description: "Enable full backtrace.  OUT can be stderr (default) or stdout"},
			{Name: "--comments", Description: "Show commented tasks only"},
			{Name: "--job-stats", Description: "Display job statistics. LEVEL=history displays a complete job list"},
			{Name: "--rules", Description: "Trace the rules resolution"},
			{Name: "--suppress-backtrace", Description: "Suppress backtrace lines matching regexp PATTERN. Ignored if --trace is on"},
		},
	})
}
