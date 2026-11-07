package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
)

var sharedHTTPClient = &http.Client{}

func NewClient(cfg config.ProviderConfig) (Client, error) {
	protocol := strings.ToLower(strings.TrimSpace(cfg.InheritedFrom))
	switch protocol {
	case "openai", "":
		return NewOpenAIClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported ai protocol: %s", cfg.InheritedFrom)
	}
}

type OpenAIClient struct {
	cfg config.ProviderConfig
}

func NewOpenAIClient(cfg config.ProviderConfig) *OpenAIClient {
	return &OpenAIClient{cfg: cfg}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

func (c *OpenAIClient) executeChat(ctx context.Context, messages []chatMessage) (string, error) {
	endpoint := strings.TrimSpace(c.cfg.Endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("empty endpoint in ai provider config")
	}
	if !strings.HasSuffix(endpoint, "/chat/completions") && !strings.Contains(endpoint, "/chat/completions?") {
		endpoint = strings.TrimRight(endpoint, "/") + "/chat/completions"
	}

	reqMap := map[string]any{
		"model":       c.cfg.Model,
		"messages":    messages,
		"max_tokens":  100,
		"temperature": 0.2,
	}
	maps.Copy(reqMap, c.cfg.ExtraRequestBody)

	bodyBytes, err := json.Marshal(reqMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ai request: %w", err)
	}

	timeoutMS := c.cfg.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 2000
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	apiKey := c.cfg.GetAPIKey()
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	res, err := sharedHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return "", fmt.Errorf("ai server returned status %d: %s", res.StatusCode, string(errBody))
	}

	resBytes, err := io.ReadAll(io.LimitReader(res.Body, 65536))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var chatRes chatResponse
	if err := json.Unmarshal(resBytes, &chatRes); err != nil {
		return "", fmt.Errorf("failed to parse ai response json: %w", err)
	}

	if len(chatRes.Choices) == 0 {
		return "", nil
	}

	return chatRes.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	userPrompt := BuildCompletionPrompt(buf, env, dynamicCtx)

	messages := []chatMessage{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: userPrompt},
	}

	rawContent, err := c.executeChat(ctx, messages)
	if err != nil {
		return nil, err
	}
	if rawContent == "" {
		return nil, nil
	}

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
	}, nil
}

func (c *OpenAIClient) SuggestFix(ctx context.Context, failedCmd string, exitCode int, env EnvSnapshot) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	userPrompt := BuildFixPrompt(failedCmd, exitCode, env)

	messages := []chatMessage{
		{Role: "system", Content: FixSystemPrompt},
		{Role: "user", Content: userPrompt},
	}

	rawContent, err := c.executeChat(ctx, messages)
	if err != nil {
		return nil, err
	}
	if rawContent == "" {
		return nil, nil
	}

	cleaned := CleanSuggestion(rawContent)
	if cleaned == "" || cleaned == strings.TrimSpace(failedCmd) {
		return nil, nil
	}

	return &spec.Suggestion{
		Cmd:        cleaned,
		Desc:       "ai fix",
		Icon:       "ai",
		Source:     string(SourceAI),
		Confidence: 85,
	}, nil
}
