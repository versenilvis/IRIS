package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "phpunit",
		Description: "Generate code coverage report in Clover XML format,",
		Options: []core.Option{
			{Name: "--coverage-clover", Description: "Generate code coverage report in Clover XML format,"},
			{Name: "--coverage-crap4j", Description: "Generate code coverage report in Crap4J XML format"},
			{Name: "--coverage-html", Description: "Generate code coverage report in HTML format"},
			{Name: "--coverage-php", Description: "Export PHP_CodeCoverage object to file"},
			{Name: "--coverage-text", Description: "Generate code coverage report in text format [default: standard output]"},
			{Name: "--coverage-xml", Description: "Generate code coverage report in PHPUnit XML format"},
			{Name: "--coverage-cache", Description: "Cache static analysis results"},
			{Name: "--warm-coverage-cache", Description: "Warm static analysis cache"},
			{Name: "--coverage-filter", Description: "Include <dir> in code coverage analysis"},
			{Name: "--path-coverage", Description: "Perform path coverage analysis"},
			{Name: "--disable-coverage-ignore", Description: "Disable annotations for ignoring code coverage"},
			{Name: "--no-coverage", Description: "Ignore code coverage configuration"},
			{Name: "--dont-report-useless-tests", Description: "Do not report tests that do not test anything"},
			{Name: "--strict-coverage", Description: "Be strict about @covers annotation usage"},
			{Name: "--strict-global-state", Description: "Be strict about changes to global state"},
			{Name: "--disallow-test-output", Description: "Be strict about output during tests"},
			{Name: "--disallow-resource-usage", Description: "Be strict about resource usage during small tests"},
			{Name: "--enforce-time-limit", Description: "Enforce time limit based on test size"},
			{Name: "--default-time-limit", Description: "Timeout in seconds for tests without @small, @medium or @large"},
			{Name: "--disallow-todo-tests", Description: "Disallow @todo-annotated tests"},
			{Name: "--log-junit", Description: "Log test execution in JUnit XML format to file"},
			{Name: "--log-teamcity", Description: "Log test execution in TeamCity format to file"},
			{Name: "--testdox-html", Description: "Write agile documentation in HTML format to file"},
			{Name: "--testdox-text", Description: "Write agile documentation in Text format to file"},
			{Name: "--testdox-xml", Description: "Write agile documentation in HTML format to file"},
			{Name: "--reverse-list", Description: "Print defects in reverse order"},
			{Name: "--no-logging", Description: "Ignore logging configuration"},
			{Name: "--prepend", Description: "A PHP script that is included as early as possible"},
			{Name: "--bootstrap", Description: "A PHP script that is included before the tests run"},
			{Name: "-c", Description: "Read configuration from XML file"},
			{Name: "--no-configuration", Description: "Ignore default configuration file (phpunit.xml)"},
			{Name: "--extensions", Description: "A comma separated list of PHPUnit extensions to load"},
			{Name: "--no-extensions", Description: "Do not load PHPUnit extensions"},
			{Name: "--include-path", Description: "Prepend PHP's include_path with given path(s)"},
			{Name: "-d", Description: "Sets a php.ini value"},
			{Name: "--cache-result-file", Description: "Specify result cache path and filename"},
			{Name: "--generate-configuration", Description: "Generate configuration file with suggested settings"},
			{Name: "--migrate-configuration", Description: "Migrate configuration file to current format"},
		},
	})
}
