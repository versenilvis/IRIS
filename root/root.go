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
			runWrapper()
		},
	}
	shellFlag string
	isReload  bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&shellFlag, "shell", "s", "", "shell to use (bash, zsh, fish)")
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

				// kill the child shell before reloading
				// this prevents multiple shells from fighting over the terminal
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
	mode := "spec"

	renderOverlay := func() {
		// don't render if shell is starting up
		if isReload && naiveBuffer == "" {
			return
		}

		results := mergeResults(naiveBuffer, mode)
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
				if b == 0x12 { // ctrl+r
					intercepted = true
					if mode == "spec" {
						mode = "history"
					} else {
						mode = "spec"
					}
					shouldOverlayDraw = true
				} else if overlay.Visible && (b == '\r' || b == 0x09) {
					intercepted = true
					selected := overlay.Items[overlay.Cursor].Cmd
					os.Stdout.Write([]byte(overlay.ClearAndDisable()))
					if mode == "spec" || b == 0x09 {
						selected += " "
					}
					naiveBuffer = selected
					mode = "spec"
					ptmx.Write(adapter.PrepareSelectSequence(selected))
					if b == '\r' {
						ptmx.Write([]byte{'\r'})
						naiveBuffer = ""
					} else {
						renderOverlay()
					}
					continue
				}

				if !intercepted {
					ptmx.Write([]byte{b})

					// arrow key monitoring with tagged switch
					if b == '\033' && i+2 < n && inputSlice[i+1] == '[' {
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

					switch b {

					case 0x09: // tab
						if !overlay.Visible {
							shouldOverlayDraw = true
						}
					case 127: // backspace
						if len(naiveBuffer) > 0 {
							naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
							shouldOverlayDraw = true
						}
					case '\r', 0x03: // enter, ctrl+c
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
		return nil
	}
	if mode == "history" {
		histResults, _ := integration.SearchHistory(query)
		cmdResults := []core.Suggestion{}
		for _, h := range histResults {
			cmdResults = append(cmdResults, core.Suggestion{
				Cmd:  h.Cmd,
				Desc: " history",
				Icon: fmt.Sprintf("%d", h.ID),
			})
		}
		if len(cmdResults) > 100 {
			cmdResults = cmdResults[:100]
		}
		return cmdResults
	}
	cmdResults := core.Lookup(query)
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
