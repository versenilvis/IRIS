package scoring

import (
	"strings"

	"github.com/versenilvis/iris/spec"
)

// ExtractSkeleton extracts the subcommand skeleton of a command string for transition tracking.
// If a spec is registered, it walks the subcommand tree.
// If no spec is registered (or fallback occurs), it returns the first token (binary name).
// Never returns empty string unless input has no tokens.
func ExtractSkeleton(buf string) string {
	if skeleton, ok := spec.TryExtractSkeleton(buf); ok && skeleton != "" {
		return skeleton
	}

	fields := strings.Fields(buf)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
