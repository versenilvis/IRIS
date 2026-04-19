package root

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type procInfo struct {
	pid  int
	ppid int
	comm string
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
