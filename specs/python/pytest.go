package python

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pytest",
		Description: "Control assertion debugging tools.",
		Options: []core.Option{
			{Name: "--assert", Description: "Remove all cache contents at start of test run"},
			{Name: "--cache-show", Description: "Per-test capturing method"},
			{Name: "--code-highlight", Description: "Whether code should be highlighted (only if --color is also enabled)"},
			{Name: "--co", Description: "Only collect tests, don't execute them"},
			{Name: "--collect-in-virtualenv", Description: "Don't ignore tests in a local virtualenv directory"},
			{Name: "--color", Description: "Color terminal output"},
			{Name: "--confcutdir", Description: "Only load conftest.py's relative to specified dir"},
			{Name: "--debug", Description: "Show N slowest setup/test durations (N=0 for all)"},
			{Name: "--durations-min", Description: "Minimal duration in seconds for inclusion in slowest list"},
			{Name: "--deselect", Description: "Deselect item (via node id prefix) during collection (multi-allowed)"},
			{Name: "--disable-warnings", Description: "Disable warnings summary"},
			{Name: "--doctest-continue-on-failure", Description: "For a given doctest, continue to run after the first failure"},
			{Name: "--doctest-modules", Description: "Run doctests in all .py modules"},
			{Name: "--doctest-report", Description: "Choose another output format for diffs on doctest failure"},
			{Name: "--doctest-glob", Description: "Doctests file matching pattern, default: test*.txt"},
			{Name: "--exitfirst", Description: "Exit instantly on first error or failed test"},
			{Name: "--failed-first", Description: "Run all tests, but run the last failures first"},
			{Name: "--fixtures", Description: "Show fixtures per test"},
			{Name: "--full-trace", Description: "Don't cut any tracebacks (default is to cut)"},
			{Name: "--help", Description: "This shows help on command line and config-line options"},
			{Name: "--ignore", Description: "Ignore path during collection (multi-allowed)"},
			{Name: "--ignore-glob", Description: "Ignore path pattern during collection (multi-allowed)"},
			{Name: "--import-mode", Description: "Create junit-xml style report file at given path"},
			{Name: "--junit-prefix", Description: "Prepend prefix to classnames in junit-xml output"},
			{Name: "-k", Description: "Ex: 'test_method or test_other'"},
			{Name: "--keep-duplicates", Description: "Keep duplicate tests"},
			{Name: "--showlocals", Description: "Show locals in tracebacks (disabled by default)"},
			{Name: "--last-failed-no-failures", Description: "Which tests to run with no previously (known) failures"},
			{Name: "--last-failed", Description: "Rerun only the tests that failed at the last run (or all if none failed)"},
			{Name: "--log-auto-indent", Description: "Cli logging level"},
			{Name: "--log-cli-format", Description: "Log format as used by the logging module"},
			{Name: "--log-cli-date-format", Description: "Log date format as used by the logging module"},
			{Name: "--log-date-format", Description: "Log date format as used by the logging module"},
			{Name: "--log-format", Description: "Log format as used by the logging module"},
			{Name: "--log-file", Description: "Path to a file where logging will be written to"},
			{Name: "--log-file-level", Description: "Log file logging level"},
			{Name: "--log-file-date-format", Description: "Log date format as used by the logging module"},
			{Name: "--log-file-format", Description: "Log format as used by the logging module"},
			{Name: "--log-level", Description: "Only run tests matching given mark expression"},
			{Name: "--markers", Description: "Show markers (builtin, plugin and per-project ones)"},
		},
	})
}
