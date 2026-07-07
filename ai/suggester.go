package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/versenilvis/iris/spec"
)

type SourceType string

const (
	SourceHistory SourceType = "history"
	SourceSpec    SourceType = "spec"
	SourceAI      SourceType = "ai"
)

type EnvSnapshot struct {
	Cwd                 string
	LastCmd             string
	LastExitCode        int
	GitStatus           string
	DirSignature        string
	RecentCmds          []string
	GitMergeInProgress  bool
	GitRebaseInProgress bool
}

func NewEnvSnapshot(cwd string, lastCmd string, lastExitCode int, recentCmds []string) EnvSnapshot {
	_, mergeErr := os.Stat(filepath.Join(cwd, ".git", "MERGE_HEAD"))
	_, rebaseErr := os.Stat(filepath.Join(cwd, ".git", "REBASE_HEAD"))
	return EnvSnapshot{
		Cwd:                 cwd,
		LastCmd:             lastCmd,
		LastExitCode:        lastExitCode,
		RecentCmds:          recentCmds,
		GitMergeInProgress:  mergeErr == nil,
		GitRebaseInProgress: rebaseErr == nil,
	}
}

func (e EnvSnapshot) Hash() string {
	raw := e.Cwd + "|" + e.LastCmd + "|" + strconv.Itoa(e.LastExitCode) + "|" + e.GitStatus + "|" + e.DirSignature + "|" + strings.Join(e.RecentCmds, ";") + "|" + strconv.FormatBool(e.GitMergeInProgress) + "|" + strconv.FormatBool(e.GitRebaseInProgress)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:8])
}

type Suggester interface {
	Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
}
