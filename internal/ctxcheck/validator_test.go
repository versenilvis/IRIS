package ctxcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/versenilvis/iris/internal/scoring"
)

func TestValidate_EndToEnd(t *testing.T) {
	tmpDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte("test:\n\techo test\n"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"scripts": {"dev": "vite"}}`), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "Makefile"), []byte("build:\n\techo build\n"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "script.sh"), []byte("#!/bin/sh\n"), 0755)

	cases := []struct {
		cmd      string
		expected Verdict
	}{
		{"just test", Valid},
		{"just missing", Invalid},
		{"npm run dev", Valid},
		{"npm run missing", Invalid},
		{"make build", Valid},
		{"make missing", Invalid},
		{"./script.sh", Valid},
		{"./missing.sh", Invalid},
		{"git add . && git commit", Free},
		{"cd /tmp", Valid},
		{"cd nonexistent_dir", Invalid},
		{"cd -", Free},
		{"cd /tmp && just test", Unknown}, // compound with cd
		{"echo $(whoami)", Unknown},       // substitution
	}

	for _, tc := range cases {
		v := Validate(tc.cmd, tmpDir, Posix)
		if v != tc.expected {
			t.Errorf("cmd %q: expected %v, got %v", tc.cmd, tc.expected, v)
		}
	}
}

func TestValidate_Dialects(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "justfile"), []byte("test:\n\techo test\n"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"scripts": {"dev": "vite"}}`), 0644)

	// UnknownDialect always yields Unknown
	for _, cmd := range []string{"just test", "cd /tmp", "npm run dev", "echo hello"} {
		if v := Validate(cmd, tmpDir, UnknownDialect); v != Unknown {
			t.Errorf("UnknownDialect for %q: expected Unknown, got %v", cmd, v)
		}
	}

	// Fish dialect specific features
	fishCases := []struct {
		cmd      string
		expected Verdict
	}{
		{"just test; and npm run dev", Valid},
		{"just test; or just missing", Invalid},
		{"echo (whoami)", Unknown},
		{"just (test)", Unknown},
		{"cd /tmp; and just test", Unknown}, // compound with cd
		{"just test", Valid},
		{"just missing", Invalid},
	}

	for _, tc := range fishCases {
		if v := Validate(tc.cmd, tmpDir, Fish); v != tc.expected {
			t.Errorf("Fish dialect for %q: expected %v, got %v", tc.cmd, tc.expected, v)
		}
	}
}

func TestAllow_Matrix(t *testing.T) {
	verdicts := []Verdict{Invalid, Unknown, Valid, Free}
	tiers := []int{0, 1, 4}
	scopesList := []int{1, 3}

	for _, v := range verdicts {
		for _, tier := range tiers {
			for _, scopes := range scopesList {
				cand := scoring.Candidate{
					Cmd:  "test-cmd",
					Tier: tier,
				}
				scopeOf := func() int { return scopes }
				allowed := Allow(cand, v, scopeOf)

				var expected bool
				switch v {
				case Invalid:
					expected = false
				case Unknown:
					expected = tier == 4
				case Valid:
					expected = true
				case Free:
					expected = tier > 0 || scopes >= scoring.GlobalScopeThreshold
				}

				if allowed != expected {
					t.Errorf("Verdict: %v, Tier: %d, Scopes: %d => expected %v, got %v",
						v, tier, scopes, expected, allowed)
				}
			}
		}
	}
}

func TestCache_Invalidation(t *testing.T) {
	tmpDir := t.TempDir()
	justPath := filepath.Join(tmpDir, "justfile")
	_ = os.WriteFile(justPath, []byte("task1:\n\techo 1\n"), 0644)

	if v := ValidateJust([]string{"just", "task1"}, tmpDir); v != Valid {
		t.Fatalf("expected task1 Valid, got %v", v)
	}
	if v := ValidateJust([]string{"just", "task2"}, tmpDir); v != Invalid {
		t.Fatalf("expected task2 Invalid before edit, got %v", v)
	}

	// append task2 and invalidate cache
	_ = os.WriteFile(justPath, []byte("task1:\n\techo 1\ntask2:\n\techo 2\n"), 0644)
	InvalidateCache()

	if v := ValidateJust([]string{"just", "task2"}, tmpDir); v != Valid {
		t.Fatalf("expected task2 Valid after invalidate, got %v", v)
	}
}
