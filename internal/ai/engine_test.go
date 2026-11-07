package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
)

func TestEnvSnapshot_Hash(t *testing.T) {
	snap1 := EnvSnapshot{
		Cwd:          "/home/user",
		LastCmd:      "ls -l",
		LastExitCode: 0,
		GitStatus:    "clean",
		DirSignature: "sig1",
	}
	hash1 := snap1.Hash()
	if len(hash1) != 16 {
		t.Fatalf("expected 16 hex chars, got len %d: %q", len(hash1), hash1)
	}

	snap2 := snap1
	snap2.LastCmd = "pwd"
	if snap1.Hash() == snap2.Hash() {
		t.Fatalf("expected different hash when field changes")
	}

	// Verify that fields containing delimiter characters do not collide
	snapA := EnvSnapshot{Cwd: "/home/user", LastCmd: "foo|bar"}
	snapB := EnvSnapshot{Cwd: "/home/user|foo", LastCmd: "bar"}
	if snapA.Hash() == snapB.Hash() {
		t.Fatalf("expected different hash for distinct snapshots containing delimiter characters")
	}
}

func TestAIEngine_Suggest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"role": "assistant", "content": "git commit -m \"feat: update main.go\""}},
			},
		}
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	cfg := config.Get()
	origAI := cfg.AI
	defer func() { cfg.AI = origAI }()

	cfg.AI.Enabled = true
	cfg.AI.Provider = "test-provider"
	cfg.AI.Providers = map[string]config.ProviderConfig{
		"test-provider": {
			InheritedFrom: "openai",
			Endpoint:      server.URL,
			Model:         "test-model",
		},
	}

	engine := NewAIEngine(nil)
	ctx := context.Background()
	snap := EnvSnapshot{GitStatus: "modified main.go"}

	sugg, err := engine.Suggest(ctx, "git commit -m \"", snap, "")
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if sugg == nil || sugg.Confidence != 85 || sugg.Source != string(SourceAI) {
		t.Fatalf("unexpected suggestion: %+v", sugg)
	}
}

func TestAIEngine_Suggest_Cancel(t *testing.T) {
	engine := NewAIEngine(func(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
		time.Sleep(50 * time.Millisecond)
		return &spec.Suggestion{Cmd: "test", Confidence: 90}, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately before calling

	sugg, err := engine.Suggest(ctx, "git commit -m \"", EnvSnapshot{}, "")
	if err == nil || sugg != nil {
		t.Fatalf("expected cancellation error and nil suggestion, got sugg: %+v, err: %v", sugg, err)
	}
}

func TestShouldOverwrite(t *testing.T) {
	newSugg := &spec.Suggestion{
		Cmd:        "git commit -m \"feat: new feature\"",
		Confidence: 85,
		Source:     string(SourceAI),
	}

	// case 1: exact match and higher confidence
	if !ShouldOverwrite("git commit -m \"", "git commit -m \"f", newSugg, 70) {
		t.Fatalf("expected true for prefix match with higher confidence")
	}

	// case 2: user backspaced / typed different prefix
	if ShouldOverwrite("git commit -m \"", "git co", newSugg, 70) {
		t.Fatalf("expected false when current buf does not match original prefix")
	}

	// case 3: new suggestion does not match what user typed
	if ShouldOverwrite("git commit -m \"", "git commit -m \"fix", newSugg, 70) {
		t.Fatalf("expected false when suggestion cmd does not match current buf")
	}

	// case 4: lower or equal confidence
	if ShouldOverwrite("git commit -m \"", "git commit -m \"", newSugg, 90) {
		t.Fatalf("expected false when new confidence is lower")
	}
}

// Verify that caller mutations on returned fresh suggestions do not corrupt the internal engine cache
func TestAIEngine_Suggest_Immutability(t *testing.T) {
	engine := NewAIEngine(func(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
		return &spec.Suggestion{Cmd: "echo original", Confidence: 90}, nil
	})
	ctx := context.Background()

	sugg1, err := engine.Suggest(ctx, "echo ", EnvSnapshot{}, "")
	if err != nil || sugg1 == nil {
		t.Fatalf("expected suggestion, got err: %v, sugg: %v", err, sugg1)
	}

	// Mutate returned suggestion
	sugg1.Cmd = "echo corrupted"

	// Fetch from cache via prefix match
	sugg2, err := engine.Suggest(ctx, "echo ", EnvSnapshot{}, "")
	if err != nil || sugg2 == nil {
		t.Fatalf("expected cached suggestion, got err: %v, sugg: %v", err, sugg2)
	}
	if sugg2.Cmd != "echo original" {
		t.Fatalf("expected cache to remain 'echo original', got corrupted cmd: %q", sugg2.Cmd)
	}
}

func TestAIEngine_SuggestFix(t *testing.T) {
	ctx := context.Background()
	engine := NewAIEngine(nil)

	// empty command or zero exit code returns nil
	if s, err := engine.SuggestFix(ctx, "", 1, EnvSnapshot{}); err != nil || s != nil {
		t.Errorf("expected nil for empty command, got %v, %v", s, err)
	}
	if s, err := engine.SuggestFix(ctx, "git push", 0, EnvSnapshot{}); err != nil || s != nil {
		t.Errorf("expected nil for zero exit code, got %v, %v", s, err)
	}

	callCount := 0
	engine.SetFixHandler(func(ctx context.Context, failedCmd string, exitCode int, env EnvSnapshot) (*spec.Suggestion, error) {
		callCount++
		return &spec.Suggestion{
			Cmd:        "git push origin main",
			Desc:       "ai fix",
			Icon:       "ai",
			Source:     "ai",
			Confidence: 85,
		}, nil
	})

	env := EnvSnapshot{Cwd: "/home/user", LastCmd: "git psh origin main", LastExitCode: 1}
	sugg1, err := engine.SuggestFix(ctx, "git psh origin main", 1, env)
	if err != nil || sugg1 == nil {
		t.Fatalf("expected fix suggestion, got sugg: %v, err: %v", sugg1, err)
	}
	if sugg1.Cmd != "git push origin main" {
		t.Errorf("expected fixed command, got %q", sugg1.Cmd)
	}
	if callCount != 1 {
		t.Errorf("expected handler called once, got %d", callCount)
	}

	// second call with same snapshot should hit cache
	sugg2, err := engine.SuggestFix(ctx, "git psh origin main", 1, env)
	if err != nil || sugg2 == nil {
		t.Fatalf("expected cached fix suggestion, got sugg: %v, err: %v", sugg2, err)
	}
	if sugg2.Cmd != "git push origin main" {
		t.Errorf("expected cached command, got %q", sugg2.Cmd)
	}
	if callCount != 1 {
		t.Errorf("expected handler not called on cache hit, callCount=%d", callCount)
	}
}
