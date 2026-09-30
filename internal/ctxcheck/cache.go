package ctxcheck

import (
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type pathEntry struct {
	path   string
	found  bool
	expiry time.Time
}

type manifestEntry struct {
	mtime time.Time
	size  int64
	data  any
}

type sfCall struct {
	wg  sync.WaitGroup
	val any
	err error
}

type singleflightGroup struct {
	mu sync.Mutex
	m  map[string]*sfCall
}

func (g *singleflightGroup) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*sfCall)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(sfCall)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

var (
	cacheMu       sync.RWMutex
	pathCache     = make(map[string]pathEntry)
	manifestCache = make(map[string]manifestEntry)
	manifestSF    singleflightGroup

	slowDirsMu sync.RWMutex
	slowDirs   = make(map[string]time.Time)
)

const manifestTTL = 1000 * time.Millisecond

func InvalidateCache() {
	cacheMu.Lock()
	clear(pathCache)
	clear(manifestCache)
	cacheMu.Unlock()

	slowDirsMu.Lock()
	clear(slowDirs)
	slowDirsMu.Unlock()
}

func markSlowDir(dir string, d time.Duration) {
	slowDirsMu.Lock()
	defer slowDirsMu.Unlock()
	now := time.Now()
	for k, exp := range slowDirs {
		if now.After(exp) {
			delete(slowDirs, k)
		}
	}
	slowDirs[dir] = now.Add(d)
}

func isSlowDir(dir string) bool {
	now := time.Now()
	slowDirsMu.RLock()
	expiry, ok := slowDirs[dir]
	if !ok {
		slowDirsMu.RUnlock()
		return false
	}
	if now.Before(expiry) {
		slowDirsMu.RUnlock()
		return true
	}
	slowDirsMu.RUnlock()

	slowDirsMu.Lock()
	if exp, stillOk := slowDirs[dir]; stillOk && now.After(exp) {
		delete(slowDirs, dir)
	}
	slowDirsMu.Unlock()
	return false
}

func getCachedPath(kind, cwd string, findFn func(string) (string, error)) (string, error) {
	key := kind + ":" + cwd
	now := time.Now()

	cacheMu.RLock()
	entry, ok := pathCache[key]
	cacheMu.RUnlock()

	if ok && now.Before(entry.expiry) {
		if !entry.found {
			return "", os.ErrNotExist
		}
		return entry.path, nil
	}

	// expired or missing
	if ok && entry.found {
		// fast stat on the known path rather than traversing upward
		if fi, err := os.Stat(entry.path); err == nil && !fi.IsDir() {
			cacheMu.Lock()
			entry.expiry = now.Add(manifestTTL)
			pathCache[key] = entry
			cacheMu.Unlock()
			return entry.path, nil
		}
	}

	// traverse upward with singleflight to avoid duplicate goroutines on hung directories
	res, err := manifestSF.Do(key, func() (any, error) {
		p, findErr := findFn(cwd)
		cacheMu.Lock()
		if findErr == nil {
			pathCache[key] = pathEntry{
				path:   p,
				found:  true,
				expiry: time.Now().Add(manifestTTL),
			}
		} else {
			pathCache[key] = pathEntry{
				found:  false,
				expiry: time.Now().Add(manifestTTL),
			}
		}
		cacheMu.Unlock()
		return p, findErr
	})

	if res != nil {
		if pStr, isStr := res.(string); isStr {
			return pStr, err
		}
	}
	return "", err
}

func getCachedManifest(path string, parseFn func(string) (any, error)) (any, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	mtime := fi.ModTime()
	size := fi.Size()

	cacheMu.RLock()
	entry, ok := manifestCache[path]
	cacheMu.RUnlock()

	if ok && entry.mtime.Equal(mtime) && entry.size == size {
		return entry.data, nil
	}

	// parse without holding lock
	data, err := parseFn(path)
	if err != nil {
		return nil, err
	}

	cacheMu.Lock()
	manifestCache[path] = manifestEntry{
		mtime: mtime,
		size:  size,
		data:  data,
	}
	cacheMu.Unlock()

	return data, nil
}

var testStatDelayHook atomic.Pointer[func(string)]

func SetTestStatDelayHook(fn func(string)) {
	if fn == nil {
		testStatDelayHook.Store(nil)
		return
	}
	testStatDelayHook.Store(&fn)
}

func ValidateWithTimeout(cmd, cwd string, d Dialect, timeout time.Duration) Verdict {
	if isSlowDir(cwd) {
		return Unknown
	}

	done := make(chan Verdict, 1)
	go func() {
		if hookPtr := testStatDelayHook.Load(); hookPtr != nil && *hookPtr != nil {
			(*hookPtr)(cwd)
		}
		done <- Validate(cmd, cwd, d)
	}()

	select {
	case v := <-done:
		return v
	case <-time.After(timeout):
		markSlowDir(cwd, 5*time.Second)
		return Unknown
	}
}
