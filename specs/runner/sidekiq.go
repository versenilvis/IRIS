package runner

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "sidekiq",
		Description: "Background job framework for Ruby",
		Options: []core.Option{
			{Name: "--concurrency", Description: "Processor threads to use"},
			{Name: "--environment", Description: "Application environment"},
			{Name: "--tag", Description: "Process tag for procline"},
			{Name: "--queue", Description: "Queues to process with optional weights"},
			{Name: "--require", Description: "Location of Rails application with jobs or file to require"},
			{Name: "--timeout", Description: "Shutdown timeout"},
			{Name: "--verbose", Description: "Print more verbose output"},
			{Name: "--config", Description: "Path to YAML config file"},
			{Name: "--help", Description: "Show help for sidekiq run"},
			{Name: "--version", Description: "Print version and exit"},
		},
	})
}
