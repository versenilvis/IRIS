package ai

import "fmt"

const SystemPrompt = "You are a concise shell command completion assistant. Provide ONLY the completed shell command line. Do not explain, do not use markdown formatting, and do not wrap the command in code blocks or backticks. Always ensure valid shell syntax: if an argument contains spaces or parentheses (such as git commit messages), you MUST wrap that argument in double quotes \"...\"."

func BuildCompletionPrompt(buf string, env EnvSnapshot, dynamicCtx string) string {
	return fmt.Sprintf("Complete this shell command line: %s\nContext:\nCwd: %s\nLastCmd: %s\nLastExitCode: %d\nGitStatus: %s\nRecentCmds: %v\nDynamicContext: %s",
		buf, env.Cwd, env.LastCmd, env.LastExitCode, env.GitStatus, env.RecentCmds, dynamicCtx)
}
