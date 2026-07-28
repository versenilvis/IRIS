package integration

import (
	"bufio"
	"os"
	"strings"
)

func parseHistoryFile(shellName, histFile string) ([]string, error) {
	file, err := os.Open(histFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if file != nil {
		defer func() { _ = file.Close() }()
	}

	var allCmds []string
	if file != nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			cmd := line

			if shellName == "zsh" {
				parts := strings.SplitN(line, ";", 2)
				if len(parts) == 2 {
					cmd = parts[1]
				}
			} else if shellName == "bash" {
				if strings.HasPrefix(line, "#") && len(line) > 1 {
					isTimestamp := true
					for _, c := range line[1:] {
						if c < '0' || c > '9' {
							isTimestamp = false
							break
						}
					}
					if isTimestamp {
						continue
					}
				}
			} else if shellName == "fish" {
				if after, ok := strings.CutPrefix(line, "- cmd: "); ok {
					cmd = after
				} else {
					continue
				}
			}

			cmd = strings.TrimSpace(cmd)
			if cmd != "" {
				allCmds = append(allCmds, cmd)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	return allCmds, nil
}
