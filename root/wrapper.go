package root

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/integration/shell"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type State struct {
	Mode string `json:"mode"`
}

func getStateFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, ".iris")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "state.json")
}

func loadMode() string {
	file := getStateFile()
	if file != "" {
		data, err := os.ReadFile(file)
		if err == nil {
			var state State
			if err := json.Unmarshal(data, &state); err == nil {
				if state.Mode == "history" || state.Mode == "spec" {
					return state.Mode
				}
			}
		}
	}
	return "spec"
}

func saveMode(mode string) {
	file := getStateFile()
	if file != "" {
		state := State{Mode: mode}
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			os.WriteFile(file, data, 0644)
		}
	}
}

// runWrapper sets up the pty environment, launches the shell,
// and manages the main input loop to provide real-time suggestions
// it handles raw terminal mode to intercept keystrokes and
// coordinates between the shell process and the suggestion overlay
func runWrapper() {
	r, w, err := os.Pipe() // pipe for ipc communication from shell to iris
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
	// pass write end of pipe to shell as fd 10 (I chose 10 just because it won't conflict with other file descriptors)
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

	// put terminal in raw mode to intercept every keystroke
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
				_ = pty.InheritSize(os.Stdin, ptmx) // handle terminal window resize
			// this is the core feature of reloading
			// it helps IRIS reload itself that you dont need to restart the shell mannually
			// SIGUSR1 is the signal to active reload when you type "just reload"
			case syscall.SIGUSR1:
				// trigger iris reload by executing itself again
				exe, _ := os.Executable()
				// this marks for the next iris process that it've just reloaded
				os.Setenv("IRIS_RELOADED", "true")

				innerShell := getActiveInnerShell(c.Process.Pid, shellName)
				if innerShell != "" {
					// to detect which is last shell (bash, zsh, fish)
					os.Setenv("IRIS_ACTIVE_SHELL", innerShell)
				}

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

	// bridge pty output to actual stdout
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

	var disableGhostText atomic.Bool
	var userNavigated bool
	var renderOverlay func()

	isExecuting := func() bool {
		pgrp, err := unix.IoctlGetInt(int(ptmx.Fd()), unix.TIOCGPGRP)
		if err != nil {
			return false
		}
		shellPGID, err := unix.Getpgid(core.ShellPID)
		if err != nil {
			return pgrp != core.ShellPID
		}
		return pgrp != shellPGID
	}

	// listen for suggestion requests from shell scripts via the ipc pipe
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

			if query == "IRIS_CMD_STOP" {
				continue
			}

			if isExecuting() {
				continue
			}

			results := mergeResults(query, "spec")
			if len(results) == 0 {
				os.Stdout.Write([]byte(overlay.ClearAndDisable()))
				continue
			}
			os.Stdout.Write([]byte(overlay.Clear()))
			overlay.UpdateItems(results)
			var rBuf strings.Builder
			if !disableGhostText.Load() {
				rBuf.WriteString(overlay.RenderGhostText(query, false))
			}
			rBuf.WriteString(overlay.Render())
			os.Stdout.Write([]byte(rBuf.String()))
		}
	}()

	var naiveBuffer string
	suggestionsEnabled := true
	mode := loadMode()

	var renderTimer *time.Timer
	var renderMu sync.Mutex

	// renderOverlay decides whether to draw the suggestion menu based on current state
	renderOverlay = func() {
		renderMu.Lock()
		defer renderMu.Unlock()

		if renderTimer != nil {
			renderTimer.Stop()
		}

		if !suggestionsEnabled || isExecuting() {
			return
		}

		bufCopy := naiveBuffer
		modeCopy := mode
		navCopy := userNavigated

		if bufCopy == "" && !navCopy {
			os.Stdout.Write([]byte(overlay.ClearAndDisable()))
			return
		}

		// Debounce for 15ms to allow PTY to process the keystroke and update the terminal cursor.
		// This completely prevents the asynchronous ghost text race condition where the PTY echo
		// overwrites the first letter of our ghost text!
		renderTimer = time.AfterFunc(15*time.Millisecond, func() {
			if isExecuting() {
				return
			}

			var b strings.Builder
			if !navCopy {
				debugLog("[Render] query: '%s', mode: %s", bufCopy, modeCopy)
				results := mergeResults(bufCopy, modeCopy)
				debugLog("[Render] results found: %d", len(results))

				if len(results) == 0 {
					b.WriteString(overlay.ClearAndDisable())
					os.Stdout.Write([]byte(b.String()))
					return
				}

				if overlay.Visible {
					b.WriteString(overlay.Clear())
				}
				overlay.UpdateItems(results)
			} else {
				if overlay.Visible {
					b.WriteString(overlay.Clear())
				}
			}

			if !disableGhostText.Load() {
				b.WriteString(overlay.RenderGhostText(bufCopy, navCopy))
			}
			b.WriteString(overlay.Render())
			os.Stdout.Write([]byte(b.String()))
		})
	}

	renderOverlay()

	renderNow := func() {
		renderMu.Lock()
		if renderTimer != nil {
			renderTimer.Stop()
		}
		renderMu.Unlock()

		var b strings.Builder
		if overlay.Visible {
			if !disableGhostText.Load() {
				b.WriteString(overlay.RenderGhostText(naiveBuffer, userNavigated))
			}
			b.WriteString(overlay.Render())
		}
		os.Stdout.Write([]byte(b.String()))
	}

	// reads from stdin and decides what to forward or intercept
	// for most cases, I just handle the already have terminal shortcuts
	// for some shortcuts like tab, enter, shift tab, ctrl r,
	// they have a little bit different behavior to match our tool
	for {
		inputSlice := make([]byte, 128)
		n, err := os.Stdin.Read(inputSlice)
		if err != nil {
			break
		}

		if n > 0 {
			if isExecuting() {
				ptmx.Write(inputSlice[:n])
				continue
			}

			shouldOverlayDraw := false
			for i := 0; i < n; i++ {
				b := inputSlice[i]
				intercepted := false

				if b == '\033' {
					// handle escape sequences like arrow keys and functional shortcuts
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						// shift tab: hide/unhide menu dropdown
						if inputSlice[i+1] == '[' && inputSlice[i+2] == 'Z' {
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

						if overlay.Visible && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') {
							intercepted = true
							userNavigated = true

							os.Stdout.Write([]byte(overlay.Clear())) // clear old menu

							if inputSlice[i+2] == 'A' { // up arrow
								overlay.Cursor--
								if overlay.Cursor < 0 {
									overlay.Cursor = 0
								}
							} else { // down arrow
								overlay.Cursor++
								if overlay.Cursor >= len(overlay.Items) {
									overlay.Cursor = len(overlay.Items) - 1
								}
							}

							selected := overlay.Items[overlay.Cursor].Cmd
							ptmx.Write([]byte{0x15}) // ctrl+u to clear line
							ptmx.Write([]byte(selected))
							naiveBuffer = selected

							renderNow()
							i += 2
							continue
						} else if !overlay.Visible && naiveBuffer == "" && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') { // up/down arrow on empty prompt
							intercepted = true
							mode = "history"
							saveMode(mode)

							results := mergeResults("", "history")
							if len(results) > 0 {
								limit := 100
								if len(results) < limit {
									limit = len(results)
								}
								var historyList []core.Suggestion

								if inputSlice[i+2] == 'A' {
									// Up arrow: Reverse the list so newest is at the bottom
									for j := limit - 1; j >= 0; j-- {
										historyList = append(historyList, results[j])
									}
								} else {
									// Down arrow: Normal order, newest is at the top
									for j := 0; j < limit; j++ {
										historyList = append(historyList, results[j])
									}
								}

								overlay.UpdateItems(historyList)

								if inputSlice[i+2] == 'A' {
									overlay.Cursor = len(historyList) - 1 // Up arrow: Start at the bottom
								} else {
									overlay.Cursor = 0 // Down arrow: Start at the top
								}

								selected := overlay.Items[overlay.Cursor].Cmd
								ptmx.Write([]byte{0x15}) // ctrl+u to clear line
								ptmx.Write([]byte(selected))
								naiveBuffer = selected

								userNavigated = true
								renderNow()
							}
							i += 2
							continue
						} else if overlay.Visible && !disableGhostText.Load() && inputSlice[i+2] == 'C' { // right arrow
							topCmd := overlay.Items[0].Cmd
							if strings.HasPrefix(strings.ToLower(topCmd), strings.ToLower(naiveBuffer)) {
								ghostText := topCmd[len(naiveBuffer):]
								if len(ghostText) > 0 {
									intercepted = true
									naiveBuffer += ghostText
									ptmx.Write([]byte(ghostText))
									shouldOverlayDraw = true
									i += 2
									continue
								}
							}
						}
					}

					// forward escape sequence to pty if not intercepted
					if !intercepted {
						os.Stdout.Write([]byte(overlay.ClearAndDisable()))
						disableGhostText.Store(true)
						naiveBuffer = ""

						ptmx.Write([]byte{b})
						// skip remaining bytes of the escape sequence to avoid misinterpretation
						for j := i + 1; j < n; j++ {
							char := inputSlice[j]
							ptmx.Write([]byte{char})
							i = j
							if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '~' {
								break
							}
						}
					}
					continue
				}

				if b == 0x12 { // ctrl+r: toggle between command specs and command history
					intercepted = true
					if mode == "spec" {
						mode = "history"
					} else {
						mode = "spec"
					}
					saveMode(mode)
					shouldOverlayDraw = true
					// enter: enter behavior is a bit different from tab suggestions in code editor
					// I want it to execute the command anyway and ignore the suggestions
					// it means only tab to select suggestions, and enter to execute
					// enter is not used to select suggestions
				} else if overlay.Visible && (b == 0x0d || b == 0x0a) {
					intercepted = true
					os.Stdout.Write([]byte(overlay.ClearAndDisable()))

					ptmx.Write([]byte{0x0d})
					naiveBuffer = ""
					disableGhostText.Store(false)
					shouldOverlayDraw = false
					userNavigated = false
					continue
				} else if b == 0x09 { // tab: select suggestions
					intercepted = true
					if !overlay.Visible {
						shouldOverlayDraw = true
					} else {
						selected := overlay.Items[overlay.Cursor].Cmd
						os.Stdout.Write([]byte(overlay.ClearAndDisable()))

						if mode == "spec" {
							selected = strings.TrimSpace(selected) + " "
						}

						naiveBuffer = selected

						ptmx.Write([]byte{0x15}) // ctrl+u to clear line
						ptmx.Write([]byte(selected))

						overlay.Cursor = 0 // this prevents when you tab, it switchs between suggestions non-stop

						shouldOverlayDraw = true // <- rerender after tab to choose, if you set to false,
						// when you press tab continuely, it will print all folder from menu suggestions
						// and make the cursor jump to next line
						userNavigated = false
					}
					continue
				}

				if !intercepted {
					ptmx.Write([]byte{b})
					// we handle line editing keys manually to keep naiveBuffer in sync
					// since terminal is in raw mode, we must update our state for every change
					switch b {
					case 127, 0x08: // backspace: remove last character from buffer
						if len(naiveBuffer) > 0 {
							naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
							shouldOverlayDraw = true
							userNavigated = false
						}
					case 0x17: // ctrl+w: delete the last word in the buffer
						trimBuf := strings.TrimRight(naiveBuffer, " ")
						lastSpace := strings.LastIndex(trimBuf, " ")
						if lastSpace >= 0 {
							naiveBuffer = trimBuf[:lastSpace+1]
						} else {
							naiveBuffer = ""
						}
						shouldOverlayDraw = true
						userNavigated = false
					case '\r', '\n', 0x03, 0x15, 0x0C: // enter, ctrl+c, ctrl+u, ctrl+l: clear buffer on line reset
						naiveBuffer = ""
						disableGhostText.Store(false)
						os.Stdout.Write([]byte(overlay.ClearAndDisable()))
						userNavigated = false
					default:
						// track normal printable characters in the buffer for matching
						if b >= 32 && b <= 126 {
							// if user presses space, check if the current word is an alias
							if b == ' ' && naiveBuffer != "" && !strings.Contains(naiveBuffer, " ") {
								if target, ok := core.GetAlias(naiveBuffer); ok {
									// clear the current alias and replace it with the full command
									ptmx.Write([]byte{0x15}) // ctrl+u to clear the current input line
									ptmx.Write([]byte(target + " "))
									naiveBuffer = target + " "
									shouldOverlayDraw = true
									continue
								}
							}
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
