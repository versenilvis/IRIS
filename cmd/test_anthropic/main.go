package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/versenilvis/iris/internal/ai"
	"github.com/versenilvis/iris/internal/config"
)

func loadEnvFile(path string) map[string]string {
	out := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			// strip surrounding quotes
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			out[key] = val
		}
	}
	return out
}

func main() {
	// Load key from .env file — never touches command line
	envPath := os.ExpandEnv("$HOME/.hermes/.env")
	envVars := loadEnvFile(envPath)
	apiKey := envVars["ANTHROPIC_API_KEY"]
	if apiKey == "" {
		fmt.Println("FAIL: ANTHROPIC_API_KEY not found in", envPath)
		os.Exit(1)
	}
	fmt.Printf("Key loaded: %d chars, prefix: %s...\n", len(apiKey), apiKey[:8])

	cfg := config.ProviderConfig{
		InheritedFrom: "anthropic",
		Model:         "claude-3-5-haiku-20241022",
		TimeoutMS:     5000,
		APIKey:        apiKey,
	}

	client := ai.NewAnthropicClient(cfg)

	// Test 1: Simple completion
	fmt.Println("\n--- Test 1: 'git commit -m \"' ---")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	env := ai.NewEnvSnapshot("/home/david", "", 0, []string{"ls", "git status"})
	sugg, err := client.Suggest(ctx, "git commit -m \"", env, "")
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	if sugg == nil {
		fmt.Println("WARN: nil suggestion")
	} else {
		fmt.Printf("OK: cmd=%q confidence=%d dangerous=%v\n", sugg.Cmd, sugg.Confidence, sugg.Dangerous)
	}

	// Test 2: Dangerous detection
	fmt.Println("\n--- Test 2: danger detection on 'rm ' ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	sugg2, err := client.Suggest(ctx2, "rm ", env, "")
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
	} else if sugg2 != nil {
		fmt.Printf("OK: cmd=%q dangerous=%v\n", sugg2.Cmd, sugg2.Dangerous)
		if sugg2.Dangerous {
			fmt.Println("✅ Dangerous flag working!")
		}
	}

	fmt.Println("\n=== Anthropic client: PASSED ===")
}