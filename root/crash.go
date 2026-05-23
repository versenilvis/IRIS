package root

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// startRescueShell starts a fallback shell if the application crashes to keep the terminal open
func startRescueShell() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	_ = syscall.Exec(shell, []string{shell}, os.Environ())
}

var (
	lastCrashFile string
	lastCrashMu   sync.Mutex
)

// writeCrashLog writes the crash info and stack trace to a new log file
func WriteCrashLog(err any) {
	home, errDir := os.UserHomeDir()
	if errDir != nil {
		return
	}
	dir := filepath.Join(home, ".iris", "crashes")
	_ = os.MkdirAll(dir, 0755)
	logFile := filepath.Join(dir, fmt.Sprintf("crash_%s.log", time.Now().Format("20060102_150405")))

	lastCrashMu.Lock()
	lastCrashFile = logFile
	lastCrashMu.Unlock()

	f, errOpen := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if errOpen != nil {
		return
	}
	defer f.Close()

	_, _ = fmt.Fprintf(f, "=== IRIS CRASH %s ===\n", time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprintf(f, "version: %s\nos: %s/%s\n\n", Version, runtime.GOOS, runtime.GOARCH)
	_, _ = fmt.Fprintf(f, "panic: %v\n\n", err)

	buf := make([]byte, 64*1024)
	n := runtime.Stack(buf, true)
	_, _ = f.Write(buf[:n])
	_, _ = fmt.Fprintln(f)
}

// getLatestCrashLog returns the path to the newest crash log file
func getLatestCrashLog() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, ".iris", "crashes")
	files, err := os.ReadDir(dir)
	if err != nil || len(files) == 0 {
		oldLog := filepath.Join(home, ".iris", "crash.log")
		if _, err := os.Stat(oldLog); err == nil {
			return oldLog
		}
		return ""
	}

	var latest string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasPrefix(name, "crash_") && strings.HasSuffix(name, ".log") {
			if name > latest {
				latest = name
			}
		}
	}
	if latest == "" {
		oldLog := filepath.Join(home, ".iris", "crash.log")
		if _, err := os.Stat(oldLog); err == nil {
			return oldLog
		}
		return ""
	}
	return filepath.Join(dir, latest)
}

// printCrashNotice prints the crash notice with absolute path to the log file
func printCrashNotice() {
	lastCrashMu.Lock()
	logFile := lastCrashFile
	lastCrashMu.Unlock()
	if logFile == "" {
		logFile = getLatestCrashLog()
	}
	_, _ = fmt.Fprintf(os.Stderr, "\n\033[31m[IRIS] crashed, report saved to %s\033[0m\n", logFile)
}

var (
	// crashCmd is the cobra command to manage crash logs
	CrashCmd = &cobra.Command{
		Use:   "crash-log",
		Short: "manage iris crash logs",
		Run: func(cmd *cobra.Command, args []string) {
			home, err := os.UserHomeDir()
			if err != nil {
				cmd.Printf("failed to get home directory: %v\n", err)
				return
			}

			if ClearLog {
				_ = os.RemoveAll(filepath.Join(home, ".iris", "crashes"))
				_ = os.Remove(filepath.Join(home, ".iris", "crash.log"))
				cmd.Println("crash log cleared")
				return
			}

			logFile := getLatestCrashLog()
			if logFile == "" {
				cmd.Println("no crash log found")
				return
			}
			cmd.Println(logFile)
		},
	}
	// clearLog is the flag to clear the crash log
	ClearLog bool
)

func init() {
	CrashCmd.Flags().BoolVar(&ClearLog, "clear", false, "clear the crash log")
	rootCmd.AddCommand(CrashCmd)
}
