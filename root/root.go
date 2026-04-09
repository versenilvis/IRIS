package root

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	_ "github.com/versenilvis/iris/commands/dev"
	_ "github.com/versenilvis/iris/commands/fs"
	_ "github.com/versenilvis/iris/commands/info"
	_ "github.com/versenilvis/iris/commands/runner"
	_ "github.com/versenilvis/iris/commands/search"
	_ "github.com/versenilvis/iris/commands/view"
	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/integration"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "iris",
	Short: "IRIS is an awesome cli auto-completion tool",
	Long: `IRIS (a.k.a Intelligent Real-time Input Suggestion) is a shell auto-autocompletion tool.
It works exactly like coding editor suggestion menu drop down.`,
	Run: func(cmd *cobra.Command, args []string) {
		runWrapper()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runWrapper() {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating pipe:", err)
		return
	}

	c := exec.Command("bash")

	// fd 0, 1, 2 are stdin, stdout, stderr (handled by pty)
	// fd 3 is our write pipe
	c.ExtraFiles = []*os.File{w}
	c.Env = append(os.Environ(), "IRIS_FD=3")

	ptmx, err := pty.Start(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting pty:", err)
		return
	}
	defer ptmx.Close()

	core.ShellPID = c.Process.Pid

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

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

				// delete line (ctrl+u)
				ptmx.Write([]byte{0x15})
				ptmx.Write([]byte(selected))
				
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
