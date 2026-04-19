package root

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh|fish]",
	Short: "Generate the autostart script for your shell",
	Long: `Add the output of this command to your shell's configuration file to start Iris automatically.
For example, add this to your ~/.zshrc:
  eval "$(iris init zsh)"`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		switch shell {
		case "bash", "zsh":
			fmt.Printf(`
# Iris Autostart Hook
if [ -z "$IRIS_PID" ]; then
    export IRIS_ACTIVE_SHELL="%s"
    exec iris
fi
`, shell)
		case "fish":
			fmt.Printf(`
# Iris Autostart Hook
if not set -q IRIS_PID
    set -gx IRIS_ACTIVE_SHELL "fish"
    exec iris
end
`)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
