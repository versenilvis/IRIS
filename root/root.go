package root

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/commands/core"
	_ "github.com/versenilvis/iris/commands"
	"github.com/versenilvis/iris/config"
	"golang.org/x/term"
)

var (
	rootCmd = &cobra.Command{
		Use:   "iris",
		Short: "IRIS is an awesome cli auto-completion tool",
		Long: `IRIS (a.k.a Intelligent Real-time Input Suggestion) is a shell auto-autocompletion tool.
It works exactly like coding editor suggestion menu drop down.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if r := recover(); r != nil {
					WriteCrashLog(r)
					restoreTerminal()
					printCrashNotice()
					startRescueShell()
					os.Exit(2)
				}
			}()
			if pidStr := os.Getenv("IRIS_PID"); pidStr != "" {
				if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
					_ = syscall.Kill(pid, syscall.SIGUSR1)
					fmt.Println("\r\033[K\033[36m[IRIS] Sent reload signal to parent session.\033[0m")
					return
				}
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

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if shellFlag != "" {
			config.Get().Core.Shell = shellFlag
		}
		if debugMode {
			config.Get().Core.Debug = true
		}
		if config.Get().Core.Debug {
			logDir, err := config.CachePath()
			if err == nil {
				_ = os.MkdirAll(logDir, 0755)
				f, _ := os.OpenFile(filepath.Join(logDir, "iris.log"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				debugLogger = f
				core.DebugWriter = f
				_, _ = fmt.Fprintf(debugLogger, "--- IRIS DEBUG LOG ---\n")
			}
		}
	}
}

func debugLog(format string, a ...any) {
	if debugLogger != nil {
		_, _ = fmt.Fprintf(debugLogger, format+"\n", a...)
	}
}

// runWatchdog spawns the watchdog parent process
func runWatchdog() {
	exe, err := os.Executable()
	if err != nil {
		runOriginal()
		return
	}

	// save original terminal settings in parent process
	watchdogOldState, errState := term.MakeRaw(int(os.Stdin.Fd()))
	if errState == nil {
		_ = term.Restore(int(os.Stdin.Fd()), watchdogOldState)
	}

	r, w, err := os.Pipe()
	if err != nil {
		runOriginal()
		return
	}

	cmd := exec.CommandContext(context.Background(), exe, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "IRIS_IS_CHILD=true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = w

	err = cmd.Start()
	if err != nil {
		runOriginal()
		return
	}

	_ = w.Close()

	// copy child stderr to both our buffer and the real stderr, filtering out panics
	var stderrBuf bytes.Buffer
	origStderr := os.Stderr
	tempBuf := make([]byte, 1024)
	suppress := false
	for {
		n, errRead := r.Read(tempBuf)
		if n > 0 {
			_, _ = stderrBuf.Write(tempBuf[:n])
			if stderrBuf.Len() > 64*1024 {
				// discard oldest bytes to avoid memory leak
				over := stderrBuf.Len() - 64*1024
				_ = stderrBuf.Next(over)
			}
			if !suppress {
				currentContent := stderrBuf.Bytes()
				searchStart := 0
				if len(currentContent) > n+12 {
					searchStart = len(currentContent) - (n + 12)
				}
				searchSlice := currentContent[searchStart:]
				idxPanic := bytes.Index(searchSlice, []byte("panic:"))
				idxFatal := bytes.Index(searchSlice, []byte("fatal error:"))
				triggerIdx := -1
				if idxPanic != -1 {
					triggerIdx = searchStart + idxPanic
				} else if idxFatal != -1 {
					triggerIdx = searchStart + idxFatal
				}

				if triggerIdx != -1 {
					suppress = true
					printedLen := len(currentContent) - n
					if triggerIdx > printedLen {
						_, _ = origStderr.Write(currentContent[printedLen:triggerIdx])
					}
				} else {
					_, _ = origStderr.Write(tempBuf[:n])
				}
			}
		}
		if errRead != nil {
			break
		}
	}

	// check if child exited abnormally or crashed
	errWait := cmd.Wait()
	if errWait != nil {
		content := stderrBuf.Bytes()
		if bytes.Contains(content, []byte("panic:")) || bytes.Contains(content, []byte("fatal error:")) {
			WriteCrashLog(string(content))
			// restore terminal state if watchdog saved it
			if watchdogOldState != nil {
				_ = term.Restore(int(os.Stdin.Fd()), watchdogOldState)
			}
			printCrashNotice()
			startRescueShell()
			os.Exit(2)
		}

		var exitErr *exec.ExitError
		if errors.As(errWait, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

// runOriginal runs the normal command execution
func runOriginal() {
	if os.Getenv("IRIS_RELOADED") == "true" {
		fmt.Printf("\r\033[K\033[35m[IRIS] reloading...\033[0m\n")
		_ = os.Unsetenv("IRIS_RELOADED")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Execute() {
	_ = config.MigrateFromLegacyJSON()
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[IRIS] config error: %v\n", err)
	}
	config.Init(cfg)

	if os.Getenv("IRIS_IS_CHILD") != "true" {
		runWatchdog()
		return
	}

	runOriginal()
}
