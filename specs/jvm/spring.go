package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "spring",
		Description: "Initialize a new project using Spring Initializr",
		Subcommands: []core.Subcommand{
			{Name: "init", Description: "Initialize a new project using Spring Initializr"},
			{Name: "gradle-build", Description: "Generate a Gradle build file"},
			{Name: "gradle-project", Description: "Generate a Gradle based project archive using the Groovy DSL"},
			{Name: "gradle-project-kotlin", Description: "Generate a Gradle based project archive using the Kotlin DSL"},
			{Name: "maven-build", Description: "Generate a Maven pom.xml"},
			{Name: "maven-project", Description: "Generate a Maven based project archive"},
			{Name: "encodepassword", Description: "Encode a password for use with Spring Security"},
			{Name: "shell", Description: "Start a nested shell"},
			{Name: "help", Description: "Show help for other commands"},
		},
		Options: []core.Option{
			{Name: "-a", Description: "Project coordinates"},
			{Name: "-b", Description: "Spring Boot version"},
			{Name: "--build", Description: "Build system to use"},
			{Name: "-d", Description: "Project description"},
			{Name: "-f", Description: "Force overwrite of existing files"},
			{Name: "--format", Description: "Format of the generated content"},
			{Name: "-g", Description: "Project coordinates"},
			{Name: "-j", Description: "Language level"},
			{Name: "--list", Description: "List the capabilities of the service"},
			{Name: "-n", Description: "Project name"},
			{Name: "-p", Description: "Project packaging"},
			{Name: "--package-name", Description: "Package name"},
			{Name: "-t", Description: "Project type"},
			{Name: "--target", Description: "URL of the service to use"},
			{Name: "-v", Description: "Project version"},
			{Name: "-x", Description: "Encode a password for use with Spring Security"},
			{Name: "--help", Description: "Show help for spring"},
			{Name: "--version", Description: "Get spring CLI version"},
			{Name: "--debug", Description: "Print additional status information for the command you are running"},
		},
	})
}
