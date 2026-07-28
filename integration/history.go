package integration

import (
	"sort"
	"strings"

	"github.com/versenilvis/fuzzy"
)

type HistResult struct {
	ID         int
	Cmd        string
	FuzzyScore int
}

func SearchHistory(query string, aliases map[string]string) ([]HistResult, error) {
	mu.Lock()
	defer mu.Unlock()

	if err := ensureCacheLoaded(); err != nil {
		return nil, err
	}

	if query == "" {
		var results []HistResult
		limit := min(len(historyCache), 100)
		for i := range limit {
			cmd := historyCache[i]
			results = append(results, HistResult{
				ID:  idMapCache[cmd],
				Cmd: cmd,
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

	addMatches := func(q string, subcmdFilter bool) {
		qLow := strings.ToLower(q)
		queryFirstWord := ""
		querySecondWord := ""
		if strings.IndexByte(qLow, ' ') != -1 {
			if fields := strings.Fields(qLow); len(fields) > 0 {
				queryFirstWord = fields[0]
				for _, f := range fields[1:] {
					if !strings.HasPrefix(f, "-") {
						querySecondWord = f
						break
					}
				}
			}
		}

		strictMatches := 0
		for _, cmd := range historyCache {
			if seenCmds[cmd] {
				continue
			}
			fields := strings.Fields(cmd)
			firstWordLow := ""
			if len(fields) > 0 {
				firstWordLow = strings.ToLower(fields[0])
			}

			if queryFirstWord != "" {
				if firstWordLow != queryFirstWord {
					continue
				}
				if subcmdFilter && querySecondWord != "" {
					if len(fields) < 2 {
						continue
					}
					secondWordLow := strings.ToLower(fields[1])
					if !strings.HasPrefix(secondWordLow, querySecondWord) {
						continue
					}
				}
			} else {
				if !strings.HasPrefix(firstWordLow, qLow) {
					continue
				}
			}

			seenCmds[cmd] = true
			results = append(results, HistResult{
				ID:         idMapCache[cmd],
				Cmd:        cmd,
				FuzzyScore: 10000,
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

			fields := strings.Fields(m.Str)
			firstWord := m.Str
			if len(fields) > 0 {
				firstWord = fields[0]
			}
			firstWordLow := strings.ToLower(firstWord)

			if queryFirstWord != "" {
				if firstWordLow != queryFirstWord {
					continue
				}
				if subcmdFilter && querySecondWord != "" {
					if len(fields) < 2 {
						continue
					}
					secondWordLow := strings.ToLower(fields[1])
					if !strings.HasPrefix(secondWordLow, querySecondWord) {
						continue
					}
				}
			} else {
				if !strings.HasPrefix(firstWordLow, qLow) {
					continue
				}
			}

			seenCmds[m.Str] = true
			results = append(results, HistResult{
				ID:         idMapCache[m.Str],
				Cmd:        m.Str,
				FuzzyScore: m.Score,
			})
		}
	}

	addMatches(query, true)
	for _, altQ := range alternativeQueries {
		addMatches(altQ, true)
	}

	if len(results) == 0 {
		addMatches(query, false)
		for _, altQ := range alternativeQueries {
			addMatches(altQ, false)
		}
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
