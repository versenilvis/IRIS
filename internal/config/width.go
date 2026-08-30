package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Width is a width setting written either as a column count (80) or as a share
// of the terminal ("80%"). Percentages only mean something once the terminal
// size is known, so both forms are carried through to draw time and resolved
// there.
type Width struct {
	n       int
	percent bool
}

// Resolve returns the column count this setting asks for on a terminal of the
// given width. Zero means unset, and is left for the caller to default.
func (w Width) Resolve(term int) int {
	if !w.percent {
		return w.n
	}
	if term <= 0 {
		return 0
	}
	return term * w.n / 100
}

func (w Width) String() string {
	if w.percent {
		return strconv.Itoa(w.n) + "%"
	}
	return strconv.Itoa(w.n)
}

func (w *Width) UnmarshalTOML(v any) error {
	switch t := v.(type) {
	case int64:
		*w = Width{n: int(t)}
		return nil
	case string:
		body, percent := strings.CutSuffix(strings.TrimSpace(t), "%")
		n, err := strconv.Atoi(strings.TrimSpace(body))
		if err != nil {
			return fmt.Errorf("invalid value %q (want a column count like 80, or a share like \"80%%\")", t)
		}
		*w = Width{n: n, percent: percent}
		return nil
	}
	return fmt.Errorf("invalid value of type %T (want a column count like 80, or a share like \"80%%\")", v)
}

func (w Width) MarshalTOML() ([]byte, error) {
	if w.percent {
		return []byte(strconv.Quote(w.String())), nil
	}
	return []byte(w.String()), nil
}

func (w Width) validate() error {
	if w.percent {
		if w.n < 1 || w.n > 100 {
			return fmt.Errorf("must be between 1%% and 100%%")
		}
		return nil
	}
	if w.n < 0 {
		return fmt.Errorf("must not be negative")
	}
	return nil
}
