package scoring

import "testing"

func TestAnchor_PrefixTier(t *testing.T) {
	candidates := []string{
		"cd ~/dev/project",
		"cd /var/log",
		"claude --resume",
	}

	anch := NewAnchor(candidates, "cd ")
	if anch.prefix != "cd " {
		t.Fatalf("expected prefix anchor 'cd ', got %q", anch.prefix)
	}

	if !anch.Allows("cd ~/dev/project") {
		t.Error("expected candidate with matching prefix to be allowed")
	}
	if anch.Allows("claude --resume") {
		t.Error("expected candidate without prefix to be rejected under tier 1")
	}
}

func TestAnchor_HeadCommandTier(t *testing.T) {
	candidates := []string{
		"git status",
		"git commit -m 'fix'",
		"ls -la",
	}

	anch := NewAnchor(candidates, "git che")
	if anch.head != "git" {
		t.Fatalf("expected head anchor 'git', got %q", anch.head)
	}

	if !anch.Allows("git checkout main") {
		t.Error("expected candidate with matching head to be allowed")
	}
	if anch.Allows("ls -la") {
		t.Error("expected unrelated command to be rejected under tier 2")
	}
}

func TestAnchor_FallbackFuzzyTier(t *testing.T) {
	candidates := []string{
		"git checkout",
		"docker compose up",
	}

	// 'gco' has no prefix match and 'gco' is not the head of any candidate
	anch := NewAnchor(candidates, "gco")
	if anch.prefix != "" || anch.head != "" {
		t.Fatalf("expected empty anchor for alias fallback, got prefix=%q head=%q", anch.prefix, anch.head)
	}

	if !anch.Allows("git checkout") {
		t.Error("expected fallback to allow any candidate")
	}
}

func TestIsSubsequenceWithGap(t *testing.T) {
	if !isSubsequenceWithGap("bl", "block", 4) {
		t.Error("expected 'bl' in 'block' to match with gap 4")
	}
	// 'c...t' in 'configure-test': c(0), o(1), n(2), f(3), i(4), g(5), u(6), r(7), e(8), -(9), t(10) -> gap is 9
	if isSubsequenceWithGap("ct", "configure-test", 4) {
		t.Error("expected large gap between c and t to be rejected")
	}
}
