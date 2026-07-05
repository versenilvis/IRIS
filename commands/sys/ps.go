package sys

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/spec"
)

func processGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// output: "pid command"
	cmd := exec.CommandContext(ctx, "ps", "-eo", "pid,comm")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		pid := parts[0]
		if pid[0] < '0' || pid[0] > '9' {
			continue
		}
		name := parts[1]

		// suggest both PID and process name
		if !seen[pid] {
			seen[pid] = true
			results = append(results, spec.Suggestion{Cmd: pid, Desc: name})
		}
	}
	return results
}

func init() {
	spec.Register(&spec.Spec{
		Name:        "ps",
		Description: "report processes",
		Options: []spec.Option{
			{Name: "-e", Description: "all processes"},
			{Name: "-f", Description: "full format"},
			{Name: "-u", Description: "by user"},
			{Name: "aux", Description: "all + user + x"},
			{Name: "-o", Description: "custom format"},
			{Name: "--pid", Description: "filter by PID"},
			{Name: "--sort", Description: "sort by field"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "kill",
		Description: "send signal to process",
		Generator:   processGenerator,
		Options: []spec.Option{
			{Name: "-9", Description: "SIGKILL (force)"},
			{Name: "-15", Description: "SIGTERM (graceful)"},
			{Name: "-2", Description: "SIGINT"},
			{Name: "-HUP", Description: "SIGHUP (reload)"},
			{Name: "-l", Description: "list signal names"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "killall",
		Description: "kill by process name",
		Generator:   processGenerator,
		Options: []spec.Option{
			{Name: "-9", Description: "SIGKILL"},
			{Name: "-s", Description: "specify signal"},
			{Name: "-u", Description: "only for user"},
			{Name: "-i", Description: "confirm each"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "pkill",
		Description: "kill by pattern",
		Options: []spec.Option{
			{Name: "-9", Description: "SIGKILL"},
			{Name: "-f", Description: "match full command"},
			{Name: "-u", Description: "match by user"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "pgrep",
		Description: "find process by pattern",
		Options: []spec.Option{
			{Name: "-l", Description: "list name and PID"},
			{Name: "-f", Description: "match full command"},
			{Name: "-u", Description: "match by user"},
			{Name: "-a", Description: "list command line"},
		},
	})
}
