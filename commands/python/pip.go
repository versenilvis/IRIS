package python

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/versenilvis/iris/spec"
)

func pipPackageGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
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

	var results []spec.Suggestion
	for line := range strings.SplitSeq(string(out), "\n") {
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
		results = append(results, spec.Suggestion{Cmd: name, Desc: desc})
	}
	return results
}

func makePipSpec(name string) *spec.Spec {
	return &spec.Spec{
		Name:        name,
		Description: "python packages",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "install package", Options: []spec.Option{
				{Name: "-r", Description: "from requirements file"},
				{Name: "-U", Description: "upgrade"},
				{Name: "--user", Description: "install for user only"},
				{Name: "--editable", Description: "editable install"},
				{Name: "--index-url", Description: "custom package index"},
				{Name: "--no-deps", Description: "skip dependencies"},
			}},
			{Name: "uninstall", Description: "remove package", Generator: pipPackageGenerator, Options: []spec.Option{
				{Name: "-y", Description: "yes to all prompts"},
			}},
			{Name: "freeze", Description: "list installed packages"},
			{Name: "list", Description: "list installed", Options: []spec.Option{
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
	spec.Register(makePipSpec("pip"))
	spec.Register(makePipSpec("pip3"))
}
