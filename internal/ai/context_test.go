package ai

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/versenilvis/iris/spec"
)

type mockAISuggester struct {
	calls int32
	ret   *spec.Suggestion
}

func (m *mockAISuggester) SuggestOnEmpty(ctx context.Context, env EnvSnapshot) (*spec.Suggestion, error) {
	atomic.AddInt32(&m.calls, 1)
	return m.ret, nil
}

func TestRuleBasedSuggester(t *testing.T) {
	rule := RuleBasedSuggester{}

	// Case 1: Last exit code != 0
	ctx := context.Background()
	sugg1, _ := rule.SuggestOnEmpty(ctx, EnvSnapshot{LastExitCode: 1, LastCmd: "make build"})
	if sugg1 == nil || sugg1.Cmd != "make build" || sugg1.Confidence != 80 {
		t.Fatalf("expected retry make build with conf 80, got: %+v", sugg1)
	}

	// case 2: git status after git status
	sugg2, _ := rule.SuggestOnEmpty(ctx, EnvSnapshot{LastCmd: "git status"})
	if sugg2 == nil || sugg2.Cmd != "git diff" || sugg2.Confidence != 75 {
		t.Fatalf("expected git diff with conf 75, got: %+v", sugg2)
	}

	// case 3: modified git files
	sugg3, _ := rule.SuggestOnEmpty(ctx, EnvSnapshot{GitStatus: " M main.go"})
	if sugg3 == nil || sugg3.Cmd != "git status" || sugg3.Confidence != 70 {
		t.Fatalf("expected git status with conf 70, got: %+v", sugg3)
	}

	// case 4: package.json signature (conf 65)
	sugg4, _ := rule.SuggestOnEmpty(ctx, EnvSnapshot{DirSignature: "package.json"})
	if sugg4 == nil || sugg4.Cmd != "npm run dev" || sugg4.Confidence != 65 {
		t.Fatalf("expected npm run dev with conf 65, got: %+v", sugg4)
	}
}

func TestContextCache_ShouldCallAI(t *testing.T) {
	cache := NewContextCache()
	snap := EnvSnapshot{Cwd: "/test", GitStatus: "clean"}

	// first call -> true
	if !cache.ShouldCallAI(snap, 50*time.Millisecond) {
		t.Fatalf("expected true for initial call")
	}
	cache.Update(snap, &spec.Suggestion{Cmd: "test"})

	// second call with same snap -> false
	if cache.ShouldCallAI(snap, 50*time.Millisecond) {
		t.Fatalf("expected false when hash has not changed")
	}

	// wait for min interval before calling again
	time.Sleep(60 * time.Millisecond)

	// third call with different snap -> true
	snap.GitStatus = "dirty"
	if !cache.ShouldCallAI(snap, 50*time.Millisecond) {
		t.Fatalf("expected true when hash changed")
	}
}

func TestEmptyLinePredictor_TwoTier(t *testing.T) {
	ctx := context.Background()
	mockAI := &mockAISuggester{ret: &spec.Suggestion{Cmd: "ai-suggested-cmd", Confidence: 85}}
	predictor := NewEmptyLinePredictor(nil, mockAI, 50*time.Millisecond)

	env1 := EnvSnapshot{LastExitCode: 1, LastCmd: "failed-cmd"}
	sugg1, _ := predictor.Predict(ctx, env1, true)
	if sugg1 == nil || sugg1.Cmd != "failed-cmd" || atomic.LoadInt32(&mockAI.calls) != 0 {
		t.Fatalf("expected rule based result and 0 ai calls, got sugg: %+v, calls: %d", sugg1, mockAI.calls)
	}

	// case 2: rule based conf < 70 -> ai called
	env2 := EnvSnapshot{DirSignature: "package.json"}
	sugg2, _ := predictor.Predict(ctx, env2, true)
	if sugg2 == nil || sugg2.Cmd != "ai-suggested-cmd" || atomic.LoadInt32(&mockAI.calls) != 1 {
		t.Fatalf("expected ai result and 1 ai call, got sugg: %+v, calls: %d", sugg2, mockAI.calls)
	}

	// case 3: same env as case 2 immediately -> ai not called again (cached)
	sugg3, _ := predictor.Predict(ctx, env2, true)
	if sugg3 == nil || sugg3.Cmd != "ai-suggested-cmd" || atomic.LoadInt32(&mockAI.calls) != 1 {
		t.Fatalf("expected cached ai result and 1 ai call (no increment), got calls: %d", mockAI.calls)
	}
}

// Verify that Makefile target extraction ignores variable assignments containing operators like := or colons in values
func TestExtractScriptsAndTargets_Makefile(t *testing.T) {
	tmp := t.TempDir()
	content := []byte("CFLAGS := -O2\nPREFIX ?= /usr/local\nPATH = /bin:/usr/bin\nall: build\nbuild:\n\techo build\n")
	_ = os.WriteFile(filepath.Join(tmp, "Makefile"), content, 0644)

	var sb strings.Builder
	ExtractScriptsAndTargets(&sb, tmp, "")
	res := sb.String()

	if !strings.Contains(res, "build") || !strings.Contains(res, "all") {
		t.Fatalf("expected real targets build and all in result, got: %q", res)
	}
	if strings.Contains(res, "CFLAGS") || strings.Contains(res, "PREFIX") || strings.Contains(res, "PATH") {
		t.Fatalf("expected variable assignments CFLAGS, PREFIX, PATH to be skipped, got: %q", res)
	}
}
