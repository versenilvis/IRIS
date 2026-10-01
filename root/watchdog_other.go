//go:build !linux

package root

import (
	"os/exec"
)

func setWatchdogSysProcAttr(cmd *exec.Cmd) {}
