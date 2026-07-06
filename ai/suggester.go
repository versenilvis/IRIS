package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/versenilvis/iris/spec"
)

type SourceType string

const (
	SourceHistory SourceType = "history"
	SourceSpec    SourceType = "spec"
	SourceAI      SourceType = "ai"
)

type EnvSnapshot struct {
	Cwd          string
	LastCmd      string
	LastExitCode int
	GitStatus    string
	DirSignature string
}

func (e EnvSnapshot) Hash() string {
	raw := e.Cwd + "|" + e.LastCmd + "|" + strconv.Itoa(e.LastExitCode) + "|" + e.GitStatus + "|" + e.DirSignature
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:8])
}

type Suggester interface {
	Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
}
