package spec

import (
	"testing"
)

func TestTryExtractSkeleton(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name: "git",
		Subcommands: []Subcommand{
			{
				Name:    "checkout",
				Aliases: []string{"co"},
			},
			{
				Name: "remote",
				Subcommands: []Subcommand{
					{Name: "add"},
				},
			},
		},
	})
	Register(&Spec{
		Name: "docker",
		Subcommands: []Subcommand{
			{
				Name: "compose",
				Subcommands: []Subcommand{
					{Name: "up"},
				},
			},
		},
	})

	tests := []struct {
		name         string
		input        string
		wantSkeleton string
		wantOk       bool
	}{
		{
			name:         "git checkout with argument",
			input:        "git checkout feature-x -b",
			wantSkeleton: "git checkout",
			wantOk:       true,
		},
		{
			name:         "git co alias with argument",
			input:        "git co feature-y",
			wantSkeleton: "git checkout",
			wantOk:       true,
		},
		{
			name:         "git remote add with arguments",
			input:        "git remote add origin https://github.com/test/test.git",
			wantSkeleton: "git remote add",
			wantOk:       true,
		},
		{
			name:         "git global flag before subcommand",
			input:        "git -C /tmp checkout main",
			wantSkeleton: "git checkout",
			wantOk:       true,
		},
		{
			name:         "docker compose up with flags",
			input:        "docker compose up -d",
			wantSkeleton: "docker compose up",
			wantOk:       true,
		},
		{
			name:         "unregistered spec",
			input:        "cargo build --release",
			wantSkeleton: "",
			wantOk:       false,
		},
		{
			name:         "empty input",
			input:        "   ",
			wantSkeleton: "",
			wantOk:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := TryExtractSkeleton(tc.input)
			if ok != tc.wantOk {
				t.Fatalf("TryExtractSkeleton(%q) ok = %v, want %v", tc.input, ok, tc.wantOk)
			}
			if got != tc.wantSkeleton {
				t.Errorf("TryExtractSkeleton(%q) = %q, want %q", tc.input, got, tc.wantSkeleton)
			}
		})
	}
}
