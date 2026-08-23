package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/internal/config"
)

var ThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "manage iris theme",
}

var ThemeInitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize default theme file with comments",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := config.ThemePath()
		if err != nil {
			fmt.Printf("failed to get theme path: %v\n", err)
			return
		}

		if _, statErr := os.Stat(path); statErr == nil {
			fmt.Printf("theme file already exists at %s\n", path)
			return
		}

		_ = os.MkdirAll(filepath.Dir(path), 0755)

		defaultContent := `# ~/.config/iris/theme.toml
# iris theme file

# overlay border color
border = "#a277ff"

# accent color for selected marker and icon
accent = "#61ffca"

# muted color for unselected icons
muted = "#6d6a7f"

# default text color
text = "#edecee"

# text color when selected
text_sel = "#ffffff"

# footer shortcut key color
key = "#a277ff"

# matched text color
match = "#61ffca"

# description text color
desc = "#9692a8"

# description text color when selected
desc_sel = "#edecee"

# selected row background
sel_bg = "#3d375e"

# selected tag pill text color
sel_text = "#110f18"

# scroll counter color
scroll_info = "#a277ff"

# inline ghost text color
ghost_text = "#4B4A4C"

# "system" tag background color
sys = "#1e1d28"

# "system" tag selected background color
sys_sel = "#a277ff"

# "history" tag background color
hist = "#1a2d36"

# "history" tag selected background color
hist_sel = "#61ffca"

# "alias" tag background color
alias = "#2a2342"

# "alias" tag selected background color (and alias word color)
alias_sel = "#a277ff"
`
		err = os.WriteFile(path, []byte(defaultContent), 0644)
		if err != nil {
			fmt.Printf("failed to write theme file: %v\n", err)
			return
		}
		fmt.Printf("initialized theme file at %s\n", path)
	},
}

func init() {
	ThemeCmd.AddCommand(ThemeInitCmd)
	rootCmd.AddCommand(ThemeCmd)
}