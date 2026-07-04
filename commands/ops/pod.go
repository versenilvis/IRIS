package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "pod",
		Description: "CocoaPods, the Cocoa library package manager",
		Subcommands: []core.Subcommand{
			{Name: "name", Description: "The name of the podspec file within the Git Repository"},
		},
		Options: []core.Option{
			{Name: "--project-directory", Description: "The path to the root of the project directory"},
			{Name: "--allow-root", Description: "Allows CocoaPods to run as root"},
			{Name: "--all", Description: "Remove all the cached pods without asking"},
			{Name: "--short", Description: "Only print the path relative to the cache root"},
			{Name: "--update", Description: "Run `pod repo update` before listing"},
			{Name: "--stats", Description: "Show additional stats (like GitHub watchers and forks)"},
			{Name: "--regex", Description: "Interpret the `QUERY` as a regular expression"},
			{Name: "--show-all", Description: "Pick from all versions of the given podspec"},
			{Name: "--quick", Description: "Lint skips checks that would require to download and build the spec"},
			{Name: "--allow-warnings", Description: "Lint validates even if warnings are present"},
			{Name: "--subspec", Description: "Lint validates only the given subspec"},
			{Name: "--no-subspecs", Description: "Lint skips validation of subspecs"},
			{Name: "--no-clean", Description: "Lint leaves the build directory intact for inspection"},
			{Name: "--fail-fast", Description: "Lint stops on the first failing platform or subspec"},
			{Name: "--use-libraries", Description: "Lint uses static libraries to install the spec"},
			{Name: "--use-modular-headers", Description: "Lint uses modular headers during installation"},
			{Name: "--use-static-frameworks", Description: "Lint uses static frameworks during installation"},
			{Name: "--sources", Description: "Lint skips checks that apply only to public specs"},
			{Name: "--swift-version", Description: "Lint skips validating that the pod can be imported"},
			{Name: "--skip-tests", Description: "Lint skips building and running tests during validation"},
			{Name: "--test-specs", Description: "List of test specs to run"},
			{Name: "--analyze", Description: "Validate with the Xcode Static Analysis tool"},
			{Name: "--configuration", Description: "Build using the given configuration (defaults to Release)"},
			{Name: "--repo-update", Description: "Force running `pod repo update` before install"},
			{Name: "--deployment", Description: "Disallow any changes to the Podfile or the Podfile.lock during installation"},
			{Name: "--clean-install", Description: "The path to the root of the project directory"},
			{Name: "--no-repo-update", Description: "Skip running `pod repo update` before install"},
			{Name: "--full", Description: "Search by name  author, and description"},
			{Name: "--template-url", Description: "The URL of the git repo containing a compatible template"},
			{Name: "--external-podspecs", Description: "Lint skips validating that the pod can be imported"},
			{Name: "--simple", Description: "Search only by name"},
			{Name: "--web", Description: "Searches on cocoapods.org"},
			{Name: "--ios", Description: "Restricts the search to Pods supported on iOS"},
			{Name: "--osx", Description: "Restricts the search to Pods supported on macOS"},
			{Name: "--watchos", Description: "Restricts the search to Pods supported on watchOS"},
			{Name: "--tvos", Description: "Restricts the search to Pods supported on tvOS"},
			{Name: "--no-pager", Description: "Do not pipe search results into a pager"},
			{Name: "--no-private", Description: "Lint includes checks that apply only to public repos"},
			{Name: "--skip-import-validation", Description: "Lint skips validating that the pod can be imported"},
			{Name: "--commit-message", Description: "Convert the podspec to JSON before pushing it to the repo"},
		},
	})
}
