package python

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/commands/core"
)

func pipPackageGenerator(tokens []string, _ string, _ string) []core.Suggestion {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pip", "list", "--format=freeze")
	out, err := cmd.Output()
	if err != nil {
		// try pip3
		cmd = exec.CommandContext(ctx, "pip3", "list", "--format=freeze")
		out, err = cmd.Output()
		if err != nil {
			return nil
		}
	}

	var results []core.Suggestion
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// format is "package==version"
		parts := strings.SplitN(line, "==", 2)
		name := parts[0]
		desc := ""
		if len(parts) == 2 {
			desc = "v" + parts[1]
		}
		results = append(results, core.Suggestion{Cmd: name, Desc: desc})
	}
	return results
}

func makePipSpec(name string) *core.Spec {
	return &core.Spec{
		Name:        name,
		Description: "python packages",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install package", Options: []core.Option{
				{Name: "-r", Description: "from requirements file"},
				{Name: "-U", Description: "upgrade"},
				{Name: "--user", Description: "install for user only"},
				{Name: "--editable", Description: "editable install"},
				{Name: "--index-url", Description: "custom package index"},
				{Name: "--no-deps", Description: "skip dependencies"},
			}},
			{Name: "uninstall", Description: "remove package", Generator: pipPackageGenerator, Options: []core.Option{
				{Name: "-y", Description: "yes to all prompts"},
			}},
			{Name: "freeze", Description: "list installed packages"},
			{Name: "list", Description: "list installed", Options: []core.Option{
				{Name: "--outdated", Description: "show outdated"},
				{Name: "--format", Description: "output format"},
			}},
			{Name: "show", Description: "show package info", Generator: pipPackageGenerator},
			{Name: "download", Description: "download packages"},
			{Name: "check", Description: "verify compatibility"},
			{Name: "config", Description: "manage config"},
			{Name: "search", Description: "search pypi"},
			{Name: "wheel", Description: "build wheel archive"},
			{Name: "hash", Description: "compute hashes"},
			{Name: "completion", Description: "generate completion"},
		},
	}
}

func init() {
	core.Register(makePipSpec("pip"))
	core.Register(makePipSpec("pip3"))
}
