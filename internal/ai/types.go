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
	// Use length-prefixed encoding for each field to prevent hash collisions when values contain delimiter characters
	var sb strings.Builder
	enc := func(s string) {
		sb.WriteString(strconv.Itoa(len(s)))
		sb.WriteByte(':')
		sb.WriteString(s)
	}
	enc(e.Cwd)
	enc(e.LastCmd)
	enc(strconv.Itoa(e.LastExitCode))
	enc(e.GitStatus)
	enc(e.DirSignature)
	enc(strconv.Itoa(len(e.RecentCmds)))
	for _, cmd := range e.RecentCmds {
		enc(cmd)
	}
	enc(strconv.FormatBool(e.GitMergeInProgress))
	enc(strconv.FormatBool(e.GitRebaseInProgress))

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:8])
}

type ContextProvider interface {
	Name() string
	Matches(buf string) bool
	Gather(ctx context.Context) (string, error)
}

type Client interface {
	Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
}

type Suggester interface {
	Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
}

type ContextSuggester interface {
	SuggestOnEmpty(ctx context.Context, env EnvSnapshot) (*spec.Suggestion, error)
}

type AIHandler func(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
