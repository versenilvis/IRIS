package spec

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
}

// Option represents a command flag or option
type Option struct {
	Name        string
	Description string
}

// Suggestion represents an item in the suggestion menu
type Suggestion struct {
	Cmd  string
	Desc string
	Icon string
}

var Registry = map[string]*Spec{}

// Register adds a new spec to the global Registry
// example: Register(&Spec{Name: "git"})
func Register(s *Spec) {
	Registry[s.Name] = s
}

// ResetRegistry clears all registered specs - use in tests only
func ResetRegistry() {
	Registry = make(map[string]*Spec)
}
