package scoring

import (
	"testing"

	"github.com/versenilvis/iris/spec"
)

func TestExtractSkeleton(t *testing.T) {
	spec.ResetRegistry()
	spec.Register(&spec.Spec{
		Name: "git",
		Subcommands: []spec.Subcommand{
			{Name: "checkout"},
			{Name: "push"},
		},
	})

	tests := []struct {
		name string
		buf  string
		want string
	}{
		{
			name: "spec command with args",
			buf:  "git checkout feature-x",
			want: "git checkout",
		},
		{
			name: "fallback no spec binary only",
			buf:  "cargo build --release",
			want: "cargo",
		},
		{
			name: "fallback single command",
			buf:  "ls -la",
			want: "ls",
		},
		{
			name: "empty input",
			buf:  "   ",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractSkeleton(tc.buf)
			if got != tc.want {
				t.Errorf("ExtractSkeleton(%q) = %q, want %q", tc.buf, got, tc.want)
			}
		})
	}
}
