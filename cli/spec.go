package cli

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

		const schemaURL = "https://raw.githubusercontent.com/versenilvis/iris/main/spec/schema.json"
		content := fmt.Sprintf(`{
  // $schema enables autocomplete and validation in VS Code / Zed
  "$schema": "%s",
  "name": "%s",
  "description": "%s command completion",
  "options": [
    {
      "name": "--verbose",
      "aliases": ["-v"],
      "description": "Enable verbose output"
    }
  ],
  "subcommands": [
    {
      "name": "run",
      "description": "Run a file",
      // suggest files by extension
      "generator": [".go", ".js", ".ts", ".py", ".sh"]
    },
    {
      "name": "help",
      "description": "Show help for a command"
    }
  ]
}
`, schemaURL, stem, stem)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to write spec file: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "initialized spec file at %s\n", filePath)
	},
}

var SpecPathCmd = &cobra.Command{
	Use:   "path",
	Short: "show the specs directory path",
	Run: func(cmd *cobra.Command, args []string) {
		dir := spec.GetUserSpecsDir()
		if dir == "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "failed to resolve specs directory")
			return
		}
		_ = os.MkdirAll(dir, 0755)
		fmt.Fprintln(cmd.OutOrStdout(), dir)
	},
}

func init() {
	SpecCmd.AddCommand(SpecInitCmd)
	SpecCmd.AddCommand(SpecPathCmd)
}
