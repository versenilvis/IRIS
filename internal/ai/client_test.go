package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/versenilvis/iris/internal/config"
)

func TestCleanSuggestion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  docker run -d nginx  ", "docker run -d nginx"},
		{"```bash\ngit status\n```", "git status"},
		{"```\nls -la\n```", "ls -la"},
		{"`npm run dev`", "npm run dev"},
		{"\"docker ps\"", "docker ps"},
		{"'git diff'", "git diff"},
		{"\"git commit -m 'hello'\"", "\"git commit -m 'hello'\""},
	}

	for _, tt := range tests {
		got := CleanSuggestion(tt.input)
		if got != tt.expected {
			t.Errorf("CleanSuggestion(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestOpenAIClient_Suggest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected post method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-secret-key" {
			t.Errorf("expected bearer token, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var reqMap map[string]any
		if err := json.Unmarshal(body, &reqMap); err != nil {
			t.Fatalf("failed to parse request json: %v", err)
		}
		if reqMap["model"] != "test-model-32b" {
			t.Errorf("expected model test-model-32b, got %v", reqMap["model"])
		}
		if reqMap["temperature"] != 0.5 {
			t.Errorf("expected extra temperature 0.5, got %v", reqMap["temperature"])
		}

		res := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "```bash\nkubectl get pods -n kube-system\n```"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	cfg := config.ProviderConfig{
		InheritedFrom: "openai",
		Endpoint:      server.URL,
		APIKey:        "test-secret-key",
		Model:         "test-model-32b",
		TimeoutMS:     1000,
		ExtraRequestBody: map[string]any{
			"temperature": 0.5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	env := EnvSnapshot{Cwd: "/home/user", LastCmd: "kubectl get", LastExitCode: 0}
	sugg, err := client.Suggest(ctx, "kubectl get p", env, "")
	if err != nil {
		t.Fatalf("suggest failed: %v", err)
	}
	if sugg == nil {
		t.Fatalf("expected suggestion, got nil")
	}
	if sugg.Cmd != "kubectl get pods -n kube-system" {
		t.Errorf("expected cleaned cmd, got %q", sugg.Cmd)
	}
	if sugg.Confidence != 85 {
		t.Errorf("expected confidence 85, got %d", sugg.Confidence)
	}
}

func TestOpenAIClient_TimeoutAndCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.ProviderConfig{
		InheritedFrom: "openai",
		Endpoint:      server.URL,
		TimeoutMS:     50,
	}

	client := NewOpenAIClient(cfg)
	ctx := context.Background()
	env := EnvSnapshot{}

	_, err := client.Suggest(ctx, "sleep", env, "")
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Suggest(ctxCancel, "sleep", env, "")
	if err == nil {
		t.Errorf("expected context canceled error, got nil")
	}
}

func TestBuildCompletionPrompt(t *testing.T) {
	env := EnvSnapshot{
		Cwd:          "/home/user/project",
		LastCmd:      "",
		LastExitCode: 0,
		GitStatus:    "",
		RecentCmds: []string{
			"git status",
			"git commit -m \"fix(auth): update\"",
		},
	}
	prompt := BuildCompletionPrompt("docker exec ", env, "Running containers: app (nginx)")

	if !strings.Contains(prompt, "Input buffer (must appear verbatim at the start of your output):\ndocker exec ") {
		t.Errorf("prompt missing verbatim input buffer instructions: %s", prompt)
	}
	if strings.Contains(prompt, "GitStatus:") || strings.Contains(prompt, "PreviousCommand (already finished, exit code ") {
		t.Errorf("prompt should omit empty GitStatus or PreviousCommand: %s", prompt)
	}
	if !strings.Contains(prompt, "  git status\n  git commit -m \"fix(auth): update\"\n") {
		t.Errorf("prompt should format RecentCmds one per line: %s", prompt)
	}
	if !strings.Contains(prompt, "DynamicContext:\nRunning containers: app (nginx)") {
		t.Errorf("prompt missing dynamic context: %s", prompt)
	}
}

func TestNormalizeSuggestion(t *testing.T) {
	tests := []struct {
		name     string
		buf      string
		raw      string
		expected string
	}{
		{
			name:     "Verbatim prefix unchanged",
			buf:      "docker exec -it ",
			raw:      "docker exec -it app-server sh",
			expected: "docker exec -it app-server sh",
		},
		{
			name:     "Case normalization of prefix",
			buf:      "docker exec ",
			raw:      "Docker exec -it app-server sh",
			expected: "docker exec -it app-server sh",
		},
		{
			name:     "Suffix only completion when buf ends in space",
			buf:      "docker exec ",
			raw:      "-it app-server sh",
			expected: "docker exec -it app-server sh",
		},
		{
			name:     "Suffix only quote completion when buf ends in quote",
			buf:      "git commit -m \"",
			raw:      "fix(auth): resolve bug\"",
			expected: "git commit -m \"fix(auth): resolve bug\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSuggestion(tt.buf, tt.raw)
			if got != tt.expected {
				t.Errorf("NormalizeSuggestion(%q, %q) = %q, want %q", tt.buf, tt.raw, got, tt.expected)
			}
		})
	}
}
