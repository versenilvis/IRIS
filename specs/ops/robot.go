package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "robot",
		Description: "Tag",
		Options: []core.Option{
			{Name: "-h", Description: "Print usage instructions"},
			{Name: "--rpa", Description: "Parse only files with this extension when executing a directory"},
			{Name: "-N", Description: "Set a name of the top level suite"},
			{Name: "-D", Description: "Set a documentation of the top level suite"},
			{Name: "-M", Description: "Set metadata of the top level suite"},
			{Name: "-G", Description: "Sets given tag to all executed tests"},
			{Name: "-t", Description: "Alias to --test. Especially applicable with --rpa"},
			{Name: "-s", Description: "Select suites by name"},
			{Name: "-i", Description: "Select test cases by tag"},
			{Name: "-e", Description: "Select test cases not to run by tag"},
			{Name: "-R", Description: "Select failed tests from an earlier output file to be re-executed"},
			{Name: "-S", Description: "Select failed suites from an earlier output file to be re-executed"},
			{Name: "--runemptysuite", Description: "Executes suite even if it contains no tests"},
			{Name: "--skip", Description: "Tests having given tag will be skipped"},
			{Name: "--skiponfailure", Description: "Tests having given tag will be skipped if they fail"},
			{Name: "-v", Description: "Set variables in the test data"},
			{Name: "-V", Description: "Python or YAML file file to read variables from"},
			{Name: "-d", Description: "XUnit compatible result file. Not created unless this option is specified"},
			{Name: "-b", Description: "Debug file written during execution. Not created unless this option is specified"},
			{Name: "-T", Description: "Split the log file into smaller pieces that open in browsers transparently"},
			{Name: "--logtitle", Description: "Title for the generated log file. The default title is `<SuiteName> Log.`"},
			{Name: "--reporttitle", Description: "Title for the generated report file. The default title is `<SuiteName> Report`"},
			{Name: "--reportbackground", Description: "Default number of lines"},
			{Name: "--maxassignlength", Description: "Default number of characters"},
			{Name: "-L", Description: "Threshold level for logging"},
			{Name: "--suitestatlevel", Description: "How many levels to show in `Statistics by Suite` in log and report"},
			{Name: "--tagstatinclude", Description: "Include only matching tags in `Statistics by Tag` in log and report"},
			{Name: "--tagstatexclude", Description: "Exclude matching tags from `Statistics by Tag`"},
			{Name: "--tagstatcombine", Description: "Add documentation to tags matching the given pattern"},
			{Name: "--tagstatlink", Description: "Matching keywords will be automatically expanded in the log file"},
			{Name: "--removekeywords", Description: "Remove keyword data from the generated log file"},
			{Name: "--flattenkeywords", Description: "Flattens matching keywords in the generated log file"},
			{Name: "--listener", Description: "Stops test execution if any test fails"},
			{Name: "--exitonerror", Description: "Causes teardowns to be skipped if test execution is stopped prematurely"},
			{Name: "--randomize", Description: "Randomizes the test execution order"},
			{Name: "--prerunmodifier", Description: "Class to programmatically modify the suite structure before execution"},
			{Name: "--prerebotmodifier", Description: "How to report execution on the console"},
			{Name: "-.", Description: "Shortcut for `--console dotted`"},
			{Name: "--quiet", Description: "Shortcut for `--console quiet`"},
			{Name: "-W", Description: "Width of the console output. Default is 78"},
		},
	})
}
