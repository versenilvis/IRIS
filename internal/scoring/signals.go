package scoring

import (
	"context"
	"strings"

	"github.com/versenilvis/iris/internal/workspace"
)

type SignalSet struct {
	Workspace      workspace.WorkspaceInfo
	LocalFrecency  []FrecencyEntry
	GlobalFrecency []FrecencyEntry
	Query          string
	RootCommand    string
	Cwd            string
}

// CollectSignals gathers environment, workspace, and historical frecency signals for the given query and directory
func CollectSignals(ctx context.Context, cwd, query, rootCmd string, frecency *FrecencyStore) SignalSet {
	ws := workspace.DetectCached(cwd)

	if ctx == nil {
		ctx = context.Background()
	}

	var local, global []FrecencyEntry
	if frecency != nil {
		local, _ = frecency.QueryLocal(ctx, cwd, query, 50)
		global, _ = frecency.QueryGlobal(ctx, query, 50)
	}

	return SignalSet{
		Workspace:      ws,
		LocalFrecency:  local,
		GlobalFrecency: global,
		Query:          strings.TrimSpace(query),
		RootCommand:    strings.TrimSpace(rootCmd),
		Cwd:            cwd,
	}
}
