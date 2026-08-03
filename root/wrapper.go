package root

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/integration/shell"
	"github.com/versenilvis/iris/internal/ai"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/logger"
	"github.com/versenilvis/iris/internal/scoring"
	"github.com/versenilvis/iris/spec"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var (
	prevRecordedCommand string
	prevCmdCwd          string
	prevCmdMu           sync.Mutex
)

func getPrevSkeleton() string {
	prevCmdMu.Lock()
	defer prevCmdMu.Unlock()
	if prevRecordedCommand == "" {
		return ""
	}
	return scoring.ExtractSkeleton(prevRecordedCommand)
}

func getPrevRecordedInfo() (string, string) {
	prevCmdMu.Lock()
	defer prevCmdMu.Unlock()
	if prevRecordedCommand == "" {
		return "", ""
	}
	return scoring.ExtractSkeleton(prevRecordedCommand), prevCmdCwd
}

func setPrevRecordedInfo(cmd, cwd string) {
	prevCmdMu.Lock()
	defer prevCmdMu.Unlock()
	prevRecordedCommand = cmd
	prevCmdCwd = cwd
}

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

func shellArgs(login bool) []string {
	if login {
		return []string{"--login"}
	}
	return nil
}

var (
	oldState     *term.State
	oldStateFd   int
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
		_ = term.Restore(oldStateFd, oldState)
		oldState = nil
	}
}

// runWrapper sets up the pty environment, launches the shell,
// and manages the main input loop to provide real-time suggestions
// it handles raw terminal mode to intercept keystrokes and
// coordinates between the shell process and the suggestion overlay
func runWrapper() {
	var naiveBuffer string
	var lastSubmittedCommand string
	cursorOffset := 0
	var bufferMu sync.Mutex
	var userNavigated atomic.Bool
	var renderMenuNow func()

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
	c := exec.CommandContext(ctx, adapter.GetShellPath(), shellArgs(config.Get().Core.ShellLogin)...)
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

	stdinFile := os.Stdin
	if !term.IsTerminal(int(stdinFile.Fd())) {
		if tty, ttyErr := os.OpenFile("/dev/tty", os.O_RDWR, 0); ttyErr == nil {
			stdinFile = tty
		}
	}

	_ = pty.InheritSize(stdinFile, ptmx)
	spec.ShellPID = c.Process.Pid

	logger.Infof("PTY child shell started: shell=%s, path=%s, pid=%d", shellName, adapter.GetShellPath(), c.Process.Pid)

	// put terminal in raw mode to intercept every keystroke
	var errMakeRaw error
	if term.IsTerminal(int(stdinFile.Fd())) {
		oldState, errMakeRaw = term.MakeRaw(int(stdinFile.Fd()))
		if errMakeRaw != nil {
			logger.Errorf("Failed to set terminal raw mode: %v", errMakeRaw)
			panic(errMakeRaw)
		}
		oldStateFd = int(stdinFile.Fd())
		logger.Debugf("Terminal set to raw mode successfully")
		defer func() {
			oldStateMu.Lock()
			defer oldStateMu.Unlock()
			if oldState != nil {
				_ = term.Restore(int(stdinFile.Fd()), oldState)
				oldState = nil
			}
		}()
	} else {
		logger.Warnf("stdinFile is not a terminal, skipping raw mode")
	}

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
				_ = pty.InheritSize(stdinFile, ptmx) // handle terminal window resize
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
					cwd, linkErr := os.Readlink(fmt.Sprintf("/proc/%d/cwd", c.Process.Pid))
					if linkErr != nil {
						ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
						out, errCmd := exec.CommandContext(ctx, "lsof", "-p", fmt.Sprintf("%d", c.Process.Pid), "-a", "-d", "cwd", "-F", "n").Output()
						cancel()
						if errCmd == nil {
							for line := range strings.SplitSeq(string(out), "\n") {
								if strings.HasPrefix(line, "n") {
									cwd = strings.TrimSpace(line[1:])
									linkErr = nil
									break
								}
							}
						}
					}
					if linkErr == nil {
						_ = os.Chdir(cwd)
					}
					_ = syscall.Kill(c.Process.Pid, syscall.SIGKILL)
					_ = ptmx.Close()
				}

				restoreTerminal()
				execArgs := []string{os.Args[0]}
				if logDir, pathErr := config.CachePath(); pathErr == nil {
					argsFile := filepath.Join(logDir, "reload-args")
					if data, readErr := os.ReadFile(argsFile); readErr == nil {
						lines := strings.SplitSeq(string(data), "\n")
						for line := range lines {
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

	shellPGID, err := unix.Getpgid(spec.ShellPID)
	if err != nil {
		shellPGID = spec.ShellPID
	}
	var isCommandActive atomic.Bool
	var disableGhostText atomic.Bool
	disableGhostText.Store(!config.Get().UI.GhostText)
	var renderOverlayFn atomic.Value // holds func()
	renderOverlayFn.Store(func() {})
	config.AutoDetectConfigChange(func(cfg *config.Config) {
		disableGhostText.Store(!cfg.UI.GhostText)
		renderOverlayFn.Load().(func())()
	})
	isExecuting := func() bool {
		if isCommandActive.Load() {
			// for bash: no preexec/precmd hooks, so fall back to TIOCGPGRP to detect when shell returns
			if shellName == "bash" {
				pgrp, pgrpErr := unix.IoctlGetInt(int(ptmx.Fd()), unix.TIOCGPGRP)
				if pgrpErr == nil && pgrp == shellPGID {
					isCommandActive.Store(false)
					return false
				}
			}
			return true
		}
		pgrp, err := unix.IoctlGetInt(int(ptmx.Fd()), unix.TIOCGPGRP)
		if err != nil {
			return false
		}
		return pgrp != shellPGID
	}

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
		var lastPromptBuf []byte
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				restoreTerminal()
				if err == io.EOF || errors.Is(err, syscall.EIO) || strings.Contains(err.Error(), "input/output error") {
					os.Exit(0)
				}
				logger.Errorf("Unexpected PTY read error: %v", err)
				os.Exit(1)
			}
			writeStdout(buf[:n])

			bufferMu.Lock()
			nbEmpty := naiveBuffer == ""
			navigated := userNavigated.Load()
			bufferMu.Unlock()

			if isExecuting() {
				lastPromptBuf = nil
			} else if nbEmpty && !navigated {
				lastPromptBuf = append(lastPromptBuf, buf[:n]...)
				if idx := bytes.LastIndexByte(lastPromptBuf, '\n'); idx >= 0 {
					lastPromptBuf = append([]byte(nil), lastPromptBuf[idx+1:]...)
				}
				pLen := integration.ComputeCursorCol(lastPromptBuf)
				if pLen >= 0 {
					overlay.SetPromptLen(pLen)
				}
			}
		}
	}()



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

			if cwd, ok := strings.CutPrefix(query, "IRIS_CWD:"); ok {
				spec.SetCWD(cwd)
				continue
			}

			if query == "IRIS_CMD_START" {
				isCommandActive.Store(true)
				bufferMu.Lock()
				naiveBuffer = ""
				cursorOffset = 0
				bufferMu.Unlock()
				writeStdout([]byte(overlay.ClearAndDisable()))
				SetCurrentAISuggestion(nil)
				continue
			}

			if query == "IRIS_CMD_STOP" || strings.HasPrefix(query, "IRIS_CMD_STOP:") {
				exitCode := 0
				if after, ok := strings.CutPrefix(query, "IRIS_CMD_STOP:"); ok {
					if code, err := strconv.Atoi(after); err == nil {
						exitCode = code
					}
				}
				isCommandActive.Store(false)
				SetCurrentAISuggestion(nil)
				bufferMu.Lock()
				cmdToRecord := lastSubmittedCommand
				lastSubmittedCommand = ""
				bufferMu.Unlock()
				if cmdToRecord != "" {
					cwd := spec.GetCWD()
					prevSkeleton, prevCwd := getPrevRecordedInfo()
					currSkeleton := scoring.ExtractSkeleton(cmdToRecord)
					go func(c, d string, code int, pSkel, pCwd, cSkel string) {
						defer func() {
							if r := recover(); r != nil {
								WriteCrashLog(r)
							}
						}()
						ctxRecord, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
						defer cancel()
						if store, err := scoring.GetFrecencyStore(); err == nil && store != nil {
							_ = store.Record(ctxRecord, c, d, code)
							if pSkel != "" && cSkel != "" {
								_ = store.RecordTransition(ctxRecord, pSkel, cSkel, d, code)
							}
						}
					}(cmdToRecord, cwd, exitCode, prevSkeleton, prevCwd, currSkeleton)
					setPrevRecordedInfo(cmdToRecord, cwd)
				}
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

			isCommandActive.Store(false)

			if overlay.GetUserNavigated() {
				continue
			}

			if query == "" {
				bufferMu.Lock()
				wasEmpty := naiveBuffer == ""
				naiveBuffer = ""
				cursorOffset = 0
				bufferMu.Unlock()
				if !wasEmpty {
					writeStdout([]byte(overlay.ClearAndDisable()))
					SetCurrentAISuggestion(nil)
				}
				continue
			}

			bufferMu.Lock()
			if naiveBuffer == query {
				bufferMu.Unlock()
				continue
			}
			naiveBuffer = query
			cursorOffset = 0
			bufferMu.Unlock()

			renderOverlayFn.Load().(func())()
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
	var aiTimer *time.Timer
	var aiCancel context.CancelFunc
	var aiMu sync.Mutex

	renderMenuNow = func() {
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

		navCopy := userNavigated.Load()

		runes := []rune(bufCopy)
		if offsetCopy > 0 && offsetCopy <= len(runes) {
			bufCopy = string(runes[:len(runes)-offsetCopy])
		}

		aiMu.Lock()
		if aiTimer != nil {
			aiTimer.Stop()
		}
		if aiCancel != nil {
			aiCancel()
			aiCancel = nil
		}
		if config.Get().AI.Enabled && bufCopy != "" && !navCopy && offsetCopy == 0 {
			queryTarget := bufCopy
			debounceMS := config.Get().AI.DebounceMS
			if debounceMS <= 0 {
				debounceMS = 500
			}
			aiTimer = time.AfterFunc(time.Duration(debounceMS)*time.Millisecond, func() {
				// Require at least 3 characters to trigger AI completion to save API quota and avoid 6000 TPM limit (Groq api docs)
				if len(strings.TrimSpace(queryTarget)) < 3 {
					return
				}
				aiMu.Lock()
				ctx, cancel := context.WithCancel(context.Background())
				aiCancel = cancel
				aiMu.Unlock()
				defer cancel()

				cwd := spec.GetCWD()
				var recentCmds []string
				var lastCmd string
				if hist, err := integration.SearchHistory("", nil); err == nil {
					// Limit to 3 recent commands to keep prompt concise and reduce token consumption
					for i := 0; i < len(hist) && i < 3; i++ {
						recentCmds = append(recentCmds, hist[i].Cmd)
					}
					if len(recentCmds) > 0 {
						lastCmd = recentCmds[0]
					}
				}
				env := ai.NewEnvSnapshot(cwd, lastCmd, 0, recentCmds)
				sugg, err := GetAIEngine().Suggest(ctx, queryTarget, env, "")
				if err != nil || sugg == nil || ctx.Err() != nil {
					return
				}
				SetCurrentAISuggestion(sugg)
				if overlay.InjectAISuggestion(*sugg) {
					renderOverlayFn.Load().(func())()
				}
			})
		}
		aiMu.Unlock()

		var b strings.Builder
		if !navCopy {
			if bufCopy == "" && !overlay.IsVisible() {
				writeStdout([]byte(overlay.ClearAndDisable()))
				return
			}
			logger.Debugf("Render query: '%s', mode: %s", bufCopy, modeCopy)
			results := MergeResults(bufCopy, modeCopy)
			logger.Debugf("Render results found: %d", len(results))

			if len(results) == 0 || (len(results) == 1 && strings.TrimSpace(results[0].Cmd) == strings.TrimSpace(bufCopy) && !strings.HasSuffix(bufCopy, " ")) {
				b.WriteString(overlay.HideMenu(bufCopy))
				writeStdout([]byte(b.String()))
				return
			}

			if overlay.IsVisible() {
				b.WriteString(overlay.Clear())
			}
			overlay.SetQueryAndItems(bufCopy, results)
		} else {
			if overlay.IsVisible() {
				b.WriteString(overlay.Clear())
			}
		}

		overlay.SetUserNavigated(navCopy)
		if !disableGhostText.Load() {
			b.WriteString(overlay.RenderGhostText(bufCopy, navCopy, offsetCopy == 0))
		}
		currentCmd := overlay.GetCurrentCmd()
		logger.Debugf("RenderOverlay nav: %v, typedQuery: '%s', currentCmd: '%s'", navCopy, overlay.GetTypedQuery(), currentCmd)
		b.WriteString(overlay.Render())
		writeStdout([]byte(b.String()))
	}

	renderOverlayFn.Store(func() {
		renderMu.Lock()
		defer renderMu.Unlock()

		if !suggestionsEnabled || isExecuting() {
			if renderTimer != nil {
				renderTimer.Stop()
				renderTimer = nil
			}
			return
		}

		if userNavigated.Load() {
			return
		}

		if renderTimer != nil {
			return
		}
		renderTimer = time.AfterFunc(20*time.Millisecond, func() {
			renderMu.Lock()
			renderTimer = nil
			renderMu.Unlock()
			renderMenuNow()
		})
	})

	renderOverlayFn.Load().(func())()

	// reads from stdin and decides what to forward or intercept
	// for most cases, I just handle the already have terminal shortcuts
	// for some shortcuts like tab, enter, shift tab, ctrl r,
	// they have a little bit different behavior to match our tool
	inBracketedPaste := false
	for {
		inputSlice := make([]byte, 128)
		n, err := stdinFile.Read(inputSlice)
		if err != nil {
			break
		}

		if n > 0 {
			if isExecuting() {
				inBracketedPaste = false
				_, _ = ptmx.Write(inputSlice[:n])
				continue
			}

			logger.Debugf("Stdin raw input: bytes=%q, hex=%x", inputSlice[:n], inputSlice[:n])

			shouldOverlayDraw := false
			for i := 0; i < n; i++ {
				b := inputSlice[i]
				intercepted := false

				if matched, consumed := config.MatchKey(inputSlice[i:], config.Get().Keybindings.ToggleMenu); matched {
					intercepted = true
					suggestionsEnabled = !suggestionsEnabled
					logger.Debugf("Intercepted ToggleMenu, suggestionsEnabled=%v", suggestionsEnabled)
					if !suggestionsEnabled {
						writeStdout([]byte(overlay.ClearAndDisable()))
					} else {
						shouldOverlayDraw = true
					}
					i += consumed - 1
					continue
				}

				if matched, consumed := config.MatchKey(inputSlice[i:], config.Get().Keybindings.ToggleMode); matched { // ctrl+r: toggle between command specs and command history
					i += consumed - 1
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
					if userNavigated.Load() {
						bufferMu.Lock()
						naiveBuffer = overlay.GetTypedQuery()
						cursorOffset = 0
						bufferMu.Unlock()
						_, _ = ptmx.Write(append([]byte{0x15}, overlay.GetTypedQuery()...))
					}
					userNavigated.Store(false)
					overlay.Show()
					shouldOverlayDraw = true
					continue
				}

				var isNavUp, isNavDown bool
				var navConsumed int
				if isNavUp, navConsumed = config.MatchKey(inputSlice[i:], config.Get().Keybindings.NavigateUp); !isNavUp {
					isNavDown, navConsumed = config.MatchKey(inputSlice[i:], config.Get().Keybindings.NavigateDown)
				}

				if isNavUp || isNavDown {
					if overlay.IsVisible() {
						intercepted = true
						userNavigated.Store(true)

						arrowDir := "down"
						if isNavUp {
							arrowDir = "up"
						}
						moved, selectedCmd := overlay.MoveCursor(arrowDir)
						if !moved {
							i += navConsumed - 1
							continue
						}

						bufferMu.Lock()
						activeModeMu.RLock()
						isHistMode := activeMode == "history"
						activeModeMu.RUnlock()
						var toWrite []byte
						if isHistMode && selectedCmd != "" {
							naiveBuffer = selectedCmd
							cursorOffset = 0
							toWrite = append([]byte{0x15}, selectedCmd...)
						}
						bufCopy := naiveBuffer
						offsetCopy := cursorOffset
						bufferMu.Unlock()

						if len(toWrite) > 0 {
							_, _ = ptmx.Write(toWrite)
						}

						var b strings.Builder
						if !disableGhostText.Load() {
							b.WriteString(overlay.RenderGhostText(bufCopy, true, offsetCopy == 0))
						}
						b.WriteString(overlay.Render())
						writeStdout([]byte(b.String()))

						i += navConsumed - 1
						continue
					} else if suggestionsEnabled {
						// hidden-overlay history navigation (only when suggestions are enabled;
						// otherwise let the navigation keys pass through to the shell)
						intercepted = true
						activeModeMu.Lock()
						if activeMode == "" {
							activeMode = loadMode()
						}
						activeModeMu.Unlock()

						activeModeMu.RLock()
						currentMode := activeMode
						activeModeMu.RUnlock()

						bufferMu.Lock()
						bufQuery := naiveBuffer
						bufferMu.Unlock()

						results := MergeResults(bufQuery, currentMode)
						if len(results) > 0 {
							limit := min(len(results), 100)
							var historyList []spec.Suggestion

							if isNavUp {
								for j := limit - 1; j >= 0; j-- {
									historyList = append(historyList, results[j])
								}
							} else {
								for j := range limit {
									historyList = append(historyList, results[j])
								}
							}

							selected := overlay.SetHistoryList(historyList, isNavUp)
							if selected != "" {
								bufferMu.Lock()
								naiveBuffer = selected
								cursorOffset = 0
								bufferMu.Unlock()

								userNavigated.Store(true)
								writeStdout([]byte(overlay.Render()))
								_, _ = ptmx.Write(append([]byte{0x15}, selected...))
							}
						}
						i += navConsumed - 1
						continue
					}
				}

				if matched, consumed := config.MatchKey(inputSlice[i:], config.Get().Keybindings.SelectSuggestion); matched && config.Get().Keybindings.SelectSuggestion != "" {
					if overlay.IsVisible() {
						intercepted = true
						selected := overlay.GetCurrentCmd()
						if selected != "" {
							activeModeMu.RLock()
							currentMode := activeMode
							activeModeMu.RUnlock()
							if currentMode == "spec" {
								s := strings.TrimSpace(selected)
								if strings.HasSuffix(s, "/") || strings.HasSuffix(s, "\\") {
									selected = s
								} else {
									selected = s + " "
								}
							}
							bufferMu.Lock()
							naiveBuffer = selected
							cursorOffset = 0
							bufferMu.Unlock()
							_, _ = ptmx.Write(append([]byte{0x15}, selected...))
							
							overlay.ClearGhostTextState()
							userNavigated.Store(false)
							writeStdout([]byte(overlay.Render()))
						}
					}
					// always consume the full binding atomically, even when the overlay is hidden
					i += consumed - 1
					continue
				}

				if b == 0x0d || b == 0x0a { // enter
					intercepted = true
					logger.Debugf("Intercepted Enter key")

					if b == 0x0d && i+1 < n && inputSlice[i+1] == 0x0a {
						i++ // consume trailing \n in \r\n to prevent matching ctrl+j
					}
					
					var selectedCmd string
					var shouldAutoExecute bool
					if overlay.IsVisible() && (config.Get().Core.AutoExecute || userNavigated.Load()) {
						selectedCmd = overlay.GetCurrentCmd()
						if selectedCmd != "" {
							shouldAutoExecute = true
						}
					}

					writeStdout([]byte(overlay.ClearAndDisable()))
					SetCurrentAISuggestion(nil)
					renderMu.Lock()
					if renderTimer != nil {
						renderTimer.Stop()
						renderTimer = nil
					}
					renderMu.Unlock()

					var cmdToSubmit string
					if shouldAutoExecute {
						activeModeMu.RLock()
						currentMode := activeMode
						activeModeMu.RUnlock()
						if currentMode == "spec" {
							s := strings.TrimSpace(selectedCmd)
							if strings.HasSuffix(s, "/") || strings.HasSuffix(s, "\\") {
								selectedCmd = s
							} else {
								selectedCmd = s + " "
							}
						}
						// update the line first
						_, _ = ptmx.Write(append([]byte{0x15}, selectedCmd...))
						cmdToSubmit = selectedCmd
					} else {
						bufferMu.Lock()
						cmdToSubmit = naiveBuffer
						bufferMu.Unlock()
					}

					if strings.TrimSpace(cmdToSubmit) == "iris reload" {
						if newCfg, err := config.Load(); err == nil {
							config.Init(newCfg)
							disableGhostText.Store(!newCfg.UI.GhostText)
						}
						msg := "echo -e '\\033[32m✓ Iris configuration reloaded successfully.\\033[0m'\r"
						_, _ = ptmx.Write(append([]byte{0x15}, []byte(msg)...))
						bufferMu.Lock()
						naiveBuffer = ""
						cursorOffset = 0
						bufferMu.Unlock()
						disableGhostText.Store(false)
						shouldOverlayDraw = false
						userNavigated.Store(false)
						continue
					}

					integration.RecordSessionCommand(cmdToSubmit)
					bufferMu.Lock()
					lastSubmittedCommand = strings.TrimSpace(cmdToSubmit)
					naiveBuffer = ""
					cursorOffset = 0
					bufferMu.Unlock()
					isCommandActive.Store(true)
					_, _ = ptmx.Write([]byte{b}) // forward enter to terminal
					disableGhostText.Store(false)
					shouldOverlayDraw = false
					userNavigated.Store(false)
					continue
				}

				if b == '\033' {
					// check for bracketed paste start/end
					if i+5 < n && inputSlice[i+1] == '[' && inputSlice[i+2] == '2' && inputSlice[i+3] == '0' {
						if (inputSlice[i+4] == '0' || inputSlice[i+4] == '1') && inputSlice[i+5] == '~' {
							intercepted = true
							inBracketedPaste = inputSlice[i+4] == '0'
							logger.Debugf("Intercepted bracketed paste event inPaste=%v", inBracketedPaste)
							_, _ = ptmx.Write(inputSlice[i : i+6])
							i += 5
							continue
						}
					}
					// handle escape sequences like arrow keys and functional shortcuts
					// left/right arrow cursor tracking
					isLeftRightArrow := false
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						if inputSlice[i+2] == 'D' {
							bufferMu.Lock()
							isEmptyQuery := naiveBuffer == "" && (!overlay.IsVisible() || overlay.GetTypedQuery() == "")
							bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}
							bufferMu.Lock()
							if naiveBuffer != "" || overlay.IsVisible() {
								cursorOffset++
								if cursorOffset > len(naiveBuffer) {
									cursorOffset = len(naiveBuffer)
								}
								shouldOverlayDraw = true
								userNavigated.Store(false)
							}
							bufferMu.Unlock()
							isLeftRightArrow = true
						} else if inputSlice[i+2] == 'C' {
							bufferMu.Lock()
							isEmptyQuery := naiveBuffer == "" && (!overlay.IsVisible() || overlay.GetTypedQuery() == "")
							bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}

							bufferMu.Lock()
							atEnd := (cursorOffset == 0)
							ghostText := ""
							if !disableGhostText.Load() {
								ghostText = overlay.GetGhostText(naiveBuffer, atEnd)
							}
							bufferMu.Unlock()

							if len(ghostText) > 0 {
								intercepted = true
								logger.Debugf("Intercepted Right Arrow (accepted ghost text: %q)", ghostText)
								bufferMu.Lock()
								naiveBuffer += ghostText
								cursorOffset = 0
								bufferMu.Unlock()
								overlay.ClearGhostTextState()
								_, _ = ptmx.Write([]byte(ghostText))
								shouldOverlayDraw = true
								i += 2
								continue
							}

							bufferMu.Lock()
							if naiveBuffer != "" || overlay.IsVisible() {
								cursorOffset--
								if cursorOffset < 0 {
									cursorOffset = 0
								}
								shouldOverlayDraw = true
								userNavigated.Store(false)
							}
							bufferMu.Unlock()
							isLeftRightArrow = true
						}
					}

					if !intercepted {
						writeStdout([]byte(overlay.ClearAndDisable()))
						disableGhostText.Store(true)
						
						// If it's a standalone ESC (n==1), don't clear the buffer because the user just wanted to hide the menu or enter vi-mode
						isStandaloneEsc := n == 1 && b == '\033'
						if !isLeftRightArrow && !isStandaloneEsc {
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
				if b == 0x03 || b == 0x15 { // ctrl+c, ctrl+u
					intercepted = true
					writeStdout([]byte(overlay.ClearAndDisable()))
					SetCurrentAISuggestion(nil)
					renderMu.Lock()
					if renderTimer != nil {
						renderTimer.Stop()
						renderTimer = nil
					}
					renderMu.Unlock()
					isCommandActive.Store(false)
					_, _ = ptmx.Write([]byte{b})
					bufferMu.Lock()
					naiveBuffer = ""
					cursorOffset = 0
					bufferMu.Unlock()
					disableGhostText.Store(false)
					shouldOverlayDraw = false
					userNavigated.Store(false)
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
						if naiveBuffer != "" || overlay.IsVisible() {
							shouldOverlayDraw = true
						}
						bufferMu.Unlock()
						userNavigated.Store(false)
					case 0x05: // ctrl+e: move to end of line
						bufferMu.Lock()
						cursorOffset = 0
						if naiveBuffer != "" || overlay.IsVisible() {
							shouldOverlayDraw = true
						}
						bufferMu.Unlock()
						userNavigated.Store(false)

					case 127, 0x08: // backspace: remove character
						bufferMu.Lock()
						wasEmpty := len(naiveBuffer) == 0
						if !wasEmpty {
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
						isEmptyNow := len(naiveBuffer) == 0
						bufferMu.Unlock()

						if wasEmpty || isEmptyNow {
							writeStdout([]byte(overlay.ClearAndDisable()))
							userNavigated.Store(false)
							continue
						}
						shouldOverlayDraw = true
						userNavigated.Store(false)
					case 0x17: // ctrl+w: delete the last word in the buffer
						bufferMu.Lock()
						wasEmpty := len(naiveBuffer) == 0
						trimBuf := strings.TrimRight(naiveBuffer, " ")
						lastSpace := strings.LastIndex(trimBuf, " ")
						if lastSpace >= 0 {
							naiveBuffer = trimBuf[:lastSpace+1]
						} else {
							naiveBuffer = ""
						}
						cursorOffset = 0
						isEmptyNow := len(naiveBuffer) == 0
						bufferMu.Unlock()

						if wasEmpty || isEmptyNow {
							writeStdout([]byte(overlay.ClearAndDisable()))
							userNavigated.Store(false)
							continue
						}
						shouldOverlayDraw = true
						userNavigated.Store(false)
					case 0x0c: // ctrl+l: clear screen but keep buffer and redraw menu
						shouldOverlayDraw = true
						userNavigated.Store(false)
					case '\r', '\n', 0x03, 0x15: // enter, ctrl+c, ctrl+u: clear buffer on line reset
						inBracketedPaste = false
						bufferMu.Lock()
						naiveBuffer = ""
						cursorOffset = 0
						bufferMu.Unlock()
						activeModeMu.Lock()
						activeMode = loadMode()
						activeModeMu.Unlock()
						disableGhostText.Store(false)
						writeStdout([]byte(overlay.ClearAndDisable()))
						SetCurrentAISuggestion(nil)
						userNavigated.Store(false)
					default:
						// track normal printable characters in the buffer for matching
						if b >= 32 && b <= 126 {
							// expand alias on space, but only when typing manually (not pasting)
							// and only if expand-alias configuration is enabled
							bufferMu.Lock()
							isSpaceAlias := config.Get().Core.ExpandAlias && !inBracketedPaste && b == ' ' && naiveBuffer != "" && !strings.Contains(naiveBuffer, " ")
							var target string
							var ok bool
							if isSpaceAlias {
								target, ok = spec.GetAlias(naiveBuffer)
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
							userNavigated.Store(false)
							overlay.SetUserNavigated(false)
						}
					}
				}
			}
			if shouldOverlayDraw {
				renderOverlayFn.Load().(func())()
			}
		}
	}
}
