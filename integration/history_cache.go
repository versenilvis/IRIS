package integration

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/versenilvis/fuzzy"
	"github.com/versenilvis/iris/integration/shell"
)

var (
	sessionHistory   []string
	sessionHistoryMu sync.Mutex

	historyCache  []string
	idMapCache    map[string]int
	searcherCache *fuzzy.Searcher
	mu            sync.Mutex
	lastModTime   int64
)

func RecordSessionCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	sessionHistoryMu.Lock()
	defer sessionHistoryMu.Unlock()

	if len(sessionHistory) > 0 && sessionHistory[len(sessionHistory)-1] == cmd {
		return
	}
	sessionHistory = append(sessionHistory, cmd)
	historyCache = nil // invalidate to merge session history on next search
}

func init() {
	idMapCache = make(map[string]int)
}

func ensureCacheLoaded() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	shellName := "bash"
	if shell.Current != nil {
		shellName = shell.Current.GetName()
	}

	var histFile string
	switch shellName {
	case "zsh":
		histFile = filepath.Join(home, ".zsh_history")
	case "fish":
		histFile = filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		histFile = filepath.Join(home, ".bash_history")
	}

	if info, err := os.Stat(histFile); err == nil {
		if info.ModTime().UnixNano() > lastModTime {
			historyCache = nil // force reload
			idMapCache = make(map[string]int)
			lastModTime = info.ModTime().UnixNano()
		}
	}

	if len(historyCache) == 0 {
		allCmds, err := parseHistoryFile(shellName, histFile)
		if err != nil {
			return err
		}

		seen := make(map[string]bool)
		historyCache = nil
		idMapCache = make(map[string]int)

		currentID := len(sessionHistory) + len(allCmds)

		sessionHistoryMu.Lock()
		for i := len(sessionHistory) - 1; i >= 0; i-- {
			cmd := sessionHistory[i]
			if !seen[cmd] {
				historyCache = append(historyCache, cmd)
				seen[cmd] = true
				idMapCache[cmd] = currentID
				currentID--
			}
		}
		sessionHistoryMu.Unlock()

		for i := len(allCmds) - 1; i >= 0; i-- {
			cmd := allCmds[i]
			if !seen[cmd] {
				historyCache = append(historyCache, cmd)
				seen[cmd] = true
				idMapCache[cmd] = currentID
				currentID--
			}
		}

		searcherCache = fuzzy.NewPlainSearcher(historyCache)
	}
	return nil
}
