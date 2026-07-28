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
type WrapperSession struct {
	ptmx      *os.File
	cmd       *exec.Cmd
	shellName string
	overlay   *integration.Overlay

	naiveBuffer          string
	lastSubmittedCommand string
	cursorOffset         int
	bufferMu             sync.Mutex
	userNavigated        atomic.Bool
	isCommandActive      atomic.Bool
	suggestionsEnabled   atomic.Bool
	disableGhostText     atomic.Bool
	updatePrinted        bool
	pendingUpdate        <-chan updateResult
	shellPGID            int

	renderTimer *time.Timer
	renderMu    sync.Mutex
	aiTimer     *time.Timer
	aiCancel    context.CancelFunc
	aiMu        sync.Mutex

	inBracketedPaste bool
}

func runWrapper() {
	r, w, err := os.Pipe()
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
	c.ExtraFiles[10] = w
	c.Env = adapter.GetEnv(13, os.Getpid())

	ptmx, err := pty.Start(c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[IRIS] failed to start PTY: %v\n", err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	_ = pty.InheritSize(os.Stdin, ptmx)
	spec.ShellPID = c.Process.Pid

	logger.Infof("PTY child shell started: shell=%s, path=%s, pid=%d", shellName, adapter.GetShellPath(), c.Process.Pid)

	var errMakeRaw error
	oldState, errMakeRaw = term.MakeRaw(int(os.Stdin.Fd()))
	if errMakeRaw != nil {
		logger.Errorf("Failed to set terminal raw mode: %v", errMakeRaw)
		panic(errMakeRaw)
	}
	logger.Debugf("Terminal set to raw mode successfully")
	defer restoreTerminal()

	startSignalMonitor(ptmx, c, shellName)

	shellPGID, err := unix.Getpgid(spec.ShellPID)
	if err != nil {
		shellPGID = spec.ShellPID
	}

	session := &WrapperSession{
		ptmx:               ptmx,
		cmd:                c,
		shellName:          shellName,
		overlay:            integration.NewOverlay(),
		shellPGID:          shellPGID,
	}
	session.suggestionsEnabled.Store(true)
	session.disableGhostText.Store(!config.Get().UI.GhostText)
	session.pendingUpdate = startBackgroundUpdateCheck()

	go session.startStdoutBridge()
	go session.startIPCListener(r)

	activeModeMu.Lock()
	activeMode = loadMode()
	activeModeMu.Unlock()

	writeStdout([]byte(session.overlay.Clear()))

	session.renderOverlay()
	session.handleInputLoop()
}

func (s *WrapperSession) isExecuting() bool {
	if s.isCommandActive.Load() {
		return true
	}
	pgrp, err := unix.IoctlGetInt(int(s.ptmx.Fd()), unix.TIOCGPGRP)
	if err != nil {
		return false
	}
	return pgrp != s.shellPGID
}

func (s *WrapperSession) bufferEmpty() bool {
	s.bufferMu.Lock()
	defer s.bufferMu.Unlock()
	return s.naiveBuffer == ""
}

func (s *WrapperSession) startStdoutBridge() {
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
		n, err := s.ptmx.Read(buf)
		if err != nil {
			if err == io.EOF {
				restoreTerminal()
				os.Exit(0)
			}
			restoreTerminal()
			break
		}
		writeStdout(buf[:n])

		s.bufferMu.Lock()
		nbEmpty := s.naiveBuffer == ""
		navigated := s.userNavigated.Load()
		s.bufferMu.Unlock()

		if s.isExecuting() {
			lastPromptBuf = nil
		} else if nbEmpty && !navigated {
			lastPromptBuf = append(lastPromptBuf, buf[:n]...)
			if idx := bytes.LastIndexByte(lastPromptBuf, '\n'); idx >= 0 {
				lastPromptBuf = append([]byte(nil), lastPromptBuf[idx+1:]...)
			}
			pLen := integration.ComputeCursorCol(lastPromptBuf)
			if pLen >= 0 {
				s.overlay.SetPromptLen(pLen)
			}
		}
	}
}

func (s *WrapperSession) startIPCListener(r *os.File) {
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

		if query == "IRIS_CMD_START" {
			s.isCommandActive.Store(true)
			s.bufferMu.Lock()
			s.naiveBuffer = ""
			s.cursorOffset = 0
			s.bufferMu.Unlock()
			writeStdout([]byte(s.overlay.ClearAndDisable()))
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
			s.isCommandActive.Store(false)
			SetCurrentAISuggestion(nil)
			s.bufferMu.Lock()
			cmdToRecord := s.lastSubmittedCommand
			s.lastSubmittedCommand = ""
			s.bufferMu.Unlock()
			if cmdToRecord != "" {
				integration.RecordSessionCommand(cmdToRecord)
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
			if !s.updatePrinted {
				select {
				case result, ok := <-s.pendingUpdate:
					if ok && result.hasUpdate {
						printUpdateNotice(result.latestVersion)
						s.updatePrinted = true
					}
				default:
				}
			}
			continue
		}

		s.isCommandActive.Store(false)

		if s.overlay.GetUserNavigated() {
			continue
		}

		if query == "" {
			s.bufferMu.Lock()
			wasEmpty := s.naiveBuffer == ""
			s.naiveBuffer = ""
			s.cursorOffset = 0
			s.bufferMu.Unlock()
			if !wasEmpty {
				writeStdout([]byte(s.overlay.ClearAndDisable()))
				SetCurrentAISuggestion(nil)
			}
			continue
		}

		s.bufferMu.Lock()
		if s.naiveBuffer == query {
			s.bufferMu.Unlock()
			continue
		}
		s.naiveBuffer = query
		s.cursorOffset = 0
		s.bufferMu.Unlock()

		s.renderOverlay()
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("IPC scanner error: %v", err)
	}
}

func (s *WrapperSession) renderMenuNow() {
	if s.isExecuting() {
		return
	}

	s.bufferMu.Lock()
	bufCopy := s.naiveBuffer
	offsetCopy := s.cursorOffset
	s.bufferMu.Unlock()

	activeModeMu.RLock()
	modeCopy := activeMode
	activeModeMu.RUnlock()

	navCopy := s.userNavigated.Load()

	runes := []rune(bufCopy)
	if offsetCopy > 0 && offsetCopy <= len(runes) {
		bufCopy = string(runes[:len(runes)-offsetCopy])
	}

	s.aiMu.Lock()
	if s.aiTimer != nil {
		s.aiTimer.Stop()
	}
	if s.aiCancel != nil {
		s.aiCancel()
		s.aiCancel = nil
	}
	if config.Get().AI.Enabled && bufCopy != "" && !navCopy && offsetCopy == 0 {
		queryTarget := bufCopy
		debounceMS := config.Get().AI.DebounceMS
		if debounceMS <= 0 {
			debounceMS = 500
		}
		s.aiTimer = time.AfterFunc(time.Duration(debounceMS)*time.Millisecond, func() {
			if len(strings.TrimSpace(queryTarget)) < 3 {
				return
			}
			s.aiMu.Lock()
			ctx, cancel := context.WithCancel(context.Background())
			s.aiCancel = cancel
			s.aiMu.Unlock()
			defer cancel()

			cwd := spec.GetCWD()
			var recentCmds []string
			var lastCmd string
			if hist, err := integration.SearchHistory("", nil); err == nil {
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
			if s.overlay.InjectAISuggestion(*sugg) {
				s.renderOverlay()
			}
		})
	}
	s.aiMu.Unlock()

	var b strings.Builder
	if !navCopy {
		if bufCopy == "" && !s.overlay.IsVisible() {
			writeStdout([]byte(s.overlay.ClearAndDisable()))
			return
		}
		logger.Debugf("Render query: '%s', mode: %s", bufCopy, modeCopy)
		results := MergeResults(bufCopy, modeCopy)
		logger.Debugf("Render results found: %d", len(results))

		if len(results) == 0 || (len(results) == 1 && strings.TrimSpace(results[0].Cmd) == strings.TrimSpace(bufCopy) && !strings.HasSuffix(bufCopy, " ")) {
			b.WriteString(s.overlay.HideMenu(bufCopy))
			writeStdout([]byte(b.String()))
			return
		}

		if s.overlay.IsVisible() {
			b.WriteString(s.overlay.Clear())
		}
		s.overlay.SetQueryAndItems(bufCopy, results)
	} else {
		if s.overlay.IsVisible() {
			b.WriteString(s.overlay.Clear())
		}
	}

	s.overlay.SetUserNavigated(navCopy)
	if !s.disableGhostText.Load() {
		b.WriteString(s.overlay.RenderGhostText(bufCopy, navCopy, offsetCopy == 0))
	}
	currentCmd := s.overlay.GetCurrentCmd()
	logger.Debugf("RenderOverlay nav: %v, typedQuery: '%s', currentCmd: '%s'", navCopy, s.overlay.GetTypedQuery(), currentCmd)
	b.WriteString(s.overlay.Render())
	writeStdout([]byte(b.String()))
}

func (s *WrapperSession) renderOverlay() {
	s.renderMu.Lock()
	defer s.renderMu.Unlock()

	if !s.suggestionsEnabled.Load() || s.isExecuting() {
		if s.renderTimer != nil {
			s.renderTimer.Stop()
			s.renderTimer = nil
		}
		return
	}

	if s.userNavigated.Load() {
		return
	}

	if s.renderTimer != nil {
		s.renderTimer.Stop()
	}
	s.renderTimer = time.AfterFunc(25*time.Millisecond, func() {
		s.renderMu.Lock()
		s.renderTimer = nil
		s.renderMu.Unlock()
		s.renderMenuNow()
	})
}

func (s *WrapperSession) handleInputLoop() {
	for {
		inputSlice := make([]byte, 128)
		n, err := os.Stdin.Read(inputSlice)
		if err != nil {
			break
		}

		if n > 0 {
			if s.isExecuting() {
				s.inBracketedPaste = false
				_, _ = s.ptmx.Write(inputSlice[:n])
				continue
			}

			logger.Debugf("Stdin raw input: bytes=%q, hex=%x", inputSlice[:n], inputSlice[:n])

			shouldOverlayDraw := false
			for i := 0; i < n; i++ {
				b := inputSlice[i]
				intercepted := false

				if b == '\033' {
					if i+5 < n && inputSlice[i+1] == '[' && inputSlice[i+2] == '2' && inputSlice[i+3] == '0' {
						if (inputSlice[i+4] == '0' || inputSlice[i+4] == '1') && inputSlice[i+5] == '~' {
							intercepted = true
							s.inBracketedPaste = inputSlice[i+4] == '0'
							logger.Debugf("Intercepted bracketed paste event inPaste=%v", s.inBracketedPaste)
							_, _ = s.ptmx.Write(inputSlice[i : i+6])
							i += 5
							continue
						}
					}
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						if inputSlice[i+1] == '[' && inputSlice[i+2] == 'Z' {
							intercepted = true
							s.suggestionsEnabled.Store(!s.suggestionsEnabled.Load())
							logger.Debugf("Intercepted Shift+Tab, suggestionsEnabled=%v", s.suggestionsEnabled.Load())
							if !s.suggestionsEnabled.Load() {
								writeStdout([]byte(s.overlay.ClearAndDisable()))
							} else {
								shouldOverlayDraw = true
							}
							i += 2
							continue
						}

						if s.overlay.IsVisible() && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') {
							intercepted = true
							s.userNavigated.Store(true)

							arrowDir := "down"
							if inputSlice[i+2] == 'A' {
								arrowDir = "up"
							}
							moved, selectedCmd := s.overlay.MoveCursor(arrowDir)
							if !moved {
								i += 2
								continue
							}

							s.bufferMu.Lock()
							activeModeMu.RLock()
							isHistMode := activeMode == "history"
							activeModeMu.RUnlock()
							var toWrite []byte
							if isHistMode && selectedCmd != "" {
								s.naiveBuffer = selectedCmd
								s.cursorOffset = 0
								toWrite = append([]byte{0x15}, selectedCmd...)
							}
							bufCopy := s.naiveBuffer
							offsetCopy := s.cursorOffset
							s.bufferMu.Unlock()

							if len(toWrite) > 0 {
								_, _ = s.ptmx.Write(toWrite)
							}

							var b strings.Builder
							if !s.disableGhostText.Load() {
								b.WriteString(s.overlay.RenderGhostText(bufCopy, true, offsetCopy == 0))
							}
							b.WriteString(s.overlay.Render())
							writeStdout([]byte(b.String()))

							i += 2
							continue
						} else if !s.overlay.IsVisible() && s.bufferEmpty() && (inputSlice[i+2] == 'A' || inputSlice[i+2] == 'B') {
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
								limit := min(len(results), 100)
								var historyList []spec.Suggestion

								if inputSlice[i+2] == 'A' {
									for j := limit - 1; j >= 0; j-- {
										historyList = append(historyList, results[j])
									}
								} else {
									for j := range limit {
										historyList = append(historyList, results[j])
									}
								}

								selected := s.overlay.SetHistoryList(historyList, inputSlice[i+2] == 'A')
								if selected != "" {
									s.bufferMu.Lock()
									s.naiveBuffer = selected
									s.cursorOffset = 0
									s.bufferMu.Unlock()

									s.userNavigated.Store(true)
									writeStdout([]byte(s.overlay.Render()))
									_, _ = s.ptmx.Write(append([]byte{0x15}, selected...))
								}
							}
							i += 2
							continue
						} else if !s.disableGhostText.Load() && inputSlice[i+2] == 'C' {
							s.bufferMu.Lock()
							atEnd := (s.cursorOffset == 0)
							ghostText := s.overlay.GetGhostText(s.naiveBuffer, atEnd)
							s.bufferMu.Unlock()

							if len(ghostText) > 0 {
								intercepted = true
								logger.Debugf("Intercepted Right Arrow (accepted ghost text: %q)", ghostText)
								s.bufferMu.Lock()
								s.naiveBuffer += ghostText
								s.cursorOffset = 0
								s.bufferMu.Unlock()
								_, _ = s.ptmx.Write([]byte(ghostText))
								shouldOverlayDraw = true
								i += 2
								continue
							}
						}
					}

					isLeftRightArrow := false
					if i+2 < n && (inputSlice[i+1] == '[' || inputSlice[i+1] == 'O') {
						if inputSlice[i+2] == 'D' {
							s.bufferMu.Lock()
							isEmptyQuery := s.naiveBuffer == "" && (!s.overlay.IsVisible() || s.overlay.GetTypedQuery() == "")
							s.bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}
							s.bufferMu.Lock()
							if s.naiveBuffer != "" || s.overlay.IsVisible() {
								s.cursorOffset++
								if s.cursorOffset > len(s.naiveBuffer) {
									s.cursorOffset = len(s.naiveBuffer)
								}
								shouldOverlayDraw = true
								s.userNavigated.Store(false)
							}
							s.bufferMu.Unlock()
							isLeftRightArrow = true
						} else if inputSlice[i+2] == 'C' {
							s.bufferMu.Lock()
							isEmptyQuery := s.naiveBuffer == "" && (!s.overlay.IsVisible() || s.overlay.GetTypedQuery() == "")
							s.bufferMu.Unlock()
							if isEmptyQuery {
								intercepted = true
								i += 2
								continue
							}
							s.bufferMu.Lock()
							if s.naiveBuffer != "" || s.overlay.IsVisible() {
								s.cursorOffset--
								if s.cursorOffset < 0 {
									s.cursorOffset = 0
								}
								shouldOverlayDraw = true
								s.userNavigated.Store(false)
							}
							s.bufferMu.Unlock()
							isLeftRightArrow = true
						}
					}

					if !intercepted {
						writeStdout([]byte(s.overlay.ClearAndDisable()))
						s.disableGhostText.Store(true)
						if !isLeftRightArrow {
							s.bufferMu.Lock()
							s.naiveBuffer = ""
							s.cursorOffset = 0
							s.bufferMu.Unlock()
						}

						_, _ = s.ptmx.Write([]byte{b})
						for j := i + 1; j < n; j++ {
							char := inputSlice[j]
							_, _ = s.ptmx.Write([]byte{char})
							i = j
							if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '~' {
								break
							}
						}
					}
					continue
				}

				if b == 0x12 {
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
					if s.userNavigated.Load() {
						s.bufferMu.Lock()
						s.naiveBuffer = s.overlay.GetTypedQuery()
						s.cursorOffset = 0
						s.bufferMu.Unlock()
						_, _ = s.ptmx.Write(append([]byte{0x15}, s.overlay.GetTypedQuery()...))
					}
					s.userNavigated.Store(false)
					s.overlay.Show()
					shouldOverlayDraw = true
				} else if b == 0x0d || b == 0x0a {
					intercepted = true
					logger.Debugf("Intercepted Enter key, navigated=%v", s.overlay.GetUserNavigated())
					var cmdToSubmit string
					if s.overlay.IsVisible() && s.overlay.GetUserNavigated() {
						selected := s.overlay.GetCurrentCmd()
						if selected != "" {
							cmdToSubmit = selected
							activeModeMu.RLock()
							currentMode := activeMode
							activeModeMu.RUnlock()
							if currentMode == "spec" {
								s_val := strings.TrimSpace(selected)
								if strings.HasSuffix(s_val, "/") || strings.HasSuffix(s_val, "\\") {
									selected = s_val
								} else {
									selected = s_val + " "
								}
							}
							_, _ = s.ptmx.Write(append([]byte{0x15}, selected...))
						}
					}
					writeStdout([]byte(s.overlay.ClearAndDisable()))
					SetCurrentAISuggestion(nil)
					s.renderMu.Lock()
					if s.renderTimer != nil {
						s.renderTimer.Stop()
						s.renderTimer = nil
					}
					s.renderMu.Unlock()

					s.isCommandActive.Store(true)
					_, _ = s.ptmx.Write([]byte{b})
					s.bufferMu.Lock()
					if cmdToSubmit == "" {
						cmdToSubmit = s.naiveBuffer
					}
					s.lastSubmittedCommand = strings.TrimSpace(cmdToSubmit)
					s.naiveBuffer = ""
					s.cursorOffset = 0
					s.bufferMu.Unlock()
					s.disableGhostText.Store(false)
					shouldOverlayDraw = false
					s.userNavigated.Store(false)
					continue
				} else if b == 0x03 || b == 0x15 {
					intercepted = true
					writeStdout([]byte(s.overlay.ClearAndDisable()))
					SetCurrentAISuggestion(nil)
					s.renderMu.Lock()
					if s.renderTimer != nil {
						s.renderTimer.Stop()
						s.renderTimer = nil
					}
					s.renderMu.Unlock()
					s.isCommandActive.Store(false)
					_, _ = s.ptmx.Write([]byte{b})
					s.bufferMu.Lock()
					s.naiveBuffer = ""
					s.cursorOffset = 0
					s.bufferMu.Unlock()
					s.disableGhostText.Store(false)
					shouldOverlayDraw = false
					s.userNavigated.Store(false)
					continue
				} else if b == 0x09 {
					intercepted = true
					logger.Debugf("Intercepted Tab key, visible=%v", s.overlay.IsVisible())
					if !s.overlay.IsVisible() {
						shouldOverlayDraw = true
					} else {
						selected := s.overlay.GetCurrentCmd()
						writeStdout([]byte(s.overlay.ClearAndDisable()))

						activeModeMu.RLock()
						currentMode := activeMode
						activeModeMu.RUnlock()
						if currentMode == "spec" {
							s_val := strings.TrimSpace(selected)
							if strings.HasSuffix(s_val, "/") || strings.HasSuffix(s_val, "\\") {
								selected = s_val
							} else {
								selected = s_val + " "
							}
						}

						s.bufferMu.Lock()
						s.naiveBuffer = selected
						s.cursorOffset = 0
						s.bufferMu.Unlock()

						_, _ = s.ptmx.Write(append([]byte{0x15}, selected...))
						s.overlay.ResetCursor()
						shouldOverlayDraw = true
						s.userNavigated.Store(false)
					}
					continue
				}

				if !intercepted {
					_, _ = s.ptmx.Write([]byte{b})
					switch b {
					case 0x01:
						s.bufferMu.Lock()
						s.cursorOffset = len(s.naiveBuffer)
						if s.naiveBuffer != "" || s.overlay.IsVisible() {
							shouldOverlayDraw = true
						}
						s.bufferMu.Unlock()
						s.userNavigated.Store(false)
					case 0x05:
						s.bufferMu.Lock()
						s.cursorOffset = 0
						if s.naiveBuffer != "" || s.overlay.IsVisible() {
							shouldOverlayDraw = true
						}
						s.bufferMu.Unlock()
						s.userNavigated.Store(false)

					case 127, 0x08:
						s.bufferMu.Lock()
						wasEmpty := len(s.naiveBuffer) == 0
						if !wasEmpty {
							runes := []rune(s.naiveBuffer)
							if s.cursorOffset <= 0 {
								if len(runes) > 0 {
									s.naiveBuffer = string(runes[:len(runes)-1])
								}
								s.cursorOffset = 0
							} else {
								if s.cursorOffset > len(runes) {
									s.cursorOffset = len(runes)
								}
								pos := len(runes) - s.cursorOffset
								if pos > 0 && pos <= len(runes) {
									s.naiveBuffer = string(append(runes[:pos-1], runes[pos:]...))
								}
							}
						}
						isEmptyNow := len(s.naiveBuffer) == 0
						s.bufferMu.Unlock()

						if wasEmpty || isEmptyNow {
							writeStdout([]byte(s.overlay.ClearAndDisable()))
							s.userNavigated.Store(false)
							continue
						}
						shouldOverlayDraw = true
						s.userNavigated.Store(false)
					case 0x17:
						s.bufferMu.Lock()
						wasEmpty := len(s.naiveBuffer) == 0
						trimBuf := strings.TrimRight(s.naiveBuffer, " ")
						lastSpace := strings.LastIndex(trimBuf, " ")
						if lastSpace >= 0 {
							s.naiveBuffer = trimBuf[:lastSpace+1]
						} else {
							s.naiveBuffer = ""
						}
						s.cursorOffset = 0
						isEmptyNow := len(s.naiveBuffer) == 0
						s.bufferMu.Unlock()

						if wasEmpty || isEmptyNow {
							writeStdout([]byte(s.overlay.ClearAndDisable()))
							s.userNavigated.Store(false)
							continue
						}
						shouldOverlayDraw = true
						s.userNavigated.Store(false)
					case 0x0c:
						shouldOverlayDraw = true
						s.userNavigated.Store(false)
					case '\r', '\n', 0x03, 0x15:
						s.inBracketedPaste = false
						s.bufferMu.Lock()
						s.naiveBuffer = ""
						s.cursorOffset = 0
						s.bufferMu.Unlock()
						activeModeMu.Lock()
						activeMode = loadMode()
						activeModeMu.Unlock()
						s.disableGhostText.Store(false)
						writeStdout([]byte(s.overlay.ClearAndDisable()))
						SetCurrentAISuggestion(nil)
						s.userNavigated.Store(false)
					default:
						if b >= 32 && b <= 126 {
							s.bufferMu.Lock()
							isSpaceAlias := !s.inBracketedPaste && b == ' ' && s.naiveBuffer != "" && !strings.Contains(s.naiveBuffer, " ")
							var target string
							var ok bool
							if isSpaceAlias {
								target, ok = spec.GetAlias(s.naiveBuffer)
							}
							s.bufferMu.Unlock()

							if isSpaceAlias && ok {
								_, _ = s.ptmx.Write(append([]byte{0x15}, target+" "...))
								s.bufferMu.Lock()
								s.naiveBuffer = target + " "
								s.cursorOffset = 0
								s.bufferMu.Unlock()
								shouldOverlayDraw = true
								continue
							}
							s.bufferMu.Lock()
							if s.cursorOffset == 0 {
								s.naiveBuffer += string(b)
							} else {
								if s.cursorOffset > len(s.naiveBuffer) {
									s.cursorOffset = len(s.naiveBuffer)
								}
								pos := len(s.naiveBuffer) - s.cursorOffset
								if pos >= 0 && pos <= len(s.naiveBuffer) {
									s.naiveBuffer = s.naiveBuffer[:pos] + string(b) + s.naiveBuffer[pos:]
								} else {
									s.naiveBuffer += string(b)
									s.cursorOffset = 0
								}
							}
							s.bufferMu.Unlock()
							shouldOverlayDraw = true
							s.userNavigated.Store(false)
							s.overlay.SetUserNavigated(false)
						}
					}
				}
			}
			if shouldOverlayDraw {
				s.renderOverlay()
			}
		}
	}
}
func startSignalMonitor(ptmx *os.File, c *exec.Cmd, shellName string) {
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
}
