//go:build linux

package root

import (
	"os/exec"
	"syscall"
)

func setWatchdogSysProcAttr(cmd *exec.Cmd) {
	// kernel kills child if watchdog parent exits abruptly
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}
}
