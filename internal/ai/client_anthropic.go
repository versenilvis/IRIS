package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
)

// AnthropicClient implements Client for Anthropic's native Messages API.
// Docs: https://docs.anthropic.com/en/api/messages
type AnthropicClient struct {
	cfg config.ProviderConfig
}

func NewAnthropicClient(cfg config.ProviderConfig) *AnthropicClient {
	return &AnthropicClient{cfg: cfg}
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature"`
	System      string             `json:"system"`
	Messages    []anthropicMessage `json:"messages"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicResponse struct {
	Content []anthropicContentBlock `json:"content"`
	StopReason string               `json:"stop_reason"`
}

func (c *AnthropicClient) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	endpoint := strings.TrimSpace(c.cfg.Endpoint)
	if endpoint == "" {
		endpoint = "https://api.anthropic.com/v1/messages"
	}

	userPrompt := BuildCompletionPrompt(buf, env, dynamicCtx)

	reqBody := anthropicRequest{
		Model:       c.cfg.Model,
		MaxTokens:   100,
		Temperature: 0.2,
		System:      SystemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal anthropic request: %w", err)
	}

	timeoutMS := c.cfg.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 3000
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	apiKey := c.cfg.GetAPIKey()
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	res, err := sharedHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, fmt.Errorf("anthropic server returned status %d: %s", res.StatusCode, string(errBody))
	}

	resBytes, err := io.ReadAll(io.LimitReader(res.Body, 65536))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var anthroRes anthropicResponse
	if err := json.Unmarshal(resBytes, &anthroRes); err != nil {
		return nil, fmt.Errorf("failed to parse anthropic response json: %w", err)
	}

	if len(anthroRes.Content) == 0 {
		return nil, nil
	}

	rawContent := anthroRes.Content[0].Text
	cleaned := NormalizeSuggestion(buf, rawContent)
	if cleaned == "" || cleaned == strings.TrimSpace(buf) {
		return nil, nil
	}

	return &spec.Suggestion{
		Cmd:        cleaned,
		Desc:       "ai suggestion",
		Icon:       "ai",
		Source:     string(SourceAI),
		Confidence: 85,
		Dangerous:  func() bool { d, _ := IsDangerous(cleaned); return d }(),
	}, nil
}