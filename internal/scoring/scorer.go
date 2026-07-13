package scoring

import (
	"math"
	"sort"
	"strings"

	"github.com/versenilvis/iris/spec"
)

type ScoreBreakdown struct {
	BasePriority int
	ContextBonus int
	Frecency     int
	MatchQuality int
}

type ScoredSuggestion struct {
	spec.Suggestion
	Score     float64
	Breakdown ScoreBreakdown
}

type ScoreConfig struct {
	WeightBasePriority float64
	WeightContextBonus float64
	WeightFrecency     float64
	WeightMatchQuality float64
}

var DefaultScoreConfig = ScoreConfig{
	WeightBasePriority: 0.30,
	WeightContextBonus: 0.25,
	WeightFrecency:     0.25,
	WeightMatchQuality: 0.20,
}

func Score(suggestions []spec.Suggestion, signals SignalSet) []ScoredSuggestion {
	return ScoreWithConfig(suggestions, signals, DefaultScoreConfig)
}

func ScoreWithConfig(suggestions []spec.Suggestion, signals SignalSet, config ScoreConfig) []ScoredSuggestion {
	if len(suggestions) == 0 {
		return nil
	}

	localMap := make(map[string]float64, len(signals.LocalFrecency))
	for _, e := range signals.LocalFrecency {
		localMap[e.Cmd] = e.RawScore
	}
	globalMap := make(map[string]float64, len(signals.GlobalFrecency))
	for _, e := range signals.GlobalFrecency {
		globalMap[e.Cmd] = e.RawScore
	}

	rawFrec := make([]float64, len(suggestions))
	for i, s := range suggestions {
		if score, ok := localMap[s.Cmd]; ok {
			rawFrec[i] = score
		} else if score, ok := globalMap[s.Cmd]; ok {
			rawFrec[i] = score * 0.7
		} else {
			rawFrec[i] = 0
		}
	}

	normFrec := normalizeFrecency(rawFrec)

	scored := make([]ScoredSuggestion, len(suggestions))
	for i, s := range suggestions {
		bp := basePriorityFor(s)
		cb := ApplyContextRules(signals.Workspace, s.Cmd)
		frec := normFrec[i]
		mq := matchQualityScore(s.Cmd, signals.Query)

		total := config.WeightBasePriority*float64(bp) +
			config.WeightContextBonus*float64(cb) +
			config.WeightFrecency*float64(frec) +
			config.WeightMatchQuality*float64(mq)

		scored[i] = ScoredSuggestion{
			Suggestion: s,
			Score:      total,
			Breakdown: ScoreBreakdown{
				BasePriority: bp,
				ContextBonus: cb,
				Frecency:     frec,
				MatchQuality: mq,
			},
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if math.Abs(scored[i].Score-scored[j].Score) > 1e-6 {
			return scored[i].Score > scored[j].Score
		}
		if scored[i].Breakdown.Frecency != scored[j].Breakdown.Frecency {
			return scored[i].Breakdown.Frecency > scored[j].Breakdown.Frecency
		}
		if scored[i].Breakdown.ContextBonus != scored[j].Breakdown.ContextBonus {
			return scored[i].Breakdown.ContextBonus > scored[j].Breakdown.ContextBonus
		}
		return scored[i].Cmd < scored[j].Cmd
	})

	return scored
}

func basePriorityFor(s spec.Suggestion) int {
	if s.Priority > 0 {
		if s.Priority > 100 {
			return 100
		}
		return s.Priority
	}

	switch s.Source {
	case "spec":
		return 60
	case "ai":
		if s.Confidence > 0 {
			if s.Confidence > 100 {
				return 100
			}
			return s.Confidence
		}
		return 50
	case "history":
		if s.Confidence > 0 {
			if s.Confidence > 100 {
				return 100
			}
			return s.Confidence
		}
		return 40
	default:
		return 50
	}
}

func matchQualityScore(cmd, query string) int {
	cmd = strings.TrimSpace(cmd)
	query = strings.TrimSpace(query)
	if query == "" {
		return 100
	}
	if cmd == query {
		return 100
	}
	if strings.HasPrefix(cmd, query) {
		return 100
	}
	if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(query)) {
		return 80
	}
	if strings.Contains(strings.ToLower(cmd), strings.ToLower(query)) {
		return 50
	}
	if isSubsequence(strings.ToLower(query), strings.ToLower(cmd)) {
		return 30
	}
	return 0
}

func isSubsequence(sub, full string) bool {
	if len(sub) == 0 {
		return true
	}
	i := 0
	for j := 0; j < len(full) && i < len(sub); j++ {
		if sub[i] == full[j] {
			i++
		}
	}
	return i == len(sub)
}

func normalizeFrecency(raw []float64) []int {
	if len(raw) == 0 {
		return nil
	}
	maxRaw := 0.0
	for _, r := range raw {
		if r > maxRaw {
			maxRaw = r
		}
	}
	if maxRaw <= 0 {
		res := make([]int, len(raw))
		return res
	}

	res := make([]int, len(raw))
	for i, r := range raw {
		val := int(math.Round((r / maxRaw) * 100.0))
		if val > 100 {
			val = 100
		} else if val < 0 {
			val = 0
		}
		res[i] = val
	}
	return res
}
