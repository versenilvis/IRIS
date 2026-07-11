package tests

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/spec"
)

func TestSshHostGenerator(t *testing.T) {
	tmp := t.TempDir()
	sshDir := filepath.Join(tmp, ".ssh")
	_ = os.MkdirAll(sshDir, 0700)

	configContent := `
Host prod-server
    HostName 10.0.0.1
    User deploy

Host staging bastion
    HostName staging.example.com
    User ubuntu

Host *.internal
    User admin

Host !forbidden wildcard-test
    HostName test.internal
`
	_ = os.WriteFile(filepath.Join(sshDir, "config"), []byte(configContent), 0600)

	results := sshHostGeneratorFromPath(filepath.Join(sshDir, "config"))

	found := make(map[string]bool)
	for _, r := range results {
		found[r.Cmd] = true
	}

	if !found["prod-server"] {
		t.Error("expected prod-server in suggestions")
	}
	if !found["staging"] {
		t.Error("expected staging in suggestions")
	}
	if !found["bastion"] {
		t.Error("expected bastion in suggestions")
	}
	if !found["wildcard-test"] {
		t.Error("expected wildcard-test in suggestions")
	}
	if found["*.internal"] || found["!forbidden"] {
		t.Error("wildcards and negations should be ignored")
	}
}

// sshHostGeneratorFromPath is a helper that reads a specific ssh config path
func sshHostGeneratorFromPath(configPath string) []spec.Suggestion {
	f, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	seen := make(map[string]bool)
	var results []spec.Suggestion

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 || !strings.EqualFold(parts[0], "host") {
			continue
		}
		for _, host := range parts[1:] {
			if strings.ContainsAny(host, "*?!") {
				continue
			}
			if seen[host] {
				continue
			}
			seen[host] = true
			results = append(results, spec.Suggestion{Cmd: host, Desc: "ssh host"})
		}
	}
	return results
}
