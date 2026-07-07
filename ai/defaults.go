package ai

import (
	"strings"

	"github.com/versenilvis/iris/spec"
)

var DefaultProviders = []*CommandContextProvider{
	{NameStr: "docker_exec", Prefixes: []string{"docker exec", "docker logs", "docker stop", "docker restart", "docker rm"},
		GatherCmd: []string{"docker", "ps", "--format", "{{.Names}}\t{{.Image}}"}, Label: "Running containers"},
	{NameStr: "docker_compose", Prefixes: []string{"docker compose exec", "docker compose logs", "docker-compose exec", "docker-compose logs"},
		GatherCmd: []string{"docker", "compose", "ps", "--format", "{{.Name}}\t{{.Service}}"}, Label: "Compose services"},
	{NameStr: "kubectl_pods", Prefixes: []string{"kubectl exec", "kubectl logs", "kubectl describe pod", "kubectl delete pod"},
		GatherCmd: []string{"kubectl", "get", "pods", "--no-headers"}, Label: "Pods"},
	{NameStr: "git_branch", Prefixes: []string{"git checkout", "git switch", "git merge", "git rebase", "git branch -d", "git branch -D"},
		GatherCmd: []string{"git", "branch", "-a", "--format=%(refname:short)"}, Label: "Branches"},
	{NameStr: "kill_proc", Prefixes: []string{"kill ", "kill -9 "},
		GatherCmd: []string{"ps", "-eo", "pid,comm,%cpu,%mem", "--sort=-%cpu"}, Label: "Top processes"},
	{NameStr: "systemctl", Prefixes: []string{"systemctl restart", "systemctl stop", "systemctl status"},
		GatherCmd: []string{"systemctl", "list-units", "--type=service", "--no-legend"}, Label: "Services"},
}

var DefaultEmptyLineRules = []EmptyLineRule{
	{Name: "merge_in_progress", Match: func(e EnvSnapshot) bool { return e.GitMergeInProgress },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: "git commit", Desc: "finish merge", Icon: "git", Source: string(SourceSpec), Confidence: 85}
		}},
	{Name: "rebase_in_progress", Match: func(e EnvSnapshot) bool { return e.GitRebaseInProgress },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: "git rebase --continue", Desc: "continue rebase", Icon: "git", Source: string(SourceSpec), Confidence: 85}
		}},
	{Name: "retry_failed", Match: func(e EnvSnapshot) bool { return e.LastExitCode != 0 && e.LastCmd != "" },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: e.LastCmd, Desc: "retry failed command", Icon: "retry", Source: string(SourceSpec), Confidence: 80}
		}},
	{Name: "git_status_diff", Match: func(e EnvSnapshot) bool { return strings.TrimSpace(e.LastCmd) == "git status" },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: "git diff", Desc: "view modifications", Icon: "git", Source: string(SourceSpec), Confidence: 75}
		}},
	{Name: "git_dirty_status", Match: func(e EnvSnapshot) bool { return e.GitStatus != "" },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: "git status", Desc: "check repository state", Icon: "git", Source: string(SourceSpec), Confidence: 70}
		}},
	{Name: "npm_run_dev", Match: func(e EnvSnapshot) bool { return strings.Contains(e.DirSignature, "package.json") },
		Suggest: func(e EnvSnapshot) *spec.Suggestion {
			return &spec.Suggestion{Cmd: "npm run dev", Desc: "start dev server", Icon: "npm", Source: string(SourceSpec), Confidence: 65}
		}},
}
