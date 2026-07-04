package ops

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "kind",
		Description: "Cluster",
		Subcommands: []core.Subcommand{
			{Name: "Build", Description: "Build one of [node-image]"},
			{Name: "node-image", Description: "Builds a node image"},
			{Name: "architecture", Description: "Architecture to build for, defaults to the host architecture"},
			{Name: "base image", Description: "Name:tag of the base image to use for the build"},
			{Name: "name:tag", Description: "Name:tag of the resulting image to be built"},
			{Name: "path", Description: "Path to the Kubernetes source directory"},
			{Name: "type", Description: "Build type, default is docker"},
			{Name: "create", Description: "Creates a cluster"},
			{Name: "cluster", Description: "Creates a cluster"},
			{Name: "completion", Description: "Generates shell completion scripts"},
			{Name: "bash", Description: "Output shell completions for bash"},
			{Name: "fish", Description: "Output shell completions for fish"},
			{Name: "zsh", Description: "Output shell completions for zsh"},
			{Name: "delete", Description: "Deletes one or more clusters"},
			{Name: "clusters", Description: "Delete Clusters"},
			{Name: "export", Description: "Exports a cluster's kubeconfig"},
			{Name: "kubeconfig", Description: "Exports a cluster's kubeconfig"},
			{Name: "logs", Description: "Exports logs to a tempdir or [output-dir] if specified"},
			{Name: "get", Description: "Gets one of [clusters, nodes, kubeconfig]"},
			{Name: "nodes", Description: "Lists existing kind nodes by their name"},
			{Name: "load", Description: "Loads images into node from an archive or image on host"},
			{Name: "docker-image", Description: "Loads docker images from host into all or specified nodes by name"},
			{Name: "image-archive", Description: "Loads docker images from archive into all or specified nodes by name"},
			{Name: "version", Description: "Prints the kind CLI version"},
		},
		Options: []core.Option{
			{Name: "--arch", Description: "Architecture to build for, defaults to the host architecture"},
			{Name: "--base-image", Description: "Base image to use for the node image"},
			{Name: "--image", Description: "Name:tag of the resulting image to be built"},
			{Name: "--kube-root", Description: "Path to the Kubernetes source directory"},
			{Name: "--type", Description: "Type of node image to build"},
			{Name: "--config", Description: "Path to a kind config file"},
			{Name: "--kubeconfig", Description: "Sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config"},
			{Name: "--name", Description: "Cluster name"},
			{Name: "--retain", Description: "Retain nodes for debugging when cluster creation fails"},
			{Name: "--wait", Description: "Wait for control-plane node to be ready"},
			{Name: "-A", Description: "Delete all clusters"},
			{Name: "--internal", Description: "Use internal address instead of externalt"},
			{Name: "--nodes", Description: "Comma separated list of nodes to load images into"},
			{Name: "-h", Description: "Help for kind"},
			{Name: "-q", Description: "Silence all stderr output"},
			{Name: "-v", Description: "Info log verbosity, higher value produces more output"},
		},
	})
}
