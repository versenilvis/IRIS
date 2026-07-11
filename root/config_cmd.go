package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage iris configuration",
}

var ConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize default configuration file with comments",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Printf("failed to get config path: %v\n", err)
			return
		}

		if _, statErr := os.Stat(path); statErr == nil {
			fmt.Printf("config file already exists at %s\n", path)
			return
		}

		_ = os.MkdirAll(filepath.Dir(path), 0755)

		defaultContent := `# ~/.config/iris/config.toml
# iris configuration file

[core]
# schema version
# do not edit this field manually
version = 1

# override shell: "bash", "zsh", "fish", keep empty for auto detection
shell = ""

# startup mode: "last", "spec", "history"
# "last" = remember last mode used
mode = "last"

# enable debug logging
debug = false

[ui]
# visual style: "modern" (icons, category pills, shortcut footer) or "classic" (minimalist, centered number, no icons)
style = "modern"

# enable Nerd Fonts icons in overlay menu
nerd-fonts = true

# enable inline ghost text
ghost-text = true

# maximum suggestions to display
max-suggestions = 100

# maximum height of the overlay
max-height = 15

[git]
# hide current branch in checkout/switch list
filter-active-branch = true

# merge remote and local branches with same name
deduplicate-branches = true

[updater]
# check for updates on startup
check-on-startup = true

# update channel: "stable", "nightly"
channel = "stable"

# interval between update checks, e.g. "24h", "6h", "30m"
check-interval = "24h"
`
		err = os.WriteFile(path, []byte(defaultContent), 0644)
		if err != nil {
			fmt.Printf("failed to write config file: %v\n", err)
			return
		}
		fmt.Printf("initialized config file at %s\n", path)
	},
}

var ConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show the resolved configuration",
	Run: func(cmd *cobra.Command, args []string) {
		enc := toml.NewEncoder(cmd.OutOrStdout())
		if err := enc.Encode(config.Get()); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to encode config: %v\n", err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(ConfigInitCmd)
	ConfigCmd.AddCommand(ConfigShowCmd)
	rootCmd.AddCommand(ConfigCmd)
}
