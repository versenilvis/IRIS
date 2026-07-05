package js

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "yalc",
		Description: "Work with yarn/npm packages locally like a boss",
		Subcommands: []core.Subcommand{
			{Name: "publish", Description: "Copy all the files that should be published in remote NPM registry"},
			{Name: "add", Description: "Copy the current version from the store to your project"},
			{Name: "package", Description: "The package you want to add"},
			{Name: "link", Description: "Alternative to 'add', instead using local '.yalc' as symlink source"},
			{Name: "update", Description: "Update package(s)"},
			{Name: "remove", Description: "Remove package info from 'package.json' & 'yalc.lock'"},
			{Name: "clean", Description: "Unpublish a package published with yalc publish"},
			{Name: "show", Description: "Show all packages to which chosen package has been added"},
			{Name: "dir", Description: "Show yalc system directory"},
			{Name: "check", Description: "Check 'package.json' for yalc packages"},
			{Name: "restore", Description: "Restore retreated packages"},
			{Name: "retreat", Description: "Remove packages from project, but leave in lock file (to be restored later)"},
		},
		Options: []core.Option{
			{Name: "--push", Description: "Publish without running scripts"},
			{Name: "--no-sig", Description: "Disable adding hash signature of all files when copying package content"},
			{Name: "--content", Description: "Show included files in the published package"},
			{Name: "--no-workspace-resolve", Description: "Do not resolve 'workspace:' protocol in dependencies"},
			{Name: "--private", Description: "Force publishing of private package"},
			{Name: "--link", Description: "Add a 'link:' dependency instead of 'file:'"},
			{Name: "--dev", Description: "Add yalc package to dev dependencies"},
			{Name: "--pure", Description: "Do not touch 'package.json' or 'node_modules'"},
			{Name: "--workspace", Description: "Add dependency with 'workspace:' protocol"},
			{Name: "--update", Description: "Run package manager's update command for packages"},
			{Name: "--all", Description: "Remove all packages from project"},
			{Name: "--help", Description: "Show help for yalc"},
			{Name: "--no-colors", Description: "Disable colors"},
			{Name: "--quiet", Description: "Fully disable output (except errors)"},
		},
	})
}
