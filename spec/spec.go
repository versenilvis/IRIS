package spec

import "slices"

type GeneratorFunc func(tokens []string, prefix string, partial string) []Suggestion

// Spec defines a top-level command structure
type Spec struct {
	Name        string
	Aliases     []string
	Description string
	Icon        string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
	MaxArgs     int
}

// Subcommand defines nested command logic
type Subcommand struct {
	Name        string
	Aliases     []string
	Description string
	Icon        string
	Subcommands []Subcommand
	Options     []Option
	Generator   GeneratorFunc
	MaxArgs     int
	Priority    int
}

// Option represents a command flag or option
type Option struct {
	Name        string
	Description string
	Priority    int
}

// Suggestion represents an item in the suggestion menu
type Suggestion struct {
	Cmd        string
	Desc       string
	Icon       string
	Source     string // "history", "spec", "ai"
	Confidence int    // 0-100
	Priority   int    // static author priority
}

var Registry = map[string]*Spec{}

// duplicateNames records specs registered under a name already taken.
//
// Registration is last-one-wins, and every spec is registered from an init()
// spread across the commands packages, so a name defined twice resolves on
// package init order -- and does so silently. `find` was registered in two
// packages, and the one that lost had the file generator, which is exactly the
// kind of difference nobody notices from reading either file.
var duplicateNames []string

// Register adds a new spec to the global Registry
// example: Register(&Spec{Name: "git"})
func Register(s *Spec) {
	if _, exists := Registry[s.Name]; exists {
		duplicateNames = append(duplicateNames, s.Name)
	}
	Registry[s.Name] = s
}

// DuplicateNames returns the spec names that were registered more than once.
// A non-empty result is a bug: which definition survives is decided by package
// init order rather than by intent.
func DuplicateNames() []string {
	return slices.Clone(duplicateNames)
}

// ResetRegistry clears all registered specs - use in tests only
func ResetRegistry() {
	Registry = make(map[string]*Spec)
	duplicateNames = nil
}
