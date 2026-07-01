package root

import (
	"bufio"
	"bytes"
	"context"
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
	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/integration/shell"
	"github.com/versenilvis/iris/logger"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func loadMode() string {
	mode := config.Get().Core.Mode
	if mode == "last" {
		state := config.LoadState()
		if state.LastMode == "history" || state.LastMode == "spec" {
			return state.LastMode
		}
		return "spec"
	}
	if mode == "history" || mode == "spec" {
		return mode
	}
	return "spec"
}

func saveMode(mode string) {
	state := config.LoadState()
	state.LastMode = mode
	_ = config.SaveState(state)
}

var (
	oldState     *term.State
	oldStateMu   sync.Mutex
	activeMode   string
	activeModeMu sync.RWMutex
	stdoutMu     sync.Mutex
)

func writeStdout(data []byte) {
	if len(data) == 0 {
		return
	}
	stdoutMu.Lock()
	defer stdoutMu.Unlock()
	_, _ = os.Stdout.Write(data)
}

// restoreTerminal restores the terminal state if needed
func restoreTerminal() {
	oldStateMu.Lock()
	defer oldStateMu.Unlock()
	if oldState != nil {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		oldState = nil
	}
}

// runWrapper sets up the pty environment, launches the shell,
// and manages the main input loop to provide real-time suggestions
// it handles raw terminal mode to intercept keystrokes and
// coordinates between the shell process and the suggestion overlay
func runWrapper() {
	var naiveBuffer string
	cursorOffset := 0
	var bufferMu sync.Mutex

	r, w, err := os.Pipe() // pipe for ipc communication from shell to iris
	if err != nil {
		return
	}

	var shellName string
	if active := os.Getenv("IRIS_ACTIVE_SHELL"); active != "" {
		shellName = active
		_ = os.Unsetenv("IRIS_ACTIVE_SHELL")
	} else if shellFlag != "" {
		shellName = shellFlag
	} else {
		shellName = detectShell()
	}

	shell.Init(shellName)
	adapter := shell.Current

	ctx := context.Background()
	c := exec.CommandContext(ctx, adapter.GetShellPath())
	c.ExtraFiles = make([]*os.File, 11)
	// pass write end of pipe to shell as fd 13 (since index 10 maps to 13)
	c.ExtraFiles[10] = w
	c.Env = adapter.GetEnv(13, os.Getpid())

	ptmx, err := pty.Start(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[IRIS] failed to start PTY: %v\n", err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	_ = pty.InheritSize(os.Stdin, ptmx)
	core.ShellPID = c.Process.Pid

	logger.Infof("PTY child shell started: shell=%s, path=%s, pid=%d", shellName, adapter.GetShellPath(), c.Process.Pid)

	// put terminal in raw mode to intercept every keystroke
	var errMakeRaw error
	oldState, errMakeRaw = term.MakeRaw(int(os.Stdin.Fd()))
	if errMakeRaw != nil {
		logger.Errorf("Failed to set terminal raw mode: %v", errMakeRaw)
		panic(errMakeRaw)
	}
	logger.Debugf("Terminal set to raw mode successfully")
	defer restoreTerminal()

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGWINCH, syscall.SIGUSR1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				WriteCrashLog(r)
				restoreTerminal()
				printCrashNotice()
				startRescueShell()
				os.Exit(2)
			}
		}()
		for s := range sigCh {
			switch s {
			case syscall.SIGWINCH:
				logger.Debugf("Received SIGWINCH terminal resize signal")
				_ = pty.InheritSize(os.Stdin, ptmx) // handle terminal window resize
			// this is the core feature of reloading
			// it helps IRIS reload itself that you dont need to restart the shell manually
			// SIGUSR1 is the signal to active reload when you type "just reload"
			case syscall.SIGUSR1:
				// trigger iris reload by executing itself again
				exe, _ := os.Executable()
				_ = os.Setenv("IRIS_RELOADED", "true")

				innerShell := getActiveInnerShell(c.Process.Pid, shellName)
				if innerShell != "" {
					// to detect which is last shell (bash, zsh, fish)
					_ = os.Setenv("IRIS_ACTIVE_SHELL", innerShell)
				}

				if c.Process != nil {
					cwd, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", c.Process.Pid))
					if err != nil {
						ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
						out, errCmd := exec.CommandContext(ctx, "lsof", "-p", fmt.Sprintf("%d", c.Process.Pid), "-a", "-d", "cwd", "-F", "n").Output()
						cancel()
						if errCmd == nil {
							for _, line := range strings.Split(string(out), "\n") {
								if strings.HasPrefix(line, "n") {
									cwd = strings.TrimSpace(line[1:])
									err = nil
									break
								}
							}
						}
					}
					if err == nil {
						_ = os.Chdir(cwd)
					}
					_ = syscall.Kill(c.Process.Pid, syscall.SIGKILL)
					_ = ptmx.Close()
				}

				restoreTerminal()
				execArgs := []string{os.Args[0]}
				if logDir, err := config.CachePath(); err == nil {
					argsFile := filepath.Join(logDir, "reload-args")
					if data, err := os.ReadFile(argsFile); err == nil {
						lines := strings.Split(string(data), "\n")
						for _, line := range lines {
							trimmed := strings.TrimSpace(line)
							if trimmed != "" {
								execArgs = append(execArgs, trimmed)
							}
						}
						_ = os.Remove(argsFile)
					} else {
						execArgs = os.Args
					}
				} else {
					execArgs = os.Args
				}
				_ = syscall.Exec(exe, execArgs, os.Environ())
			}
		}
	}()

	overlay := integration.NewOverlay()

	// start background update check (async)
	pendingUpdate = startBackgroundUpdateCheck()
	updatePrinted := false

	// bridge pty output to actual stdout
	go func() {
		defer func() {
			if r := recover(); r != nil {
				WriteCrashLog(r)
				restoreTerminal()
				printCrashNotice()
				startRescueShell()
				os.Exit(2)
			}
		}()
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					restoreTerminal()
					os.Exit(0)
				}
				continue
			}
			writeStdout(buf[:n])
		}
	}()

	var disableGhostText atomic.Bool
	disableGhostText.Store(!config.Get().UI.GhostText)
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
		defer func() {
			if r := recover(); r != nil {
				WriteCrashLog(r)
				restoreTerminal()
				printCrashNotice()
				startRescueShell()
				os.Exit(2)
			}
		}()
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
				// hook: after user executes a command, print the update notice exactly once per session
				if !updatePrinted {
					select {
					case result, ok := <-pendingUpdate:
						if ok && result.hasUpdate {
							printUpdateNotice(result.latestVersion)
							updatePrinted = true
						}
					default:
					}
				}
				continue
			}

			if overlay.UserNavigated {
				continue
			}

			// sync local buffer with actual command line
			bufferMu.Lock()
			naiveBuffer = query
			cursorOffset = 0
			bufferMu.Unlock()

			activeModeMu.RLock()
			currentMode := activeMode
			activeModeMu.RUnlock()

			results := MergeResults(query, currentMode)
			if len(results) == 0 {
				writeStdout([]byte(overlay.ClearAndDisable()))
				continue
			}
			writeStdout([]byte(overlay.Clear()))
			overlay.TypedQuery = query
			overlay.UserNavigated = false
			overlay.UpdateItems(results)
			var rBuf strings.Builder
			if !disableGhostText.Load() {
				rBuf.WriteString(overlay.RenderGhostText(query, false))
			}
			rBuf.WriteString(overlay.Render())
			writeStdout([]byte(rBuf.String()))
		}
		if err := scanner.Err(); err != nil {
			logger.Errorf("IPC scanner error: %v", err)
		}
	}()

	suggestionsEnabled := true
	activeModeMu.Lock()
	activeMode = loadMode()
	activeModeMu.Unlock()

	writeStdout([]byte(overlay.Clear()))

	var renderTimer *time.Timer
	var renderMu sync.Mutex

	renderMenuNow := func() {
		if isExecuting() {
			return
		}

		// copy state safely inside timer
		bufferMu.Lock()
		bufCopy := naiveBuffer
		offsetCopy := cursorOffset
		bufferMu.Unlock()

		activeModeMu.RLock()
		modeCopy := activeMode
		activeModeMu.RUnlock()

		navCopy := userNavigated

		runes := []rune(bufCopy)
		if offsetCopy > 0 && offsetCopy <= len(runes) {
			bufCopy = string(runes[:len(runes)-offsetCopy])
		}

		var b strings.Builder
		if !navCopy {
			if bufCopy == "" {
				writeStdout([]byte(overlay.ClearAndDisable()))
				return
			}
			logger.Debugf("Render query: '%s', mode: %s", bufCopy, modeCopy)
			results := MergeResults(bufCopy, modeCopy)
			logger.Debugf("Render results found: %d", len(results))

			if len(results) == 0 || (len(results) == 1 && strings.TrimSpace(results[0].Cmd) == strings.TrimSpace(bufCopy) && !strings.HasSuffix(bufCopy, " ")) {
				b.WriteString(overlay.ClearAndDisable())
				writeStdout([]byte(b.String()))
				return
			}

			if overlay.Visible {
				b.WriteString(overlay.Clear())
			}
			overlay.TypedQuery = bufCopy
			overlay.UpdateItems(results)
		} else {
			if overlay.Visible {
				b.WriteString(overlay.Clear())
			}
		}

		overlay.UserNavigated = navCopy
		if !disableGhostText.Load() {
			b.WriteString(overlay.RenderGhostText(bufCopy, navCopy))
		}
		currentCmd := ""
		if len(overlay.Items) > 0 && overlay.Cursor >= 0 && overlay.Cursor < len(overlay.Items) {
			currentCmd = overlay.Items[overlay.Cursor].Cmd
		}
		logger.Debugf("RenderOverlay nav: %v, cursor: %d, typedQuery: '%s', currentCmd: '%s'", navCopy, overlay.Cursor, overlay.TypedQuery, currentCmd)
		b.WriteString(overlay.Render())
		writeStdout([]byte(b.String()))
	}

	var isNavTimerRunning bool

	renderOverlay = func() {
		renderMu.Lock()
		defer renderMu.Unlock()

		if !suggestionsEnabled || isExecuting() {
			if renderTimer != nil {
				renderTimer.Stop()
				renderTimer = nil
			}
			isNavTimerRunning = false
			return
		}

		navCopy := userNavigated

		if navCopy {
			if isNavTimerRunning {
				return
			}
			if renderTimer != nil {
				renderTimer.Stop()
			}
			isNavTimerRunning = true
			// IMPORTANT: please dont change this to above 24ms or under 19ms
			// I still can get the reason why it only works stably between 20-23ms
			renderTimer = time.AfterFunc(23*time.Millisecond, func() {
				renderMu.Lock()
				renderTimer = nil
				isNavTimerRunning = false
				renderMu.Unlock()
				renderMenuNow()
			})
		} else {
			if renderTimer != nil {
				renderTimer.Stop()
			}
			isNavTimerRunning = false
			renderTimer = time.AfterFunc(25*time.Millisecond, func() {
				renderMu.Lock()
				renderTimer = nil
				renderMu.Unlock()
				renderMenuNow()
			})
		}
	}

	renderOverlay()

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
				_, _ = ptmx.Write(inputSlice[:n])
				continue
			}

			logger.Debugf("Stdin raw input: bytes=%q, hex=%x", inputSlice[:n], inputSlice[:n])

			shouldOverlayDraw := false
			for i := 0; i < n; i++ {
				b := inputSlice[i]
				intercepted := false

				if b == '\033' {
					// check for bracketed paste start/end
					if i+5 < n && inputSlice[i+1] == '[' && inputSlice[i+2] == '2' && inputSlice[i+3] == '0' {
						if (inputSlice[i+4] == '0' || inputSlice[i+4] == '1') && inputSlice[i+5] == '~' {
							intercepted = true
							logger.Debugf("Intercepted bracketed paste event")
							_, _ = ptmx.Write(inputSlice[i : i+6])
							i += 5
							continue
						}
					}
					// handle escape sequences like arrow keys and functional shortcuts
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						// shift tab: hide/unhide menu dropdown
						if inputSlice[i+1] == '[' && inputSlice[i+2] == 'Z' {
							intercepted = true
							suggestionsEnabled = !suggestionsEnabled
							logger.Debugf("Intercepted Shift+Tab, suggestionsEnabled=%v", suggestionsEnabled)
							if !suggestionsEnabled {
								writeStdout([]byte(overlay.ClearAndDisable()))
							} else {
								shouldOverlayDraw = true
							}
							i += 2
							continue
						}

						if overlay.Visible && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') {
							intercepted = true
							userNavigated = true
							overlay.UserNavigated = true

							// clear ghost text synchronously
							if overlay.LastGhostLen > 0 {
								var gs strings.Builder
								gs.WriteString("\0337")
								gs.WriteString(strings.Repeat(" ", overlay.LastGhostLen+10))
								gs.WriteString("\0338")
								writeStdout([]byte(gs.String()))
								overlay.LastGhostLen = 0
							}

							oldCursor := overlay.Cursor
							arrowDir := "down"
							if inputSlice[i+2] == 'A' {
								arrowDir = "up"
								overlay.Cursor--
								if overlay.Cursor < 0 {
									overlay.Cursor = 0
								}
							} else {
								overlay.Cursor++
								if overlay.Cursor >= len(overlay.Items) {
									overlay.Cursor = len(overlay.Items) - 1
								}
							}
							logger.Debugf("Intercepted %s Arrow, cursor moved %d -> %d", arrowDir, oldCursor, overlay.Cursor)

							// boundary hit - ignore redundant write to avoid PTY flooding
							if overlay.Cursor == oldCursor {
								i += 2
								continue
							}

							selected := overlay.Items[overlay.Cursor].Cmd
							bufferMu.Lock()
							naiveBuffer = selected
							cursorOffset = 0
							bufferMu.Unlock()

							_, _ = ptmx.Write(append([]byte{0x15}, selected...))

							userNavigated = true
							overlay.UserNavigated = true
							renderOverlay()
							i += 2
							continue
						} else if !overlay.Visible && naiveBuffer == "" && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') { // up/down arrow on empty prompt
							intercepted = true
							activeModeMu.Lock()
							activeMode = "history"
							saveMode(activeMode)
							activeModeMu.Unlock()

							activeModeMu.RLock()
							currentMode := activeMode
							activeModeMu.RUnlock()
							results := MergeResults("", currentMode)
							if len(results) > 0 {
								limit := 100
								if len(results) < limit {
									limit = len(results)
								}
								var historyList []core.Suggestion

								if inputSlice[i+2] == 'A' {
									// up arrow: reverse the list so newest is at the bottom
									for j := limit - 1; j >= 0; j-- {
										historyList = append(historyList, results[j])
									}
								} else {
									// down arrow: normal order, newest is at the top
									for j := 0; j < limit; j++ {
										historyList = append(historyList, results[j])
									}
								}

								overlay.UpdateItems(historyList)

								if inputSlice[i+2] == 'A' {
									overlay.Cursor = len(historyList) - 1 // up arrow: start at the bottom
								} else {
									overlay.Cursor = 0 // down arrow: start at the top
								}

								selected := overlay.Items[overlay.Cursor].Cmd
								bufferMu.Lock()
								naiveBuffer = selected
								cursorOffset = 0
								bufferMu.Unlock()

								_, _ = ptmx.Write(append([]byte{0x15}, selected...))

								userNavigated = true
								overlay.UserNavigated = true
								renderOverlay()
							}
							i += 2
							continue
						} else if overlay.Visible && !disableGhostText.Load() && inputSlice[i+2] == 'C' { // right arrow
							bufferMu.Lock()
							topCmd := overlay.Items[0].Cmd
							hasMatch := strings.HasPrefix(strings.ToLower(topCmd), strings.ToLower(naiveBuffer))
							var ghostText string
							if hasMatch {
								ghostText = topCmd[len(naiveBuffer):]
							}
							bufferMu.Unlock()

							if hasMatch && len(ghostText) > 0 {
								intercepted = true
								logger.Debugf("Intercepted Right Arrow (accepted ghost text: %q)", ghostText)
								bufferMu.Lock()
								naiveBuffer += ghostText
								cursorOffset = 0
								bufferMu.Unlock()
								_, _ = ptmx.Write([]byte(ghostText))
								shouldOverlayDraw = true
								i += 2
								continue
							}
						}
					}

					// left/right arrow cursor tracking
					isLeftRightArrow := false
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						if inputSlice[i+2] == 'D' {
							bufferMu.Lock()
							isEmptyQuery := (overlay.Visible && overlay.TypedQuery == "") || (!overlay.Visible && naiveBuffer == "")
							bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}
							bufferMu.Lock()
							if naiveBuffer != "" || overlay.Visible {
								cursorOffset++
								if cursorOffset > len(naiveBuffer) {
									cursorOffset = len(naiveBuffer)
								}
								shouldOverlayDraw = true
								userNavigated = false
							}
							bufferMu.Unlock()
							isLeftRightArrow = true
						} else if inputSlice[i+2] == 'C' {
							bufferMu.Lock()
							isEmptyQuery := (overlay.Visible && overlay.TypedQuery == "") || (!overlay.Visible && naiveBuffer == "")
							bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}
							bufferMu.Lock()
							if naiveBuffer != "" || overlay.Visible {
								cursorOffset--
								if cursorOffset < 0 {
									cursorOffset = 0
								}
								shouldOverlayDraw = true
								userNavigated = false
							}
							bufferMu.Unlock()
							isLeftRightArrow = true
						}
					}

					if !intercepted {
						writeStdout([]byte(overlay.ClearAndDisable()))
						disableGhostText.Store(true)
						if !isLeftRightArrow {
							bufferMu.Lock()
							naiveBuffer = ""
							cursorOffset = 0
							bufferMu.Unlock()
						}

						_, _ = ptmx.Write([]byte{b})
						// skip remaining bytes of the escape sequence to avoid misinterpretation
						for j := i + 1; j < n; j++ {
							char := inputSlice[j]
							_, _ = ptmx.Write([]byte{char})
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
					activeModeMu.Lock()
					if activeMode == "spec" {
						activeMode = "history"
					} else {
						activeMode = "spec"
					}
					saveMode(activeMode)
					activeModeMu.Unlock()
					logger.Debugf("Intercepted Ctrl+R, toggled mode to %q", activeMode)
					shouldOverlayDraw = true
					// enter: enter behavior is a bit different from tab suggestions in code editor
					// I want it to execute the command anyway and ignore the suggestions
					// it means only tab to select suggestions, and enter to execute
					// enter is not used to select suggestions
				} else if overlay.Visible && (b == 0x0d || b == 0x0a) {
					intercepted = true
					logger.Debugf("Intercepted Enter key, navigated=%v", overlay.UserNavigated)
					if overlay.UserNavigated && len(overlay.Items) > 0 && overlay.Cursor >= 0 && overlay.Cursor < len(overlay.Items) {
						selected := overlay.Items[overlay.Cursor].Cmd
						activeModeMu.RLock()
						currentMode := activeMode
						activeModeMu.RUnlock()
						if currentMode == "spec" {
							selected = strings.TrimSpace(selected) + " "
						}
						_, _ = ptmx.Write(append([]byte{0x15}, selected...))
					}
					writeStdout([]byte(overlay.ClearAndDisable()))

					_, _ = ptmx.Write([]byte{0x0d})
					bufferMu.Lock()
					naiveBuffer = ""
					cursorOffset = 0
					bufferMu.Unlock()
					disableGhostText.Store(false)
					shouldOverlayDraw = false
					userNavigated = false
					continue
				} else if b == 0x09 { // tab: select suggestions
					intercepted = true
					logger.Debugf("Intercepted Tab key, visible=%v, cursor=%d", overlay.Visible, overlay.Cursor)
					if !overlay.Visible {
						shouldOverlayDraw = true
					} else {
						selected := overlay.Items[overlay.Cursor].Cmd
						writeStdout([]byte(overlay.ClearAndDisable()))

						activeModeMu.RLock()
						currentMode := activeMode
						activeModeMu.RUnlock()
						if currentMode == "spec" {
							selected = strings.TrimSpace(selected) + " "
						}

						bufferMu.Lock()
						naiveBuffer = selected
						cursorOffset = 0
						bufferMu.Unlock()

						_, _ = ptmx.Write(append([]byte{0x15}, selected...))

						overlay.Cursor = 0 // this prevents when you tab, it switches between suggestions non-stop

						shouldOverlayDraw = true // <- rerender after tab to choose, if you set to false,
						// when you press tab continually, it will print all folder from menu suggestions
						// and make the cursor jump to next line
						userNavigated = false
					}
					continue
				}

				if !intercepted {
					_, _ = ptmx.Write([]byte{b})
					// we handle line editing keys manually to keep naiveBuffer in sync
					// since terminal is in raw mode, we must update our state for every change
					switch b {
					case 0x01: // ctrl+a: move to beginning of line
						bufferMu.Lock()
						cursorOffset = len(naiveBuffer)
						if naiveBuffer != "" || overlay.Visible {
							shouldOverlayDraw = true
						}
						bufferMu.Unlock()
						userNavigated = false
					case 0x05: // ctrl+e: move to end of line
						bufferMu.Lock()
						cursorOffset = 0
						if naiveBuffer != "" || overlay.Visible {
							shouldOverlayDraw = true
						}
						bufferMu.Unlock()
						userNavigated = false

					case 127, 0x08: // backspace: remove character
						bufferMu.Lock()
						if len(naiveBuffer) > 0 {
							runes := []rune(naiveBuffer)
							if cursorOffset <= 0 {
								if len(runes) > 0 {
									naiveBuffer = string(runes[:len(runes)-1])
								}
								cursorOffset = 0
							} else {
								if cursorOffset > len(runes) {
									cursorOffset = len(runes)
								}
								pos := len(runes) - cursorOffset
								if pos > 0 && pos <= len(runes) {
									naiveBuffer = string(append(runes[:pos-1], runes[pos:]...))
								}
							}
						}
						bufferMu.Unlock()
						shouldOverlayDraw = true
						userNavigated = false
					case 0x17: // ctrl+w: delete the last word in the buffer
						bufferMu.Lock()
						trimBuf := strings.TrimRight(naiveBuffer, " ")
						lastSpace := strings.LastIndex(trimBuf, " ")
						if lastSpace >= 0 {
							naiveBuffer = trimBuf[:lastSpace+1]
						} else {
							naiveBuffer = ""
						}
						cursorOffset = 0
						bufferMu.Unlock()
						shouldOverlayDraw = true
						userNavigated = false
					case 0x0c: // ctrl+l: clear screen but keep buffer and redraw menu
						shouldOverlayDraw = true
						userNavigated = false
					case '\r', '\n', 0x03, 0x15: // enter, ctrl+c, ctrl+u: clear buffer on line reset
						bufferMu.Lock()
						naiveBuffer = ""
						cursorOffset = 0
						bufferMu.Unlock()
						activeModeMu.Lock()
						activeMode = loadMode()
						activeModeMu.Unlock()
						disableGhostText.Store(false)
						writeStdout([]byte(overlay.ClearAndDisable()))
						userNavigated = false
					default:
						// track normal printable characters in the buffer for matching
						if b >= 32 && b <= 126 {
							// if user presses space, check if the current word is an alias
							bufferMu.Lock()
							isSpaceAlias := b == ' ' && naiveBuffer != "" && !strings.Contains(naiveBuffer, " ")
							var target string
							var ok bool
							if isSpaceAlias {
								target, ok = core.GetAlias(naiveBuffer)
							}
							bufferMu.Unlock()

							if isSpaceAlias && ok {
								// clear the current alias and replace it with the full command
								_, _ = ptmx.Write(append([]byte{0x15}, target+" "...))
								bufferMu.Lock()
								naiveBuffer = target + " "
								cursorOffset = 0
								bufferMu.Unlock()
								shouldOverlayDraw = true
								continue
							}
							bufferMu.Lock()
							if cursorOffset == 0 {
								naiveBuffer += string(b)
							} else {
								if cursorOffset > len(naiveBuffer) {
									cursorOffset = len(naiveBuffer)
								}
								pos := len(naiveBuffer) - cursorOffset
								if pos >= 0 && pos <= len(naiveBuffer) {
									naiveBuffer = naiveBuffer[:pos] + string(b) + naiveBuffer[pos:]
								} else {
									naiveBuffer += string(b)
									cursorOffset = 0
								}
							}
							bufferMu.Unlock()
							shouldOverlayDraw = true
							userNavigated = false
							overlay.UserNavigated = false
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
