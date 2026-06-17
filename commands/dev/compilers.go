package dev

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
	// gcc / g++
	cOptions := []core.Option{
		{Name: "-o", Description: "output file"},
		{Name: "-c", Description: "compile only (no link)"},
		{Name: "-g", Description: "debug symbols"},
		{Name: "-O0", Description: "no optimization"},
		{Name: "-O1", Description: "basic optimization"},
		{Name: "-O2", Description: "standard optimization"},
		{Name: "-O3", Description: "aggressive optimization"},
		{Name: "-Os", Description: "optimize for size"},
		{Name: "-Wall", Description: "enable common warnings"},
		{Name: "-Wextra", Description: "extra warnings"},
		{Name: "-Werror", Description: "treat warnings as errors"},
		{Name: "-std", Description: "language standard (e.g. c11, c++17)"},
		{Name: "-I", Description: "include directory"},
		{Name: "-L", Description: "library directory"},
		{Name: "-l", Description: "link library"},
		{Name: "-D", Description: "define macro"},
		{Name: "-shared", Description: "build shared library"},
		{Name: "-static", Description: "static linking"},
		{Name: "-fPIC", Description: "position independent code"},
		{Name: "-march", Description: "target architecture"},
		{Name: "-pthread", Description: "POSIX threads"},
		{Name: "-v", Description: "verbose"},
		{Name: "-E", Description: "preprocess only"},
		{Name: "-S", Description: "compile to assembly"},
		{Name: "-M", Description: "output dependencies"},
		{Name: "-MMD", Description: "write dep file"},
		{Name: "--version", Description: "print version"},
	}

	core.Register(&core.Spec{
		Name:        "gcc",
		Description: "GNU C compiler",
		Generator:   core.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	core.Register(&core.Spec{
		Name:        "g++",
		Description: "GNU C++ compiler",
		Generator:   core.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})

	core.Register(&core.Spec{
		Name:        "cc",
		Description: "C compiler (alias)",
		Generator:   core.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	core.Register(&core.Spec{
		Name:        "c++",
		Description: "C++ compiler (alias)",
		Generator:   core.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})

	// rustc — direct Rust compiler (cargo is already registered separately)
	core.Register(&core.Spec{
		Name:        "rustc",
		Description: "Rust compiler",
		Generator:   core.FileGenerator(".rs"),
		Options: []core.Option{
			{Name: "-o", Description: "output file"},
			{Name: "--edition", Description: "Rust edition (2015/2018/2021/2024)"},
			{Name: "--crate-type", Description: "crate type (bin/lib/dylib/cdylib/rlib)"},
			{Name: "--crate-name", Description: "crate name"},
			{Name: "-C", Description: "codegen option"},
			{Name: "-L", Description: "add to library search path"},
			{Name: "--extern", Description: "extern crate"},
			{Name: "-g", Description: "debug symbols"},
			{Name: "-O", Description: "optimize"},
			{Name: "--emit", Description: "output type (asm/llvm-ir/obj/mir)"},
			{Name: "--explain", Description: "explain error code"},
			{Name: "--test", Description: "build test harness"},
			{Name: "-W", Description: "set lint warning"},
			{Name: "-A", Description: "set lint allow"},
			{Name: "-D", Description: "set lint deny"},
			{Name: "--edition", Description: "Rust edition"},
			{Name: "--verbose", Description: "verbose"},
			{Name: "--version", Description: "print version"},
		},
	})

	// Java
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

	// Kotlin
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

	// clang
	core.Register(&core.Spec{
		Name:        "clang",
		Description: "LLVM C compiler",
		Generator:   core.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	core.Register(&core.Spec{
		Name:        "clang++",
		Description: "LLVM C++ compiler",
		Generator:   core.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})
}
