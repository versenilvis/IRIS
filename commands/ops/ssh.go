package ops

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/versenilvis/iris/commands/core"
)

func sshHostGenerator(tokens []string, _ string, _ string) []core.Suggestion {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	paths := []string{
		filepath.Join(home, ".ssh", "config"),
		"/etc/ssh/ssh_config",
	}

	seen := make(map[string]bool)
	var results []core.Suggestion

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(strings.ToLower(line), "host ") {
				continue
			}

			// "Host foo bar baz" → each token is a separate alias
			parts := strings.Fields(line)
			for _, host := range parts[1:] {
				// skip wildcards
				if strings.ContainsAny(host, "*?!") {
					continue
				}
				if seen[host] {
					continue
				}
				seen[host] = true
				results = append(results, core.Suggestion{Cmd: host, Desc: "ssh host"})
			}
		}
		_ = f.Close()
	}

	return results
}

func init() {
	sshOptions := []core.Option{
		{Name: "-p", Description: "port"},
		{Name: "-i", Description: "identity file"},
		{Name: "-L", Description: "local port forward"},
		{Name: "-R", Description: "remote port forward"},
		{Name: "-D", Description: "dynamic socks proxy"},
		{Name: "-N", Description: "no remote command"},
		{Name: "-f", Description: "background"},
		{Name: "-v", Description: "verbose"},
		{Name: "-A", Description: "agent forwarding"},
		{Name: "-X", Description: "x11 forwarding"},
	}

	core.Register(&core.Spec{
		Name:        "ssh",
		Description: "secure shell",
		Generator:   sshHostGenerator,
		Options:     sshOptions,
	})

	core.Register(&core.Spec{
		Name:        "scp",
		Description: "secure copy",
		Generator:   sshHostGenerator,
		Options: []core.Option{
			{Name: "-r", Description: "recursive"},
			{Name: "-p", Description: "port"},
			{Name: "-i", Description: "identity file"},
			{Name: "-P", Description: "remote port"},
		},
	})

	core.Register(&core.Spec{
		Name:        "rsync",
		Description: "remote sync",
		Generator:   sshHostGenerator,
		Options: []core.Option{
			{Name: "-av", Description: "archive + verbose"},
			{Name: "-z", Description: "compress"},
			{Name: "--delete", Description: "delete extraneous"},
			{Name: "--progress", Description: "show progress"},
			{Name: "-e", Description: "remote shell command"},
			{Name: "--exclude", Description: "exclude pattern"},
			{Name: "--dry-run", Description: "simulation only"},
		},
	})
}
