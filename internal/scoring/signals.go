package scoring

import (
	"context"
	"strings"

	"github.com/versenilvis/iris/internal/workspace"
)

type SignalSet struct {
	Workspace         workspace.WorkspaceInfo
	LocalFrecency     []FrecencyEntry
	GlobalFrecency    []FrecencyEntry
	TransitionEntries []TransitionEntry
	TransitionIsLocal bool
	Query             string
	RootCommand       string
	Cwd               string
}

// CollectSignals gathers environment, workspace, and historical frecency/transition signals for the given query and directory
func CollectSignals(ctx context.Context, cwd, query, rootCmd string, frecency *FrecencyStore, prevCmdSkeleton string) SignalSet {
	ws := workspace.DetectCached(cwd)

	if ctx == nil {
		ctx = context.Background()
	}

	var local, global []FrecencyEntry
	var trans []TransitionEntry
	var transIsLocal bool

	if frecency != nil {
		local, _ = frecency.QueryLocal(ctx, cwd, query, 50)
		global, _ = frecency.QueryGlobal(ctx, query, 50)
		if prevCmdSkeleton != "" {
			trans, transIsLocal = frecency.QueryTransitionsWithFallback(ctx, prevCmdSkeleton, cwd)
		}
	}

	return SignalSet{
		Workspace:         ws,
		LocalFrecency:     local,
		GlobalFrecency:    global,
		TransitionEntries: trans,
		TransitionIsLocal: transIsLocal,
		Query:             strings.TrimSpace(query),
		RootCommand:       strings.TrimSpace(rootCmd),
		Cwd:               cwd,
	}
}
