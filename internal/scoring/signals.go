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
	SequenceEntries   []SequenceEntry
	SequenceIsLocal   bool
	Query             string
	RootCommand       string
	Cwd               string
}

// CollectSignals gathers environment, workspace, and historical frecency/transition signals for the given query and directory
func CollectSignals(ctx context.Context, cwd, query, rootCmd string, frecency *FrecencyStore, prevCmdSkeleton string, prevCmd ...string) SignalSet {
	ws := workspace.DetectCached(cwd)

	if ctx == nil {
		ctx = context.Background()
	}

	var local, global []FrecencyEntry
	var trans []TransitionEntry
	var transIsLocal bool
	var seqs []SequenceEntry
	var seqsIsLocal bool

	if frecency != nil {
		local, _ = frecency.QueryLocal(ctx, cwd, query, 50)
		global, _ = frecency.QueryGlobal(ctx, query, 50)
		if len(prevCmd) > 0 && prevCmd[0] != "" {
			seqs, seqsIsLocal = frecency.QuerySequencesWithFallback(ctx, prevCmd[0], cwd)
		}
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
		SequenceEntries:   seqs,
		SequenceIsLocal:   seqsIsLocal,
		Query:             strings.TrimSpace(query),
		RootCommand:       strings.TrimSpace(rootCmd),
		Cwd:               cwd,
	}
}
