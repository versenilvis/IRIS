package scoring

import (
	"strings"

	"github.com/versenilvis/iris/internal/workspace"
)

type ContextRule interface {
	Match(ws workspace.WorkspaceInfo, cmd string) bool
	Bonus() int
}

type SimpleContextRule struct {
	check func(ws workspace.WorkspaceInfo, cmd string) bool
	bonus int
}

func (r *SimpleContextRule) Match(ws workspace.WorkspaceInfo, cmd string) bool {
	return r.check(ws, cmd)
}

func (r *SimpleContextRule) Bonus() int {
	return r.bonus
}

var DefaultContextRules = []ContextRule{
	// Git rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasGit && (strings.HasPrefix(cmd, "git status") || strings.HasPrefix(cmd, "git diff") ||
				strings.HasPrefix(cmd, "git add") || strings.HasPrefix(cmd, "git push") ||
				strings.HasPrefix(cmd, "git pull") || strings.HasPrefix(cmd, "git commit") ||
				strings.HasPrefix(cmd, "git switch") || strings.HasPrefix(cmd, "git checkout") ||
				strings.HasPrefix(cmd, "git branch"))
		},
		bonus: 40,
	},
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasGit && (strings.HasPrefix(cmd, "git init") || strings.HasPrefix(cmd, "git clone"))
		},
		bonus: -50,
	},

	// Node.js & Bun rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasNodeProject && (strings.HasPrefix(cmd, "npm run ") || strings.HasPrefix(cmd, "pnpm run ") ||
				strings.HasPrefix(cmd, "yarn run ") || strings.HasPrefix(cmd, "bun run ") ||
				strings.HasPrefix(cmd, "npm test") || strings.HasPrefix(cmd, "npm start") ||
				strings.HasPrefix(cmd, "bun test") || strings.HasPrefix(cmd, "bun start"))
		},
		bonus: 50,
	},
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasNodeProject && (strings.HasPrefix(cmd, "npm install") || strings.HasPrefix(cmd, "npm i ") ||
				strings.HasPrefix(cmd, "pnpm add") || strings.HasPrefix(cmd, "yarn add") ||
				strings.HasPrefix(cmd, "bun install") || strings.HasPrefix(cmd, "bun add"))
		},
		bonus: 40,
	},

	// Go rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasGoProject && (strings.HasPrefix(cmd, "go test") || strings.HasPrefix(cmd, "go run") ||
				strings.HasPrefix(cmd, "go build") || strings.HasPrefix(cmd, "go mod tidy") ||
				strings.HasPrefix(cmd, "go vet") || strings.HasPrefix(cmd, "go fmt"))
		},
		bonus: 50,
	},

	// Rust rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasRustProject && (strings.HasPrefix(cmd, "cargo test") || strings.HasPrefix(cmd, "cargo run") ||
				strings.HasPrefix(cmd, "cargo build") || strings.HasPrefix(cmd, "cargo check") ||
				strings.HasPrefix(cmd, "cargo clippy"))
		},
		bonus: 50,
	},

	// Python rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasPythonProject && (strings.HasPrefix(cmd, "pytest") || strings.HasPrefix(cmd, "python main.py") ||
				strings.HasPrefix(cmd, "pip install") || strings.HasPrefix(cmd, "poetry run ") ||
				strings.HasPrefix(cmd, "uv run "))
		},
		bonus: 50,
	},

	// Justfile rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasJustfile && strings.HasPrefix(cmd, "just ")
		},
		bonus: 50,
	},

	// Makefile & C/C++ rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasMakefile && strings.HasPrefix(cmd, "make")
		},
		bonus: 50,
	},
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasMakefile && (strings.HasPrefix(cmd, "gcc ") || strings.HasPrefix(cmd, "g++ ") ||
				strings.HasPrefix(cmd, "clang ") || strings.HasPrefix(cmd, "cmake "))
		},
		bonus: 40,
	},

	// Docker & K8s rules
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasDockerfile && (strings.HasPrefix(cmd, "docker build") || strings.HasPrefix(cmd, "docker compose up") ||
				strings.HasPrefix(cmd, "docker compose down") || strings.HasPrefix(cmd, "docker compose logs"))
		},
		bonus: 40,
	},
	&SimpleContextRule{
		check: func(ws workspace.WorkspaceInfo, cmd string) bool {
			return ws.HasK8s && (strings.HasPrefix(cmd, "kubectl get ") || strings.HasPrefix(cmd, "kubectl apply -f ") ||
				strings.HasPrefix(cmd, "kubectl logs ") || strings.HasPrefix(cmd, "kubectl describe ") ||
				strings.HasPrefix(cmd, "helm upgrade ") || strings.HasPrefix(cmd, "helm install "))
		},
		bonus: 40,
	},
}

func ApplyContextRules(ws workspace.WorkspaceInfo, cmd string) int {
	return ApplyCustomContextRules(ws, cmd, DefaultContextRules)
}

func ApplyCustomContextRules(ws workspace.WorkspaceInfo, cmd string, rules []ContextRule) int {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return 0
	}

	total := 0
	for _, rule := range rules {
		if rule.Match(ws, cmd) {
			total += rule.Bonus()
		}
	}

	if total > 100 {
		return 100
	}
	if total < -100 {
		return -100
	}
	return total
}
