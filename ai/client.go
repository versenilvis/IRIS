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

	"github.com/versenilvis/iris/config"
	"github.com/versenilvis/iris/spec"
)

var sharedHTTPClient = &http.Client{}

type Client interface {
	Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error)
}

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

func (c *OpenAIClient) Suggest(ctx context.Context, buf string, env EnvSnapshot, dynamicCtx string) (*spec.Suggestion, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	endpoint := strings.TrimSpace(c.cfg.Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("empty endpoint in ai provider config")
	}
	if !strings.HasSuffix(endpoint, "/chat/completions") && !strings.Contains(endpoint, "/chat/completions?") {
		endpoint = strings.TrimRight(endpoint, "/") + "/chat/completions"
	}

	systemPrompt := "You are a concise shell command completion assistant. Provide ONLY the completed shell command line. Do not explain, do not use markdown formatting, and do not wrap the command in code blocks or backticks. Always ensure valid shell syntax: if an argument contains spaces or parentheses (such as git commit messages), you MUST wrap that argument in double quotes \"...\"."
	userPrompt := fmt.Sprintf("Complete this shell command line: %s\nContext:\nCwd: %s\nLastCmd: %s\nLastExitCode: %d\nGitStatus: %s\nRecentCmds: %v\nDynamicContext: %s",
		buf, env.Cwd, env.LastCmd, env.LastExitCode, env.GitStatus, env.RecentCmds, dynamicCtx)

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	reqMap := map[string]any{
		"model":       c.cfg.Model,
		"messages":    messages,
		"max_tokens":  100,
		"temperature": 0.2,
	}
	for k, v := range c.cfg.ExtraRequestBody {
		reqMap[k] = v
	}

	bodyBytes, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ai request: %w", err)
	}

	timeoutMS := c.cfg.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 2000
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	apiKey := c.cfg.GetAPIKey()
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	res, err := sharedHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, fmt.Errorf("ai server returned status %d: %s", res.StatusCode, string(errBody))
	}

	resBytes, err := io.ReadAll(io.LimitReader(res.Body, 65536))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var chatRes chatResponse
	if err := json.Unmarshal(resBytes, &chatRes); err != nil {
		return nil, fmt.Errorf("failed to parse ai response json: %w", err)
	}

	if len(chatRes.Choices) == 0 {
		return nil, nil
	}

	rawContent := chatRes.Choices[0].Message.Content
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

func CleanSuggestion(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) > 1 {
			endIdx := len(lines)
			if strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
				endIdx = len(lines) - 1
			}
			s = strings.TrimSpace(strings.Join(lines[1:endIdx], "\n"))
		}
	}
	if len(s) >= 2 && strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`") && !strings.HasPrefix(s, "``") {
		s = s[1 : len(s)-1]
	}
	if len(s) >= 2 && ((strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))) {
		inner := s[1 : len(s)-1]
		if !strings.ContainsAny(inner, "\"'") {
			s = inner
		}
	}
	return strings.TrimSpace(s)
}
