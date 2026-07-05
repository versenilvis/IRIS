package ops

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/spec"
)

func dockerContainerGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	cwd := spec.GetCWD()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--format={{.Names}}")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		results = append(results, spec.Suggestion{Cmd: line, Desc: "container"})
	}
	return results
}

func dockerRunningContainerGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	cwd := spec.GetCWD()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "--format={{.Names}}")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		results = append(results, spec.Suggestion{Cmd: line, Desc: "running"})
	}
	return results
}

func dockerImageGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	cwd := spec.GetCWD()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "images", "--format={{.Repository}}:{{.Tag}}")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line == "<none>:<none>" || seen[line] {
			continue
		}
		seen[line] = true
		results = append(results, spec.Suggestion{Cmd: line, Desc: "image"})
	}
	return results
}

func init() {
	spec.Register(&spec.Spec{
		Name:        "docker",
		Description: "container engine",
		Subcommands: []spec.Subcommand{
			{
				Name:        "ps",
				Description: "list containers",
				Options: []spec.Option{
					{Name: "-a", Description: "show all"},
					{Name: "-q", Description: "only show IDs"},
					{Name: "--format", Description: "format output"},
					{Name: "--filter", Description: "filter by condition"},
				},
			},
			{
				Name:        "build",
				Description: "build image",
				Generator:   spec.FileGenerator(".dockerfile", ".Dockerfile"),
				Options: []spec.Option{
					{Name: "-t", Description: "tag name"},
					{Name: "-f", Description: "dockerfile path"},
					{Name: "--no-cache", Description: "no build cache"},
					{Name: "--platform", Description: "target platform"},
					{Name: "--build-arg", Description: "build argument"},
				},
			},
			{
				Name:        "run",
				Description: "run container",
				Generator:   dockerImageGenerator,
				Options: []spec.Option{
					{Name: "-d", Description: "detached mode"},
					{Name: "-p", Description: "port mapping"},
					{Name: "-v", Description: "volume mount"},
					{Name: "--rm", Description: "auto remove"},
					{Name: "--name", Description: "container name"},
					{Name: "-it", Description: "interactive tty"},
					{Name: "-e", Description: "set env variable"},
					{Name: "--network", Description: "network mode"},
					{Name: "--restart", Description: "restart policy"},
				},
			},
			{
				Name:        "pull",
				Description: "pull image",
			},
			{
				Name:        "push",
				Description: "push image",
				Generator:   dockerImageGenerator,
			},
			{
				Name:        "exec",
				Description: "exec in container",
				Generator:   dockerRunningContainerGenerator,
				Options: []spec.Option{
					{Name: "-it", Description: "interactive tty"},
					{Name: "-e", Description: "set env variable"},
					{Name: "-u", Description: "run as user"},
				},
			},
			{
				Name:        "stop",
				Description: "stop container",
				Generator:   dockerRunningContainerGenerator,
			},
			{
				Name:        "start",
				Description: "start container",
				Generator:   dockerContainerGenerator,
			},
			{
				Name:        "restart",
				Description: "restart container",
				Generator:   dockerContainerGenerator,
			},
			{
				Name:        "rm",
				Description: "remove container",
				Generator:   dockerContainerGenerator,
				Options: []spec.Option{
					{Name: "-f", Description: "force remove"},
					{Name: "-v", Description: "remove volumes"},
				},
			},
			{
				Name:        "rmi",
				Description: "remove image",
				Generator:   dockerImageGenerator,
				Options: []spec.Option{
					{Name: "-f", Description: "force remove"},
				},
			},
			{
				Name:        "images",
				Description: "list images",
				Options: []spec.Option{
					{Name: "-a", Description: "show all"},
					{Name: "-q", Description: "only IDs"},
					{Name: "--format", Description: "format output"},
				},
			},
			{
				Name:        "logs",
				Description: "view logs",
				Generator:   dockerContainerGenerator,
				Options: []spec.Option{
					{Name: "-f", Description: "follow output"},
					{Name: "--tail", Description: "last n lines"},
					{Name: "--since", Description: "since timestamp"},
					{Name: "-t", Description: "show timestamps"},
				},
			},
			{
				Name:        "inspect",
				Description: "show low-level info",
				Generator: func(tokens []string, prefix string, partial string) []spec.Suggestion {
					containers := dockerContainerGenerator(tokens, prefix, partial)
					images := dockerImageGenerator(tokens, prefix, partial)
					return append(containers, images...)
				},
			},
			{
				Name:        "compose",
				Description: "multi-container",
				Subcommands: []spec.Subcommand{
					{Name: "up", Description: "start services", Options: []spec.Option{{Name: "-d", Description: "detached"}, {Name: "--build", Description: "rebuild images"}}},
					{Name: "down", Description: "stop services", Options: []spec.Option{{Name: "-v", Description: "remove volumes"}}},
					{Name: "build", Description: "build services"},
					{Name: "logs", Description: "view logs", Options: []spec.Option{{Name: "-f", Description: "follow"}}},
					{Name: "ps", Description: "list services"},
					{Name: "exec", Description: "execute command"},
					{Name: "restart", Description: "restart services"},
					{Name: "pull", Description: "pull images"},
					{Name: "stop", Description: "stop services"},
					{Name: "start", Description: "start services"},
					{Name: "config", Description: "validate config", Generator: spec.FileGenerator(".yml", ".yaml")},
				},
				Options: []spec.Option{
					{Name: "-f", Description: "compose file"},
					{Name: "--project-name", Description: "project name"},
				},
			},
			{
				Name:        "network",
				Description: "manage networks",
				Subcommands: []spec.Subcommand{
					{Name: "ls", Description: "list networks"},
					{Name: "create", Description: "create network"},
					{Name: "rm", Description: "remove network"},
					{Name: "inspect", Description: "show details"},
					{Name: "prune", Description: "remove unused"},
				},
			},
			{
				Name:        "volume",
				Description: "manage volumes",
				Subcommands: []spec.Subcommand{
					{Name: "ls", Description: "list volumes"},
					{Name: "create", Description: "create volume"},
					{Name: "rm", Description: "remove volume"},
					{Name: "prune", Description: "remove unused"},
					{Name: "inspect", Description: "show details"},
				},
			},
			{
				Name:        "system",
				Description: "manage docker system",
				Subcommands: []spec.Subcommand{
					{Name: "prune", Description: "remove unused data", Options: []spec.Option{{Name: "-a", Description: "remove all unused"}, {Name: "--volumes", Description: "include volumes"}}},
					{Name: "df", Description: "show disk usage"},
					{Name: "info", Description: "system info"},
				},
			},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "docker-compose",
		Description: "multi-container (legacy)",
		Subcommands: []spec.Subcommand{
			{Name: "up", Description: "start services", Options: []spec.Option{{Name: "-d", Description: "detached"}, {Name: "--build", Description: "rebuild"}}},
			{Name: "down", Description: "stop services", Options: []spec.Option{{Name: "-v", Description: "remove volumes"}}},
			{Name: "build", Description: "build services"},
			{Name: "logs", Description: "view logs", Options: []spec.Option{{Name: "-f", Description: "follow"}}},
			{Name: "ps", Description: "list services"},
			{Name: "exec", Description: "execute command"},
			{Name: "restart", Description: "restart services"},
			{Name: "pull", Description: "pull images"},
		},
		Options: []spec.Option{
			{Name: "-f", Description: "compose file"},
		},
	})
}
