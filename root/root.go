package root

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"github.com/versenilvis/iris/commands/core"
	_ "github.com/versenilvis/iris/commands/dev"
	_ "github.com/versenilvis/iris/commands/fs"
	_ "github.com/versenilvis/iris/commands/info"
	_ "github.com/versenilvis/iris/commands/runner"
	_ "github.com/versenilvis/iris/commands/search"
	_ "github.com/versenilvis/iris/commands/view"
	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/integration/shell"
	"golang.org/x/term"
)

var (
	rootCmd = &cobra.Command{
		Use:   "iris",
		Short: "IRIS is an awesome cli auto-completion tool",
		Long: `IRIS (a.k.a Intelligent Real-time Input Suggestion) is a shell auto-autocompletion tool.
It works exactly like coding editor suggestion menu drop down.`,
		Run: func(cmd *cobra.Command, args []string) {
			if debugMode {
				f, _ := os.OpenFile("iris.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				debugLogger = f
				core.DebugWriter = f // Link core logger
				fmt.Fprintf(debugLogger, "--- IRIS Debug Started ---\n")
			}
			runWrapper()
		},
	}
	shellFlag   string
	debugMode   bool
	debugLogger *os.File
	isReload    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&shellFlag, "shell", "s", "", "shell to use (bash, zsh, fish)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug logging to iris.log")
}

func debugLog(format string, a ...interface{}) {
	if debugLogger != nil {
		fmt.Fprintf(debugLogger, format+"\n", a...)
	}
}

func Execute() {
	if os.Getenv("IRIS_RELOADED") == "true" {
		isReload = true
		fmt.Printf("\r\033[K\033[35m[IRIS] reloading...\033[0m\n")
		os.Unsetenv("IRIS_RELOADED")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func detectShell() string {
	pid := os.Getppid()
	for i := 0; i < 5 && pid > 1; i++ {
		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err == nil {
			comm := strings.ToLower(strings.TrimSpace(string(data)))
			if strings.Contains(comm, "zsh") {
				return "zsh"
			}
			if strings.Contains(comm, "bash") {
				return "bash"
			}
			if strings.Contains(comm, "fish") {
				return "fish"
			}
		}

		data, err = os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil {
			break
		}
		fields := strings.Fields(string(data))
		if len(fields) > 3 {
			ppid, _ := strconv.Atoi(fields[3])
			if ppid == pid || ppid <= 1 {
				break
			}
			pid = ppid
		} else {
			break
		}
	}

	s := os.Getenv("SHELL")
	if strings.Contains(s, "zsh") {
		return "zsh"
	}
	return "bash"
}

type procInfo struct {
	pid  int
	ppid int
	comm string
}

func getActiveInnerShell(rootPid int, defaultShell string) string {
	cmd := exec.Command("ps", "-e", "-o", "pid,ppid,comm")
	out, err := cmd.Output()
	if err != nil {
		return defaultShell
	}

	lines := strings.Split(string(out), "\n")
	childrenMap := make(map[int][]procInfo)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] != "PID" {
			pid, _ := strconv.Atoi(fields[0])
			ppid, _ := strconv.Atoi(fields[1])
			comm := strings.ToLower(strings.Join(fields[2:], " "))
			childrenMap[ppid] = append(childrenMap[ppid], procInfo{pid, ppid, comm})
		}
	}

	var findDeepest func(pid int, current string) string
	findDeepest = func(pid int, current string) string {
		shell := current
		for _, child := range childrenMap[pid] {
			childShell := shell
			if strings.Contains(child.comm, "zsh") {
				childShell = "zsh"
			}
			if strings.Contains(child.comm, "bash") {
				childShell = "bash"
			}
			if strings.Contains(child.comm, "fish") {
				childShell = "fish"
			}
			if deepest := findDeepest(child.pid, childShell); deepest != "" {
				shell = deepest
			}
		}
		return shell
	}
	return findDeepest(rootPid, defaultShell)
}

func runWrapper() {
	r, w, err := os.Pipe()
	if err != nil {
		return
	}

	var shellName string
	if active := os.Getenv("IRIS_ACTIVE_SHELL"); active != "" {
		shellName = active
		os.Unsetenv("IRIS_ACTIVE_SHELL")
	} else if shellFlag != "" {
		shellName = shellFlag
	} else {
		shellName = detectShell()
	}

	shell.Init(shellName)
	adapter := shell.Current

	c := exec.Command(adapter.GetShellPath())
	c.ExtraFiles = make([]*os.File, 11)
	c.ExtraFiles[10] = w
	c.Env = adapter.GetEnv(10, os.Getpid())

	ptmx, err := pty.Start(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[IRIS] failed to start PTY: %v\n", err)
		return
	}
	defer ptmx.Close()

	_ = pty.InheritSize(os.Stdin, ptmx)
	core.ShellPID = c.Process.Pid

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGWINCH, syscall.SIGUSR1)
	go func() {
		for s := range sigCh {
			switch s {
			case syscall.SIGWINCH:
				_ = pty.InheritSize(os.Stdin, ptmx)
			case syscall.SIGUSR1:
				exe, _ := os.Executable()
				os.Setenv("IRIS_RELOADED", "true")

				// capture the current shell state
				innerShell := getActiveInnerShell(c.Process.Pid, shellName)
				if innerShell != "" {
					os.Setenv("IRIS_ACTIVE_SHELL", innerShell)
				}

				// kill the child process before reloading
				if c.Process != nil {
					_ = syscall.Kill(c.Process.Pid, syscall.SIGKILL)
					ptmx.Close()
				}

				if oldState != nil {
					_ = term.Restore(int(os.Stdin.Fd()), oldState)
				}
				_ = syscall.Exec(exe, os.Args, os.Environ())
			}
		}
	}()

	overlay := integration.NewOverlay()

	// pty -> stdout
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					_ = term.Restore(int(os.Stdin.Fd()), oldState)
					os.Exit(0)
				}
				continue
			}
			os.Stdout.Write(buf[:n])
		}
	}()

	// ipc pipe
	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.IndexByte(data, '\x00'); i >= 0 {
				return i + 1, data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		})

		for scanner.Scan() {
			query := scanner.Text()
			results := mergeResults(query, "spec")
			if len(results) == 0 {
				os.Stdout.Write([]byte(overlay.ClearAndDisable()))
				continue
			}
			os.Stdout.Write([]byte(overlay.Clear()))
			overlay.UpdateItems(results)
			os.Stdout.Write([]byte(overlay.Render()))
		}
	}()

	var naiveBuffer string
	suggestionsEnabled := true
	mode := "spec"

	renderOverlay := func() {
		if !suggestionsEnabled {
			return
		}
		
		// If buffer is empty (e.g. backspace to 0 or after reload), clear overlay immediately
		if naiveBuffer == "" {
			os.Stdout.Write([]byte(overlay.ClearAndDisable()))
			return
		}

		debugLog("[Render] query: '%s', mode: %s", naiveBuffer, mode)
		results := mergeResults(naiveBuffer, mode)
		debugLog("[Render] results found: %d", len(results))

		if len(results) == 0 {
			os.Stdout.Write([]byte(overlay.ClearAndDisable()))
		} else {
			var buf strings.Builder
			if overlay.Visible {
				buf.WriteString(overlay.Clear())
			}
			overlay.UpdateItems(results)
			buf.WriteString(overlay.Render())
			os.Stdout.Write([]byte(buf.String()))
		}
	}

	// Trigger initial render to clear any potential artifacts from previous session
	renderOverlay()

	for {
		inputSlice := make([]byte, 128)
		n, err := os.Stdin.Read(inputSlice)
		if err != nil {
			break
		}

		if n > 0 {
			shouldOverlayDraw := false
			for i := 0; i < n; i++ {
				b := inputSlice[i]
				intercepted := false

				// Detect Escape sequence (e.g. arrows, Shift+Tab) early
				if b == '\033' {
					// Shift+Tab: \033[Z
					if i+2 < n && inputSlice[i+1] == '[' && inputSlice[i+2] == 'Z' {
						intercepted = true
						suggestionsEnabled = !suggestionsEnabled
						if !suggestionsEnabled {
							os.Stdout.Write([]byte(overlay.ClearAndDisable()))
						} else {
							shouldOverlayDraw = true
						}
						i += 2
						continue
					}

					ptmx.Write([]byte{b})
					// arrow key monitoring inside escape
					if i+2 < n && inputSlice[i+1] == '[' {
						if overlay.Visible {
							switch inputSlice[i+2] {
							case 'A': // up
								overlay.Cursor--
								if overlay.Cursor < 0 {
									overlay.Cursor = 0
								}
								os.Stdout.Write([]byte(overlay.Render()))
							case 'B': // down
								overlay.Cursor++
								if overlay.Cursor >= len(overlay.Items) {
									overlay.Cursor = len(overlay.Items) - 1
								}
								os.Stdout.Write([]byte(overlay.Render()))
							}
						}
					}
					// skip sequence
					for j := i + 1; j < n; j++ {
						char := inputSlice[j]
						ptmx.Write([]byte{char})
						i = j
						if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '~' {
							break
						}
					}
					continue
				}

				// Iris shortcuts
				if b == 0x12 { // ctrl+r
					intercepted = true
					if mode == "spec" {
						mode = "history"
					} else {
						mode = "spec"
					}
					shouldOverlayDraw = true
				} else if b == 0x1b && n == 1 { // Standalone ESC (Dismiss CURRENT menu only)
					intercepted = true
					os.Stdout.Write([]byte(overlay.ClearAndDisable()))
					shouldOverlayDraw = false // Just hide it for this render
					continue
				} else if overlay.Visible && (b == 0x0d || b == 0x0a) { // Enter while menu open
					intercepted = true
					os.Stdout.Write([]byte(overlay.ClearAndDisable()))
					// Send Enter to PTY to execute current buffer as-is
					ptmx.Write([]byte{0x0d})
					naiveBuffer = ""
					shouldOverlayDraw = false
					continue
				} else if b == 0x09 { // Tab (Universal Intercept - Select/Commit)
					intercepted = true
					if !overlay.Visible {
						debugLog("[Input] Tab pressed: opening menu")
						shouldOverlayDraw = true
					} else {
						// SELECT and COMMIT
						selected := overlay.Items[overlay.Cursor].Cmd
						debugLog("[Input] Tab pressed: committing '%s'", selected)

						os.Stdout.Write([]byte(overlay.ClearAndDisable()))

						// Auto-add space for spec mode
						if mode == "spec" {
							selected = strings.TrimSpace(selected) + " "
						}

						naiveBuffer = selected
						mode = "spec"
						// Sync PTY with the new buffer (Ctrl+A then Ctrl+K to clear rest)
						ptmx.Write([]byte{0x01, 0x0b}) // bash/zsh style home and clear
						ptmx.Write([]byte(selected))

						shouldOverlayDraw = false
					}
					continue
				}

				// normal typing
				if !intercepted {
					ptmx.Write([]byte{b})
					switch b {
					case 127, 0x08: // backspace
						if len(naiveBuffer) > 0 {
							naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
							shouldOverlayDraw = true
						}
					case 0x17: // Ctrl+W: delete one word
						trimBuf := strings.TrimRight(naiveBuffer, " ")
						lastSpace := strings.LastIndex(trimBuf, " ")
						if lastSpace >= 0 {
							naiveBuffer = trimBuf[:lastSpace+1]
						} else {
							naiveBuffer = ""
						}
						shouldOverlayDraw = true
					case '\r', 0x03, 0x15, 0x0C: // Enter, Ctrl+C, Ctrl+U, Ctrl+L
						naiveBuffer = ""
						mode = "spec"
						os.Stdout.Write([]byte(overlay.ClearAndDisable()))
					default:
						if b >= 32 && b <= 126 {
							naiveBuffer += string(b)
							shouldOverlayDraw = true
						}
					}
				}
			}
			if shouldOverlayDraw {
				renderOverlay()
			}
		}
	}
}

func mergeResults(query string, mode string) []core.Suggestion {
	if query == "" {
		debugLog("[Merge] Query empty, returning nil")
		return nil
	}

	normalizedQuery := strings.TrimSpace(query)
	seen := make(map[string]bool)
	deduped := []core.Suggestion{}

	if mode == "history" {
		histResults, _ := integration.SearchHistory(query)
		for _, h := range histResults {
			normalizedCmd := strings.TrimSpace(h.Cmd)
			if seen[normalizedCmd] {
				continue
			}
			seen[normalizedCmd] = true
			deduped = append(deduped, core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history",
				Icon: fmt.Sprintf("%d", h.ID),
			})
			if len(deduped) >= 10 {
				break
			}
		}
		debugLog("[Merge] History mode found %d items", len(deduped))
		return deduped
	}

	debugLog("[Merge] Calling Lookup for '%s'", query)
	cmdResults := core.Lookup(query)
	debugLog("[Merge] Lookup returned %d raw items", len(cmdResults))

	for _, s := range cmdResults {
		normalizedCmd := strings.TrimSpace(s.Cmd)
		// CRITICAL: Filter out exact matches to prevent infinite tab loops
		if normalizedCmd == normalizedQuery {
			debugLog("[Merge] Filtered EXACT MATCH: '%s'", normalizedCmd)
			continue
		}

		if !seen[s.Cmd] {
			seen[s.Cmd] = true
			deduped = append(deduped, s)
		}
	}
	if len(deduped) > 10 {
		deduped = deduped[:10]
	}
	return deduped
}
