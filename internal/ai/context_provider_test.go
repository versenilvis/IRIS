package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
	cache := NewProviderCache(50 * time.Millisecond)
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

	engine := NewAIEngine(func(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
		if dynamicCtx != "test-container\tnginx" {
			t.Fatalf("expected dynamicCtx to be passed to handler, got: %q", dynamicCtx)
		}
		return &spec.Suggestion{Cmd: "docker exec -it test-container bash", Confidence: 85}, nil
	})
	engine.RegisterProvider(provider)

	ctx := context.Background()
	sugg, err := engine.Suggest(ctx, "docker exec ", EnvSnapshot{}, "")
	if err != nil || sugg == nil {
		t.Fatalf("expected suggestion, got err: %v, sugg: %+v", err, sugg)
	}
	if sugg.Cmd != "docker exec -it test-container bash" {
		t.Fatalf("unexpected cmd: %q", sugg.Cmd)
	}
}

// Verify that cache evicts expired entries and resets when exceeding 50 items to prevent unbounded memory growth
func TestProviderCache_Eviction(t *testing.T) {
	cache := NewProviderCache(15 * time.Millisecond)
	ctx := context.Background()

	// 1. Fill cache to capacity (50 items) rapidly before any TTL expires
	for i := 0; i < 50; i++ {
		p := &mockProvider{
			name:      fmt.Sprintf("prov-%d", i),
			matchPref: "test",
			gatherRet: "data",
		}
		cache.GetOrGather(ctx, p)
	}

	cache.mu.Lock()
	if len(cache.entries) != 50 {
		t.Fatalf("expected cache to hold exactly 50 entries, got %d", len(cache.entries))
	}
	cache.mu.Unlock()

	// 2. Add the 51st item immediately when none are expired -> triggers map reset to bound memory
	p51 := &mockProvider{
		name:      "prov-50",
		matchPref: "test",
		gatherRet: "data",
	}
	cache.GetOrGather(ctx, p51)

	cache.mu.Lock()
	if len(cache.entries) != 1 {
		t.Fatalf("expected cache reset when exceeding 50 unexpired items (expected len 1, got %d)", len(cache.entries))
	}
	if _, ok := cache.entries["prov-50"]; !ok {
		t.Fatalf("expected 'prov-50' to exist after reset")
	}
	cache.mu.Unlock()

	// 3. Fill up to 50 items again
	for i := 51; i < 100; i++ {
		p := &mockProvider{
			name:      fmt.Sprintf("prov-%d", i),
			matchPref: "test",
			gatherRet: "data",
		}
		cache.GetOrGather(ctx, p)
	}

	// 4. Wait for all items to expire
	time.Sleep(25 * time.Millisecond)

	// 5. Add a new item when cache has 50 expired items -> should delete expired entries without complete map recreation
	pNext := &mockProvider{
		name:      "prov-100",
		matchPref: "test",
		gatherRet: "data",
	}
	cache.GetOrGather(ctx, pNext)

	cache.mu.Lock()
	if len(cache.entries) != 1 {
		t.Fatalf("expected 50 expired items to be evicted before adding new entry (expected len 1, got %d)", len(cache.entries))
	}
	if _, ok := cache.entries["prov-100"]; !ok {
		t.Fatalf("expected new entry 'prov-100' to exist after eviction")
	}
	cache.mu.Unlock()
}

// Verify that CommandContextProvider caps gathered output to 1000 characters to protect token budget
func TestCommandContextProvider_Truncation(t *testing.T) {
	provider := &CommandContextProvider{
		NameStr:   "test-trunc",
		Prefixes:  []string{"echo"},
		GatherCmd: []string{"go", "env"},
		Label:     "GoEnv",
	}
	ctx := context.Background()
	res, err := provider.Gather(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res, "GoEnv:\n") {
		t.Fatalf("expected label prefix, got: %q", res)
	}
	if len(res) > 1100 {
		t.Fatalf("expected gathered output to be truncated around 1000 characters, got len: %d", len(res))
	}
}

// Verify that concurrent provider registration and context gathering do not cause data races
func TestAIEngine_ConcurrentRegistrationAndGather(t *testing.T) {
	engine := NewAIEngine(nil)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p := &mockProvider{
				name:      fmt.Sprintf("conc-prov-%d", idx),
				matchPref: "docker",
				gatherRet: "conc-data",
			}
			engine.RegisterProvider(p)
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			engine.GatherDynamicContext(ctx, "docker ps", "/tmp")
		}()
	}

	wg.Wait()
}
