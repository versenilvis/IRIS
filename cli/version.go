package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current Iris version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("iris %s\n", Version)
	},
}
