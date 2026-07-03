package root

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/config"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Iris and remove shell integrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling Iris...")
		home, _ := os.UserHomeDir()

		configFiles := []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".config", "fish", "config.fish"),
		}

		for _, file := range configFiles {
			if cleanShellConfig(file) {
				fmt.Printf("✓ Removed integration from %s\n", file)
			}
		}

		// Remove config, state, and cache directories
		if cfgPath, err := config.ConfigPath(); err == nil {
			if cfgDir := filepath.Dir(cfgPath); os.RemoveAll(cfgDir) == nil {
				fmt.Printf("✓ Removed config directory: %s\n", cfgDir)
			}
		}
		if statePath, err := config.StatePath(); err == nil {
			if stateDir := filepath.Dir(statePath); os.RemoveAll(stateDir) == nil {
				fmt.Printf("✓ Removed state directory: %s\n", stateDir)
			}
		}
		if cachePath, err := config.CachePath(); err == nil {
			if os.RemoveAll(cachePath) == nil {
				fmt.Printf("✓ Removed cache directory: %s\n", cachePath)
			}
		}

		binLocations := []string{
			filepath.Join(home, ".local", "bin", "iris"),
			"/usr/local/bin/iris",
		}
		if exe, err := os.Executable(); err == nil && exe != "" {
			binLocations = append(binLocations, exe)
		}

		removedBin := false
		for _, loc := range binLocations {
			if _, err := os.Stat(loc); err == nil {
				if errRemove := os.Remove(loc); errRemove == nil {
					fmt.Printf("✓ Removed binary: %s\n", loc)
					removedBin = true
				} else {
					fmt.Printf("! Could not remove binary at %s (try with sudo): %v\n", loc, errRemove)
				}
			}
		}

		if !removedBin {
			fmt.Println("✓ No leftover binary files found")
		}

		_ = os.Remove("iris.log")

		fmt.Println("\n✓ Iris has been successfully uninstalled")
		if os.Getenv("IRIS_PID") != "" {
			fmt.Println("\n⚠️  You are currently inside an active Iris session.")
			fmt.Println("Iris runs as the parent process of this terminal - do NOT run 'pkill iris'")
			fmt.Println("as it will immediately close this terminal window.")
			fmt.Println("\nTo fully exit, simply close this terminal window and open a new one.")
			fmt.Println("Iris will not start again since the shell config has been cleaned up.")
		} else {
			fmt.Println("Please close and reopen your terminal to complete the uninstall.")
		}
	},
}

func cleanShellConfig(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	var lines []string
	scanner := bufio.NewScanner(f)
	modified := false
	skipNext := false

	for scanner.Scan() {
		line := scanner.Text()
		if skipNext {
			skipNext = false
			modified = true
			continue
		}
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "# iris autocomplete") || strings.Contains(lowerLine, "# iris autostart") {
			modified = true
			skipNext = true
			continue
		}
		if strings.Contains(lowerLine, "iris init") {
			modified = true
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return false
	}

	if !modified {
		return false
	}

	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	output := strings.Join(lines, "\n")
	if len(lines) > 0 {
		output += "\n"
	}

	err = os.WriteFile(filePath, []byte(output), 0644)
	return err == nil
}
