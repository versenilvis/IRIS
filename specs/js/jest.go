package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "jest",
		Description: "A delightful JavaScript Testing Framework with a focus on simplicity",
		Options: []core.Option{
			{Name: "--bail", Description: "Whether to use the cache"},
			{Name: "--no-cache", Description: "Whether to use the cache"},
			{Name: "--changedFilesWithAncestor", Description: "Runs tests related to the changes since the provided branch or commit hash"},
			{Name: "--ci", Description: "Deletes the Jest cache directory and then exits without running tests"},
			{Name: "--collectCoverageFrom", Description: "Forces test results output highlighting even if stdout is not a TTY"},
			{Name: "--config", Description: "The path to a Jest config file specifying how to find and execute tests"},
			{Name: "--coverage", Description: "Enable or disable coverage, disabled by default"},
			{Name: "--coverageProvider", Description: "Indicates which provider should be used to instrument code for coverage"},
			{Name: "--debug", Description: "Print debugging info about your Jest config"},
			{Name: "--detectOpenHandles", Description: "Attempt to collect and print open handles preventing Jest from exiting cleanly"},
			{Name: "--env", Description: "The test environment used for all tests"},
			{Name: "--errorOnDeprecated", Description: "Make calling deprecated APIs throw helpful error messages"},
			{Name: "--expand", Description: "Use this flag to show full diffs and errors instead of a patch"},
			{Name: "--findRelatedTests", Description: "Force Jest to exit after all tests have completed running"},
			{Name: "--help", Description: "Show the help information"},
			{Name: "--init", Description: "Generate a basic configuration file"},
			{Name: "--injectGlobals", Description: "Prints the test results in JSON"},
			{Name: "--outputFile", Description: "Write test results to a file when the --json option is also specified"},
			{Name: "--lastCommit", Description: "Run all tests affected by file changes in the last commit made"},
			{Name: "--listTests", Description: "Lists all tests as JSON that Jest will run given the arguments, and exits"},
			{Name: "--logHeapUsage", Description: "Logs the heap usage after every test"},
			{Name: "--maxConcurrency", Description: "Disables stack trace in test results output"},
			{Name: "--notify", Description: "Activates notifications for test results"},
			{Name: "--onlyChanged", Description: "Allows the test suite to pass when no files are found"},
			{Name: "--projects", Description: "Run tests with specified reporters"},
			{Name: "--roots", Description: "A list of paths to directories that Jest should use to search for files in"},
			{Name: "--runInBand", Description: "Run only the tests of the specified projects"},
			{Name: "--runTestsByPath", Description: "Run only the tests that were specified with their exact paths"},
			{Name: "--silent", Description: "Prevent tests from printing messages through the console"},
			{Name: "--testNamePattern", Description: "Run only tests with a name that matches the regex"},
			{Name: "--testLocationInResults", Description: "Adds a location field to test results"},
			{Name: "--testPathPattern", Description: "Lets you specify a custom test runner"},
			{Name: "--testSequencer", Description: "Lets you specify a custom test sequencer"},
			{Name: "--testTimeout", Description: "Default timeout of a test in milliseconds"},
			{Name: "--updateSnapshot", Description: "Use this flag to re-record every snapshot that fails during this test run"},
			{Name: "--useStderr", Description: "Divert all output to stderr"},
			{Name: "--verbose", Description: "Display individual test results with the test suite hierarchy"},
			{Name: "--version", Description: "Print the version and exit"},
			{Name: "--watch", Description: "Watch files for changes and rerun tests related to changed files"},
			{Name: "--watchAll", Description: "Watch files for changes and rerun all tests when something changes"},
		},
	})
}
