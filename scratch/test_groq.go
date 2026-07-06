package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/versenilvis/iris/ai"
	"github.com/versenilvis/iris/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	config.Init(cfg)

	fmt.Printf("AI Enabled: %v\n", cfg.AI.Enabled)
	fmt.Printf("Active Provider: %q\n", cfg.AI.Provider)

	pCfg, ok := cfg.AI.GetActiveProvider()
	if !ok {
		fmt.Println("Error: active provider not found in config")
		os.Exit(1)
	}

	fmt.Printf("Endpoint: %s\n", pCfg.Endpoint)
	fmt.Printf("Model: %s\n", pCfg.Model)
	fmt.Printf("API Key Present: %v (len %d)\n", pCfg.GetAPIKey() != "", len(pCfg.GetAPIKey()))

	client, err := ai.NewClient(pCfg)
	if err != nil {
		fmt.Printf("Error creating AI client: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	env := ai.EnvSnapshot{
		Cwd:          "/home/verse/dev/github/iris",
		LastCmd:      "git status",
		LastExitCode: 0,
		GitStatus:    "modified: ai/client.go, ai/engine.go",
	}

	fmt.Println("\nSending live test completion request to Groq for: \"git commit -m \" ...")
	start := time.Now()
	sugg, err := client.Suggest(ctx, "git commit -m \"", env, "")
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("AI Suggest Error: %v\n", err)
		os.Exit(1)
	}
	if sugg == nil {
		fmt.Println("AI returned no suggestion (nil)")
		os.Exit(0)
	}

	fmt.Printf("SUCCESS! Response received in %v:\n", duration)
	fmt.Printf("  Cmd:        %s\n", sugg.Cmd)
	fmt.Printf("  Confidence: %d\n", sugg.Confidence)
	fmt.Printf("  Source:     %s\n", sugg.Source)
}
