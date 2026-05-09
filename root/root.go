package root

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/commands/core"
	_ "github.com/versenilvis/iris/commands/dev"
	_ "github.com/versenilvis/iris/commands/fs"
	_ "github.com/versenilvis/iris/commands/info"
	_ "github.com/versenilvis/iris/commands/runner"
	_ "github.com/versenilvis/iris/commands/search"
	_ "github.com/versenilvis/iris/commands/view"
)

var (
	rootCmd = &cobra.Command{
		Use:   "iris",
		Short: "IRIS is an awesome cli auto-completion tool",
		Long: `IRIS (a.k.a Intelligent Real-time Input Suggestion) is a shell auto-autocompletion tool.
It works exactly like coding editor suggestion menu drop down.`,
		Run: func(cmd *cobra.Command, args []string) {
			if pidStr := os.Getenv("IRIS_PID"); pidStr != "" {
				if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
					_ = syscall.Kill(pid, syscall.SIGUSR1)
					fmt.Println("\r\033[K\033[36m[IRIS] Sent reload signal to parent session.\033[0m")
					return
				}
			}
			if debugMode {
				f, _ := os.OpenFile("iris.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				debugLogger = f
				core.DebugWriter = f
				_, _ = fmt.Fprintf(debugLogger, "--- IRIS DEBUG LOG ---\n")
			}
			runWrapper()
		},
	}
	shellFlag   string
	debugMode   bool
	debugLogger *os.File
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&shellFlag, "shell", "s", "", "shell to use (bash, zsh, fish)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug logging to iris.log")
}

func debugLog(format string, a ...interface{}) {
	if debugLogger != nil {
		_, _ = fmt.Fprintf(debugLogger, format+"\n", a...)
	}
}

func Execute() {
	if os.Getenv("IRIS_RELOADED") == "true" {
		fmt.Printf("\r\033[K\033[35m[IRIS] reloading...\033[0m\n")
		_ = os.Unsetenv("IRIS_RELOADED")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
