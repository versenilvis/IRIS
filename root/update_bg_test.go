package root

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestChangelogSummaryLines(t *testing.T) {
	body := "## Changelog\n### Bug fixes\n* abc1234  fix something\n* def5678  fix another thing\n* ghi9012  a third fix\n## Update\n```bash\niris update\n```\n"

	got := changelogSummaryLines(body, 2)
	want := []string{"fix something", "fix another thing"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(got), got)
	}
	for i, line := range want {
		if got[i] != line {
			t.Errorf("line %d: expected %q, got %q", i, line, got[i])
		}
	}
}

func TestChangelogSummaryLinesEmptyBody(t *testing.T) {
	if got := changelogSummaryLines("", 2); len(got) != 0 {
		t.Errorf("expected no lines for an empty body, got %v", got)
	}
}

func TestPrintUpdateNoticeIncludesChangelogSummary(t *testing.T) {
	originalVersion := Version
	Version = "v1.0.0"
	t.Cleanup(func() { Version = originalVersion })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	originalStdout := os.Stdout
	os.Stdout = w

	printUpdateNotice("v1.1.0", "## Changelog\n### Bug fixes\n* abc1234  fix something\n* def5678  fix another thing\n* ghi9012  a third fix\n")

	_ = w.Close()
	os.Stdout = originalStdout
	out, _ := io.ReadAll(r)
	got := string(out)

	if !strings.Contains(got, "v1.0.0") || !strings.Contains(got, "v1.1.0") {
		t.Errorf("expected the version transition in the notice, got %q", got)
	}
	if !strings.Contains(got, "fix something") || !strings.Contains(got, "fix another thing") {
		t.Errorf("expected the first two changelog entries in the notice, got %q", got)
	}
	if strings.Contains(got, "a third fix") {
		t.Errorf("expected the summary to cap at 2 entries, got %q", got)
	}
}

func TestPrintUpdateNoticeWithoutNotes(t *testing.T) {
	originalVersion := Version
	Version = "v1.0.0"
	t.Cleanup(func() { Version = originalVersion })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	originalStdout := os.Stdout
	os.Stdout = w

	printUpdateNotice("v1.1.0", "")

	_ = w.Close()
	os.Stdout = originalStdout
	out, _ := io.ReadAll(r)
	got := string(out)

	if !strings.Contains(got, "v1.1.0") {
		t.Errorf("expected the notice to still show without changelog notes, got %q", got)
	}
}
