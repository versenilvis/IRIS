package cc

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "xcodebuild",
		Description: "Build Xcode projects",
		Options: []spec.Option{
			{Name: "-usage", Description: "Print brief usage"},
			{Name: "-help", Description: "Print complete usage"},
			{Name: "-verbose", Description: "Provide additional status output"},
			{Name: "-license", Description: "Show the Xcode and SDK license agreements"},
			{Name: "-checkFirstLaunchStatus", Description: "Check if any First Launch tasks need to be performed"},
			{Name: "-runFirstLaunch", Description: "Install packages and agree to the license"},
			{Name: "-project", Description: "Build the project NAME"},
			{Name: "-target", Description: "Build the target NAME"},
			{Name: "-alltargets", Description: "Build all targets"},
			{Name: "-workspace", Description: "Build the workspace NAME"},
			{Name: "-scheme", Description: "Build the scheme NAME"},
			{Name: "-configuration", Description: "Use the build configuration NAME for building each target"},
			{Name: "-xcconfig", Description: "Apply the build settings defined in the file at PATH as overrides"},
			{Name: "-arch", Description: "Use SDK as the name or path of the base SDK when building the project"},
			{Name: "-toolchain", Description: "Use the toolchain with identifier or name NAME"},
			{Name: "-destination", Description: "Wait for TIMEOUT seconds while searching for the destination device"},
			{Name: "-parallelizeTargets", Description: "Build independent targets in parallel"},
			{Name: "-jobs", Description: "Specify the maximum number of concurrent build operations"},
			{Name: "-parallel-testing-enabled", Description: "Overrides the per-target setting in the scheme"},
			{Name: "-dry-run", Description: "Do everything except actually running the commands"},
			{Name: "-quiet", Description: "Do not print any output except for warnings and errors"},
			{Name: "-hideShellScriptEnvironment", Description: "Don't show shell script environment variables in build log"},
			{Name: "-showsdks", Description: "Display a compact list of the installed SDKs"},
			{Name: "-showdestinations", Description: "Display a list of destinations"},
			{Name: "-showTestPlans", Description: "Display a list of test plans"},
			{Name: "-showBuildSettings", Description: "Display a list of build settings and values"},
			{Name: "-showBuildSettingsForIndex", Description: "Display build settings for the index service"},
			{Name: "-list", Description: "Lists the targets and configurations in a project, or the schemes in a workspace"},
			{Name: "-find-executable", Description: "Display the full path to executable NAME in the provided SDK and toolchain"},
			{Name: "-find-library", Description: "Display the full path to library NAME in the provided SDK and toolchain"},
			{Name: "-version", Description: "Turn the address sanitizer on or off"},
			{Name: "-enableThreadSanitizer", Description: "Turn the thread sanitizer on or off"},
			{Name: "-resultBundlePath", Description: "Specifies which result bundle version should be used"},
			{Name: "-clonedSourcePackagesDirPath", Description: "Specifies the directory where build products and other derived data will go"},
			{Name: "-archivePath", Description: "Specifies that an archive should be exported"},
			{Name: "-exportNotarizedApp", Description: "Export an archive that has been notarized by Apple"},
			{Name: "-exportOptionsPlist", Description: "Specifies a path to a plist file that configures archive exporting"},
			{Name: "-enableCodeCoverage", Description: "Turn code coverage on or off when testing"},
			{Name: "-exportPath", Description: "Specifies the destination for the product exported from an archive"},
			{Name: "-skipUnavailableActions", Description: "Exports completed and outstanding project localizations"},
		},
	})
}
