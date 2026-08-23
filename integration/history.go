package integration

import (
	"bufio"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/versenilvis/fuzzy"
	"github.com/versenilvis/iris/integration/shell"
	"github.com/versenilvis/iris/internal/config"
	_ "modernc.org/sqlite"
)

var (
	sessionHistory   []string
	sessionHistoryMu sync.Mutex

	historyCache  []string
	idMapCache    map[string]int
	sourceMapCache map[string]string
	searcherCache *fuzzy.Searcher
	mu            sync.Mutex
	lastModTime   int64

	atuinCmds    []string
	atuinLastMod int64
	lastAtuinMode int = -1
)

func RecordSessionCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
	cmd = sanitizeUTF8(cmd)
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

type HistResult struct {
	ID         int
	Cmd        string
	FuzzyScore int
	Source     string
}

func init() {
	idMapCache = make(map[string]int)
	sourceMapCache = make(map[string]string)
}

func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	var result strings.Builder
	for _, r := range s {
		if r == utf8.RuneError {
			result.WriteRune('�')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func loadAtuinCmds() ([]string, error) {
	dbPath, err := config.AtuinDBPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), `SELECT command FROM history WHERE deleted_at IS NULL ORDER BY timestamp DESC LIMIT 10000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seen := make(map[string]bool)
	var cmds []string
	for rows.Next() {
		var cmd string
		if err := rows.Scan(&cmd); err != nil {
			continue
		}
		cmd = strings.TrimSpace(sanitizeUTF8(cmd))
		cmd = strings.ReplaceAll(cmd, "\n", " ")
		cmd = strings.ReplaceAll(cmd, "\r", "")
		if cmd != "" && !seen[cmd] {
			seen[cmd] = true
			cmds = append(cmds, cmd)
		}
	}
	// reverse array to be oldest-first
	for i, j := 0, len(cmds)-1; i < j; i, j = i+1, j-1 {
		cmds[i], cmds[j] = cmds[j], cmds[i]
	}
	return cmds, rows.Err()
}

func SearchHistory(query string, aliases map[string]string) ([]HistResult, error) {
	mu.Lock()
	defer mu.Unlock()

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	shellName := "bash"
	if shell.Current != nil {
		shellName = shell.Current.GetName()
	}

	var histFile string
	if envHist := os.Getenv("HISTFILE"); envHist != "" {
		histFile = envHist
	} else {
		switch shellName {
		case "zsh":
			histFile = filepath.Join(shell.GetZshConfigDir(), ".zsh_history")
		case "fish":
			histFile = filepath.Join(shell.GetFishDataDir(), "fish_history")
		default:
			histFile = filepath.Join(home, ".bash_history")
		}
	}

	atuinMode := config.Get().Core.Atuin
	if lastAtuinMode != -1 && atuinMode != lastAtuinMode {
		historyCache = nil
	}
	lastAtuinMode = atuinMode

	if atuinMode > 0 {
		dbPath, _ := config.AtuinDBPath()
		if info, err := os.Stat(dbPath); err == nil {
			mod := info.ModTime().UnixNano()
			if mod != atuinLastMod {
				atuinLastMod = mod
				atuinCmds = nil
				historyCache = nil
			}
		}
		if atuinCmds == nil {
			atuinCmds, _ = loadAtuinCmds()
		}
	}

	if info, err := os.Stat(histFile); err == nil {
		if info.ModTime().UnixNano() > lastModTime {
			historyCache = nil // force reload
			idMapCache = make(map[string]int)
			lastModTime = info.ModTime().UnixNano()
		}
	}

	// lazy load history if cache is empty
	if len(historyCache) == 0 {
		file, err := os.Open(histFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if file != nil {
			defer func() { _ = file.Close() }()
		}

		var allCmds []string
		if atuinMode == 1 && len(atuinCmds) > 0 {
			// atuin only — prepend newest-first so the merge loop below works
			allCmds = atuinCmds
		} else {
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
						cmd = sanitizeUTF8(cmd)
						allCmds = append(allCmds, cmd)
					}
				}
				if scanner.Err() != nil {
					_ = scanner.Err()
				}
			}
			// mode 2: append atuin (newer, higher priority) after shell file
			if atuinMode == 2 && len(atuinCmds) > 0 {
				allCmds = append(allCmds, atuinCmds...)
			}
		}

		// build historyCache backwards so newest commands come first
		seen := make(map[string]bool)
		atuinSeen := make(map[string]bool)
		for _, c := range atuinCmds {
			atuinSeen[c] = true
		}

		historyCache = nil
		idMapCache = make(map[string]int)
		sourceMapCache = make(map[string]string)

		currentID := len(sessionHistory) + len(allCmds)

		sessionHistoryMu.Lock()
		for i := len(sessionHistory) - 1; i >= 0; i-- {
			cmd := sessionHistory[i]
			if !seen[cmd] {
				historyCache = append(historyCache, cmd)
				seen[cmd] = true
				idMapCache[cmd] = currentID
				if atuinSeen[cmd] {
					sourceMapCache[cmd] = "atuin"
				} else {
					sourceMapCache[cmd] = "session"
				}
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
				if atuinSeen[cmd] {
					sourceMapCache[cmd] = "atuin"
				} else {
					sourceMapCache[cmd] = "history"
				}
				currentID--
			}
		}

		searcherCache = fuzzy.NewPlainSearcher(historyCache)
	}

	if query == "" {
		var results []HistResult
		limit := min(len(historyCache), 100)

		for i := range limit {
			cmd := historyCache[i]
			results = append(results, HistResult{
				ID:  idMapCache[cmd],
				Cmd: cmd,
				Source: sourceMapCache[cmd],
			})
		}
		return results, nil
	}

	var alternativeQueries []string
	for name, target := range aliases {
		if target != "" {
			qLow := strings.ToLower(query)
			tLow := strings.ToLower(target)
			nLow := strings.ToLower(name)

			if qLow == tLow {
				alternativeQueries = append(alternativeQueries, name)
			} else if strings.HasPrefix(qLow, tLow+" ") {
				suffix := query[len(target):]
				alternativeQueries = append(alternativeQueries, name+suffix)
			}

			if qLow == nLow {
				alternativeQueries = append(alternativeQueries, target)
			} else if strings.HasPrefix(qLow, nLow+" ") {
				suffix := query[len(name):]
				alternativeQueries = append(alternativeQueries, target+suffix)
			}
		}
	}

	var results []HistResult
	seenCmds := make(map[string]bool)

	addMatches := func(q string) {
		qLow := strings.ToLower(q)

		// extract pure substring matches (all words present) based strictly on recency order (historyCache is newest-first)
		// this ensures that long commands with exact substrings are never truncated by the fuzzy searcher's limit
		strictMatches := 0
		words := strings.Fields(qLow)
		if len(words) == 0 {
			words = []string{qLow}
		}

		for _, cmd := range historyCache {
			if seenCmds[cmd] {
				continue
			}

			cmdLow := strings.ToLower(cmd)
			matchAll := true
			for _, w := range words {
				if !strings.Contains(cmdLow, w) {
					matchAll = false
					break
				}
			}

			if !matchAll {
				continue
			}

			seenCmds[cmd] = true
			results = append(results, HistResult{
				ID:         idMapCache[cmd],
				Cmd:        cmd,
				FuzzyScore: 10000,
				Source:     sourceMapCache[cmd],
			})
			strictMatches++
			if strictMatches >= 200 {
				break
			}
		}

		matches := searcherCache.SearchWithScores(q, &fuzzy.SearchOptions{Limit: 1000})
		for _, m := range matches {
			if seenCmds[m.Str] {
				continue
			}

			// filter out extremely weak fuzzy matches (e.g. random garbage typing that
			// loosely matches across a very long command)
			if len(q) > 0 && m.Score/len(q) < 150 {
				continue
			}

			seenCmds[m.Str] = true
			results = append(results, HistResult{
				ID:         idMapCache[m.Str],
				Cmd:        m.Str,
				FuzzyScore: m.Score,
				Source:     sourceMapCache[m.Str],
			})
		}
	}

	addMatches(query)
	for _, altQ := range alternativeQueries {
		addMatches(altQ)
	}

	getTier := func(cmd, q string) int {
		bestTier := 4
		check := func(ql string) {
			cmdLow := strings.ToLower(cmd)
			qlLow := strings.ToLower(ql)
			tier := 4
			if cmdLow == qlLow {
				tier = 1
			} else if strings.HasPrefix(cmdLow, qlLow) {
				tier = 2
			} else if strings.Contains(cmdLow, qlLow) {
				tier = 3
			}
			if tier < bestTier {
				bestTier = tier
			}
		}
		check(q)
		for _, altQ := range alternativeQueries {
			check(altQ)
		}
		return bestTier
	}

	tiers := make([]int, len(results))
	for i, r := range results {
		tiers[i] = getTier(r.Cmd, query)
	}

	sort.SliceStable(results, func(i, j int) bool {
		tI := tiers[i]
		tJ := tiers[j]
		if tI != tJ {
			return tI < tJ
		}

		if tI == 4 && results[i].FuzzyScore != results[j].FuzzyScore {
			return results[i].FuzzyScore > results[j].FuzzyScore
		}

		return results[i].ID > results[j].ID
	})

	return results, nil
}
