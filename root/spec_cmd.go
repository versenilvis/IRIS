package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/spec"
)

var SpecCmd = &cobra.Command{
	Use:   "spec",
	Short: "manage custom completion specs",
}

var SpecInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "initialize a new completion spec file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		specsDir := spec.GetUserSpecsDir()
		if specsDir == "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "failed to resolve specs directory")
			return
		}

		name := "example"
		if len(args) > 0 {
			trimmed := strings.TrimSpace(args[0])
			if trimmed != "" {
				name = trimmed
			}
		}

		baseName := filepath.Base(name)
		stem := strings.TrimSuffix(strings.TrimSuffix(baseName, ".jsonc"), ".json")
		if stem == "" || stem == "." {
			stem = "example"
		}
		filePath := filepath.Join(specsDir, stem+".jsonc")

		if _, err := os.Stat(filePath); err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "spec file already exists at %s\n", filePath)
			return
		}

		if err := os.MkdirAll(specsDir, 0755); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to create specs directory: %v\n", err)
			return
		}

		content := fmt.Sprintf(`{
  "name": "%s",
  "description": "%s command completion",
  "options": [
    {
      "name": "--help",
      "aliases": ["-h"],
      "description": "Show help message"
    }
  ],
  "subcommands": [
    {
      "name": "run",
      "description": "Run target",
      "generator": "file"
    }
  ]
}
`, stem, stem)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to write spec file: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "initialized spec file at %s\n", filePath)
	},
}

func init() {
	SpecCmd.AddCommand(SpecInitCmd)
	rootCmd.AddCommand(SpecCmd)
}
