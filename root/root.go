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
	"time"

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
			runWrapper()
		},
	}
	shellFlag string
	isReload  bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&shellFlag, "shell", "s", "", "Shell to use (bash, zsh, fish)")
}

func Execute() {
	if os.Getenv("IRIS_RELOADED") == "true" {
		isReload = true
		// Clear it but keep the knowledge for the current execution
		os.Unsetenv("IRIS_RELOADED")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func detectShell() string {
	// Look up to 5 levels to find the nearest shell
	pid := os.Getppid()
	for i := 0; i < 5 && pid > 1; i++ {
		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err == nil {
			comm := strings.ToLower(strings.TrimSpace(string(data)))
			if strings.Contains(comm, "zsh") { return "zsh" }
			if strings.Contains(comm, "bash") { return "bash" }
			if strings.Contains(comm, "fish") { return "fish" }
		}

		// Move up to parent
		data, err = os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil { break }
		fields := strings.Fields(string(data))
		if len(fields) > 3 {
			ppid, _ := strconv.Atoi(fields[3])
			if ppid == pid || ppid <= 1 { break }
			pid = ppid
		} else { break }
	}

	// Final fallback to system default
	s := os.Getenv("SHELL")
	if strings.Contains(s, "zsh") { return "zsh" }
	return "bash"
}

type procInfo struct {
	pid  int
	ppid int
	comm string
}

// getActiveInnerShell scans the process tree to find the deepest shell running inside the PTY
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
			pid, err1 := strconv.Atoi(fields[0])
			ppid, err2 := strconv.Atoi(fields[1])
			if err1 == nil && err2 == nil {
				comm := strings.ToLower(strings.Join(fields[2:], " "))
				childrenMap[ppid] = append(childrenMap[ppid], procInfo{pid, ppid, comm})
			}
		}
	}

	var findDeepest func(pid int, current string) string
	findDeepest = func(pid int, current string) string {
		shell := current
		for _, child := range childrenMap[pid] {
			childShell := shell
			if strings.Contains(child.comm, "zsh") {
				childShell = "zsh"
			} else if strings.Contains(child.comm, "bash") {
				childShell = "bash"
			} else if strings.Contains(child.comm, "fish") {
				childShell = "fish"
			}

			// recursively find even deeper shells
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
		fmt.Fprintln(os.Stderr, "Error creating pipe:", err)
		return
	}

	var shellName string
	// 1. Priority: Shell from previous reload (inner deepest)
	if active := os.Getenv("IRIS_ACTIVE_SHELL"); active != "" {
		shellName = active
		os.Unsetenv("IRIS_ACTIVE_SHELL")
	// 2. Explicit flag
	} else if shellFlag != "" {
		shellName = shellFlag
	// 3. Dynamic detection (outer parent)
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
		fmt.Fprintln(os.Stderr, "Error starting pty:", err)
		return
	}
	defer ptmx.Close()

	core.ShellPID = c.Process.Pid

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	// signal handling for resize and reload
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
				
				// capture the actual shell running deep inside the PTY
				innerShell := getActiveInnerShell(c.Process.Pid, shellName)
				if innerShell != "" {
					os.Setenv("IRIS_ACTIVE_SHELL", innerShell)
				}

				if oldState != nil {
					_ = term.Restore(int(os.Stdin.Fd()), oldState)
				}
				_ = syscall.Exec(exe, os.Args, os.Environ())
			}
		}
	}()
	sigCh <- syscall.SIGWINCH

	overlay := integration.NewOverlay()

	// PTY -> Stdout
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					// Clean up the terminal and exit when the inner shell dies
					_ = term.Restore(int(os.Stdin.Fd()), oldState)
					os.Exit(0)
				}
				continue
			}
			TermWrite(buf[:n])
		}
	}()

	// Pipe IPC -> Logic -> Render
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
			results := mergeResults(query, "spec") // TODO: implement mode sync

			if len(results) == 0 {
				TermWrite([]byte(overlay.ClearAndDisable()))
				continue
			}

			TermWrite([]byte(overlay.Clear()))
			overlay.UpdateItems(results)
			TermWrite([]byte(overlay.Render()))
		}
	}()

	// Stdin -> PTY
	buf := make([]byte, 1)
	var escSeq []byte
	var naiveBuffer string // fallback naive tracker for bash
	mode := "spec"         // can be "spec" or "history"

	renderOverlay := func() {
		results := mergeResults(naiveBuffer, mode)
		if len(results) == 0 {
			TermWrite([]byte(overlay.ClearAndDisable()))
		} else {
			var buf strings.Builder
			if overlay.Visible {
				buf.WriteString(overlay.Clear())
			}
			overlay.UpdateItems(results)
			buf.WriteString(overlay.Render())
			TermWrite([]byte(buf.String()))
		}
	}

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		if n > 0 {
			b := buf[0]

			// escape sequence parsing
			if b == '\033' {
				escSeq = append(escSeq, b)
				continue
			}
			if len(escSeq) > 0 {
				escSeq = append(escSeq, b)
				if len(escSeq) == 3 {
					seq := string(escSeq)
					escSeq = nil

					if overlay.Visible {
						if seq == "\033[A" { // Up
							overlay.Cursor--
							if overlay.Cursor < 0 {
								overlay.Cursor = 0
							}
							TermWrite([]byte(overlay.Render()))
							continue
						} else if seq == "\033[B" { // Down
							overlay.Cursor++
							if overlay.Cursor >= len(overlay.Items) {
								overlay.Cursor = len(overlay.Items) - 1
							}
							TermWrite([]byte(overlay.Render()))
							continue
						}
					}
					ptmx.Write([]byte(seq))
				}
				continue
			}

			// capture naive typing for bash test
			shouldOverlayDraw := false

			// enter 0x0D, tab 0x09
			if overlay.Visible && (b == '\r' || b == 0x09) {
				selected := overlay.Items[overlay.Cursor].Cmd

				TermWrite([]byte(overlay.ClearAndDisable()))

				// auto add space after tab for next suggestion in spec mode
				if mode == "spec" || b == 0x09 {
					selected += " "
				}

				// sync our naive buffer with the selected result
				naiveBuffer = selected
				mode = "spec" // reset to spec mode after selection

				ptmx.Write(adapter.PrepareSelectSequence(selected))

				if b == '\r' {
					ptmx.Write([]byte{'\r'})
					naiveBuffer = ""
				} else {
					// if it was tab, we want to immediately show next suggestions
					renderOverlay()
				}
				continue
			}

			if b == 0x12 { // Ctrl+R to switch between showing commands mode and history mode
				if mode == "spec" {
					mode = "history"
				} else {
					mode = "spec"
				}
				shouldOverlayDraw = true
			} else if b >= 32 && b <= 126 {
				naiveBuffer += string(b)
				shouldOverlayDraw = true
			} else if b == 127 { // Backspace
				if len(naiveBuffer) > 0 {
					naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
					shouldOverlayDraw = true
				}
			} else if b == '\r' || b == 0x03 {
				naiveBuffer = ""
				mode = "spec" // reset on clear
				TermWrite([]byte(overlay.ClearAndDisable()))
			}

			if b != 0x12 {
				ptmx.Write([]byte{b})
				// small delay to let PTY echo arrive before overlay render
				time.Sleep(3 * time.Millisecond)
			}

			if shouldOverlayDraw {
				renderOverlay()
			}
		}
	}
}

// mergeResults returns suggestions based on the active mode
func mergeResults(query string, mode string) []core.Suggestion {
	if query == "" {
		return nil
	}

	if mode == "history" {
		histResults, _ := integration.SearchHistory(query)
		cmdResults := []core.Suggestion{}
		for _, h := range histResults {
			cmdResults = append(cmdResults, core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history", // please dont remove space here
				Icon: fmt.Sprintf("%d", h.ID),
			})
		}
		if len(cmdResults) > 100 {
			cmdResults = cmdResults[:100]
		}
		return cmdResults
	}

	// Spec mode
	if query == "" {
		return nil
	}

	cmdResults := core.Lookup(query)

	// deduplicate by Cmd
	seen := make(map[string]bool)
	deduped := []core.Suggestion{}
	for _, s := range cmdResults {
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
