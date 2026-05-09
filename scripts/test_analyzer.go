//go:build ignore
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TestEvent represents a single line of output from 'go test -json'
type TestEvent struct {
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Action  string  `json:"Action"` // pass, fail, skip, run...
	Output  string  `json:"Output"`
}

func main() {
	fmt.Println("\033[1;34m[IRIS] Automated Project Analysis & Test Reporter\033[0m")
	fmt.Println("\033[2mRunning all tests and analyzing results...\033[0m")

	// Run go test with -json flag
	cmd := exec.Command("go", "test", "./...", "-json")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	// Parse JSON stream
	decoder := json.NewDecoder(stdout)
	
	type result struct {
		name     string
		category string
		status   string
	}
	results := []result{}
	
	for {
		var event TestEvent
		if err := decoder.Decode(&event); err != nil {
			break
		}

		// We only care about the final status of individual tests
		if event.Test != "" && (event.Action == "pass" || event.Action == "fail" || event.Action == "skip") {
			pkgParts := strings.Split(event.Package, "/")
			category := pkgParts[len(pkgParts)-1]
			
			results = append(results, result{
				name:     event.Test,
				category: category,
				status:   event.Action,
			})
		}
	}
	cmd.Wait()

	fmt.Println("\n==================================================================================")
	fmt.Printf("%-8s | %-12s | %-45s\n", "STATUS", "CATEGORY", "TEST CASE")
	fmt.Println("==================================================================================")

	overallPass := true
	for _, res := range results {
		status := ""
		switch res.status {
		case "pass":
			status = "\033[32mPASS\033[0m"
		case "fail":
			status = "\033[31mFAIL\033[0m"
			overallPass = false
		case "skip":
			status = "\033[33mSKIP\033[0m"
		}
		
		fmt.Printf("%-18s | %-12s | %-45s\n", status, res.category, res.name)
	}

	fmt.Println("==================================================================================")
	if len(results) == 0 {
		fmt.Println("\033[1;33mWARNING: No tests found! Did you write any TestXxx functions?\033[0m")
	} else if overallPass {
		fmt.Printf("\033[1;32mANALYSIS SUCCESS: All %d tests passed successfully!\033[0m\n", len(results))
	} else {
		fmt.Println("\033[1;31mANALYSIS FAILED: Some tests did not pass. Run 'just test' for details.\033[0m")
		os.Exit(1)
	}
}