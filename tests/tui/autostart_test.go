package tui

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/creack/pty"
)

func TestAutostartDoesNotSwallowCommand(t *testing.T) {
	bin := binary(t)

	t.Run("bash", func(t *testing.T) {
		if _, err := exec.LookPath("bash"); err != nil {
			t.Skip("bash not installed")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		initScript, err := exec.CommandContext(ctx, bin, "init", "bash").Output()
		if err != nil {
			t.Fatalf("iris init bash: %v", err)
		}
		cmd := exec.CommandContext(ctx, "bash", "--norc", "-i", "-c", string(initScript)+"\necho RAN_BASH_42")
		ptmx, err := pty.Start(cmd)
		if err != nil {
			t.Fatalf("starting bash in pty: %v", err)
		}
		defer func() { _ = ptmx.Close() }()

		var out bytes.Buffer
		done := make(chan error, 1)
		go func() {
			_, errCopy := io.Copy(&out, ptmx)
			done <- errCopy
		}()

		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			t.Fatalf("timed out waiting for bash command, output:\n%s", out.String())
		case <-done:
		}

		if !bytes.Contains(out.Bytes(), []byte("RAN_BASH_42")) {
			t.Fatalf("expected output to contain RAN_BASH_42, got:\n%s", out.String())
		}
	})

	t.Run("zsh", func(t *testing.T) {
		if _, err := exec.LookPath("zsh"); err != nil {
			t.Skip("zsh not installed")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		initScript, err := exec.CommandContext(ctx, bin, "init", "zsh").Output()
		if err != nil {
			t.Fatalf("iris init zsh: %v", err)
		}
		cmd := exec.CommandContext(ctx, "zsh", "-f", "-i", "-c", string(initScript)+"\necho RAN_ZSH_42")
		ptmx, err := pty.Start(cmd)
		if err != nil {
			t.Fatalf("starting zsh in pty: %v", err)
		}
		defer func() { _ = ptmx.Close() }()

		var out bytes.Buffer
		done := make(chan error, 1)
		go func() {
			_, errCopy := io.Copy(&out, ptmx)
			done <- errCopy
		}()

		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			t.Fatalf("timed out waiting for zsh command, output:\n%s", out.String())
		case <-done:
		}

		if !bytes.Contains(out.Bytes(), []byte("RAN_ZSH_42")) {
			t.Fatalf("expected output to contain RAN_ZSH_42, got:\n%s", out.String())
		}
	})
}
