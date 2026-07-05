package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "tuist",
		Description: "Build the project in the current directory",
		Subcommands: []spec.Subcommand{
			{Name: "build", Description: "Build the project in the current directory"},
			{Name: "warm", Description: "Warms the local and remote cache"},
			{Name: "print-hashes", Description: "Print the hashes of the cacheable frameworks in the given project"},
			{Name: "template", Description: "Name of template you want to use"},
			{Name: "list", Description: "Lists available scaffold templates"},
		},
		Options: []spec.Option{
			{Name: "--build-output-path", Description: "Build the project to a specific device, example usage: 'tuist build --device \\"},
			{Name: "--os", Description: "Build the project to a specific OS, example usage: 'tuist build --os 14.0':"},
			{Name: "--path", Description: "The path to the directory that contains the project whose targets will be cached"},
			{Name: "--profile", Description: "The name of the profile to be used when warming up the cache"},
			{Name: "--xcframeworks", Description: "When passed it caches the targets for simulator and device using xcframeworks"},
			{Name: "--dependencies-only", Description: "Print the hashes of the cacheable frameworks in the given project"},
			{Name: "--skip-external-dependencies", Description: "Excludes external dependencies from the generated graph"},
			{Name: "--format", Description: "If set, the generated graph is not opened automatically. Default is yes"},
			{Name: "--platform", Description: "Don't open the project after generating it. Default is false"},
			{Name: "--xcodeproj-path", Description: "Required. Path to the Xcode project whose build settings will be extracted"},
			{Name: "--xcconfig-path", Description: "Required. Path to the Xcode project whose build settings will be extracted"},
			{Name: "--target", Description: "Required. Path to the Xcode project whose build settings will be checked"},
			{Name: "--name", Description: "The name of the generate project"},
			{Name: "--clean", Description: "When passed, it cleans the project before testing it"},
			{Name: "--device", Description: "Test on a specific device"},
			{Name: "--configuration", Description: "The configuration to be used when testing the scheme"},
			{Name: "--skip-ui-tests", Description: "When passed, it skips testing UI Tests targets"},
			{Name: "--result-bundle-path", Description: "Path where test result bundle will be saved"},
			{Name: "--retry-count", Description: "Show help for tuist"},
		},
	})
}
