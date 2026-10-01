package ctxcheck

import (
	"strings"

	"github.com/versenilvis/iris/internal/scoring"
)

type Verdict int

const (
	Free Verdict = iota
	Valid
	Invalid
	Unknown
)

func (v Verdict) String() string {
	switch v {
	case Free:
		return "Free"
	case Valid:
		return "Valid"
	case Invalid:
		return "Invalid"
	case Unknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

type Dialect int

const (
	Posix Dialect = iota
	Fish
	UnknownDialect
)

func DialectFromShell(sh string) (Dialect, bool) {
	switch strings.ToLower(sh) {
	case "bash", "zsh":
		return Posix, true
	case "fish":
		return Fish, true
	default:
		return UnknownDialect, false
	}
}

func Allow(c scoring.Candidate, v Verdict, scopeOf func() int) bool {
	switch v {
	case Invalid:
		return false
	case Unknown:
		return c.Tier == 4
	case Valid:
		return true
	default:
		return c.Tier > 0 || (scopeOf != nil && scopeOf() >= scoring.GlobalScopeThreshold)
	}
}
