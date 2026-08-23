package main

import "github.com/spf13/cobra"

// minimal Cobra CLI, built at test time as a fixture for TestLooksLikeCobraBinary_RealCobraBinary
func main() {
	root := &cobra.Command{Use: "cobrafixture"}
	_ = root.Execute()
}
