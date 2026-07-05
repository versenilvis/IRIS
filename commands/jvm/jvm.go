package jvm

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "javac",
		Description: "Java compiler",
		Generator:   spec.FileGenerator(".java"),
		Options: []spec.Option{
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

	spec.Register(&spec.Spec{
		Name:        "java",
		Description: "Java runtime",
		Generator:   spec.FileGenerator(".jar", ".class"),
		Options: []spec.Option{
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

	spec.Register(&spec.Spec{
		Name:        "kotlinc",
		Description: "Kotlin compiler",
		Generator:   spec.FileGenerator(".kt", ".kts"),
		Options: []spec.Option{
			{Name: "-d", Description: "output dir or jar"},
			{Name: "-cp", Description: "class path"},
			{Name: "-include-runtime", Description: "include runtime in jar"},
			{Name: "-jvm-target", Description: "JVM target version"},
			{Name: "-script", Description: "run as script"},
		},
	})
}
