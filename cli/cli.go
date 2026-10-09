package cli

import "github.com/spf13/cobra"

var Version = "dev"

func Register(rootCmd *cobra.Command, version string) {
	if version != "" {
		Version = version
	}
	rootCmd.AddCommand(
		ConfigCmd,
		ThemeCmd,
		SpecCmd,
		InitCmd,
		SetupCmd,
		UninstallCmd,
		UpdateCmd,
		VersionCmd,
		ChangelogCmd,
	)
}
