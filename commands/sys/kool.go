package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "kool",
		Description: "Script",
		Subcommands: []spec.Subcommand{
			{Name: "create", Description: "Create a new project using a preset"},
			{Name: "preset", Description: "Preset that will be used to create the project"},
			{Name: "folder", Description: "Folder where the project will be created"},
			{Name: "deploy", Description: "Deploy a local application to a Kool Cloud environment"},
			{Name: "destroy", Description: "Destroy an environment deployed to Kool Cloud"},
			{Name: "exec", Description: "Execute a command inside a running service container deployed to Kool Cloud"},
			{Name: "logs", Description: "See the logs of running service container deployed to Kool Cloud"},
			{Name: "docker", Description: "Create a new container (a powered up 'docker run')"},
			{Name: "image", Description: "Docker image"},
			{Name: "command", Description: "Command to execute inside the container"},
			{Name: "service", Description: "Service you want to execute a command"},
			{Name: "info", Description: "Print out information about the local environment"},
			{Name: "recipe", Description: "Adds configuration for some recipe in the current work directory"},
			{Name: "run", Description: "Execute a script defined in kool.yml"},
			{Name: "script", Description: "Script to be executed"},
			{Name: "self-update", Description: "Update kool to the latest version"},
			{Name: "share", Description: "Live share your local environment on the Internet using an HTTP tunnel"},
			{Name: "start", Description: "Start service containers defined in docker-compose.yml"},
			{Name: "status", Description: "Show the status of all service containers"},
			{Name: "stop", Description: "Stop and destroy running service containers"},
		},
		Options: []spec.Option{
			{Name: "--verbose", Description: "Increases output verbosity"},
			{Name: "--help", Description: "Help for create"},
			{Name: "--container", Description: "Container target"},
			{Name: "--follow", Description: "Follow log output"},
			{Name: "--tail", Description: "Help for logs"},
			{Name: "--env", Description: "Environment variables"},
			{Name: "--network", Description: "Connect a container to a network"},
			{Name: "--publish", Description: "Publish a container's port(s) to the host"},
			{Name: "--volume", Description: "Bind mount a volume"},
			{Name: "--detach", Description: "Detached mode: Run command in the background"},
			{Name: "--purge", Description: "Remove all persistent data from volume mounts on containers"},
			{Name: "--rebuild", Description: "Updates and builds service's images"},
			{Name: "--port", Description: "The name of the local service container you want to share"},
			{Name: "--subdomain", Description: "The subdomain used to generate your public https://subdomain.kool.live URL"},
			{Name: "--foreground", Description: "Start containers in foreground mode"},
			{Name: "--profile", Description: "Specify a profile to enable"},
		},
	})
}
