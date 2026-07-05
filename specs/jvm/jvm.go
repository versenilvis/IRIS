package jvm

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	core.Register(&core.Spec{
		Name:        "javac",
		Description: "Java compiler",
		Generator:   core.FileGenerator(".java"),
		Options: []core.Option{
			{Name: "-d", Description: "output directory"},
			{Name: "-cp", Description: "class path"},
			{Name: "-classpath", Description: "class path (long form)"},
			{Name: "-sourcepath", Description: "source path"},
			{Name: "-g", Description: "debug info"},
			{Name: "-verbose", Description: "verbose"},
			{Name: "-encoding", Description: "source file encoding"},
			{Name: "--release", Description: "target release version"},
			{Name: "-source", Description: "source version"},
			{Name: "-target", Description: "target version"},
			{Name: "-Xlint", Description: "enable warnings"},
		},
	})

	core.Register(&core.Spec{
		Name:        "java",
		Description: "Java runtime",
		Generator:   core.FileGenerator(".jar", ".class"),
		Options: []core.Option{
			{Name: "-jar", Description: "run jar file"},
			{Name: "-cp", Description: "class path"},
			{Name: "-Xmx", Description: "max heap size (e.g. -Xmx512m)"},
			{Name: "-Xms", Description: "initial heap size"},
			{Name: "-D", Description: "set system property"},
			{Name: "-ea", Description: "enable assertions"},
			{Name: "--version", Description: "print version"},
			{Name: "-server", Description: "server VM"},
			{Name: "-verbose", Description: "verbose class loading"},
		},
	})

	core.Register(&core.Spec{
		Name:        "kotlinc",
		Description: "Kotlin compiler",
		Generator:   core.FileGenerator(".kt", ".kts"),
		Options: []core.Option{
			{Name: "-d", Description: "output dir or jar"},
			{Name: "-cp", Description: "class path"},
			{Name: "-include-runtime", Description: "include runtime in jar"},
			{Name: "-jvm-target", Description: "JVM target version"},
			{Name: "-script", Description: "run as script"},
		},
	})
}
