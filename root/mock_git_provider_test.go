package root

import "github.com/versenilvis/iris/spec/alias"

type mockGitProvider struct{}

func (m *mockGitProvider) ToolName() string { return "git" }
func (m *mockGitProvider) GetAliases(cwd string) []alias.AliasEntry {
	return []alias.AliasEntry{{Name: "co", Expansion: "checkout", Scope: "local"}}
}
