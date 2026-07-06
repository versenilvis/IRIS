package ai_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/versenilvis/iris/ai"
	"github.com/versenilvis/iris/spec"
)

type mockProvider struct {
	name      string
	matchPref string
	gatherRet string
	calls     int32
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Matches(buf string) bool {
	return len(buf) >= len(m.matchPref) && buf[:len(m.matchPref)] == m.matchPref
}

func (m *mockProvider) Gather(ctx context.Context) (string, error) {
	atomic.AddInt32(&m.calls, 1)
	return m.gatherRet, nil
}

func TestProviderCache_TTL(t *testing.T) {
	cache := ai.NewProviderCache(50 * time.Millisecond)
	provider := &mockProvider{
		name:      "test-prov",
		matchPref: "test",
		gatherRet: "cached-data",
	}
	ctx := context.Background()

	// call 1 -> should gather
	res1 := cache.GetOrGather(ctx, provider)
	if res1 != "cached-data" || atomic.LoadInt32(&provider.calls) != 1 {
		t.Fatalf("expected gather call 1, got res: %q, calls: %d", res1, provider.calls)
	}

	// call 2 immediately -> should hit cache
	res2 := cache.GetOrGather(ctx, provider)
	if res2 != "cached-data" || atomic.LoadInt32(&provider.calls) != 1 {
		t.Fatalf("expected cache hit (calls stay 1), got calls: %d", provider.calls)
	}

	// wait for ttl to expire
	time.Sleep(60 * time.Millisecond)

	// call 3 after ttl -> should gather again
	res3 := cache.GetOrGather(ctx, provider)
	if res3 != "cached-data" || atomic.LoadInt32(&provider.calls) != 2 {
		t.Fatalf("expected gather call 2 after ttl, got calls: %d", provider.calls)
	}
}

func TestAIEngine_DynamicContext(t *testing.T) {
	provider := &mockProvider{
		name:      "docker-mock",
		matchPref: "docker exec",
		gatherRet: "test-container\tnginx",
	}

	engine := ai.NewAIEngine(func(ctx context.Context, buf string, env ai.EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
		if dynamicCtx != "test-container\tnginx" {
			t.Fatalf("expected dynamicCtx to be passed to handler, got: %q", dynamicCtx)
		}
		return &spec.Suggestion{Cmd: "docker exec -it test-container bash", Confidence: 85}, nil
	})
	engine.RegisterProvider(provider)

	ctx := context.Background()
	sugg, err := engine.Suggest(ctx, "docker exec ", ai.EnvSnapshot{}, "")
	if err != nil || sugg == nil {
		t.Fatalf("expected suggestion, got err: %v, sugg: %+v", err, sugg)
	}
	if sugg.Cmd != "docker exec -it test-container bash" {
		t.Fatalf("unexpected cmd: %q", sugg.Cmd)
	}
}
