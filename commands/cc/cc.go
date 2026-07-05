package cc

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	cOptions := []spec.Option{
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

	spec.Register(&spec.Spec{
		Name:        "gcc",
		Description: "GNU C compiler",
		Generator:   spec.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	spec.Register(&spec.Spec{
		Name:        "g++",
		Description: "GNU C++ compiler",
		Generator:   spec.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})

	spec.Register(&spec.Spec{
		Name:        "cc",
		Description: "C compiler (alias)",
		Generator:   spec.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	spec.Register(&spec.Spec{
		Name:        "c++",
		Description: "C++ compiler (alias)",
		Generator:   spec.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})

	spec.Register(&spec.Spec{
		Name:        "clang",
		Description: "LLVM C compiler",
		Generator:   spec.FileGenerator(".c", ".h"),
		Options:     cOptions,
	})

	spec.Register(&spec.Spec{
		Name:        "clang++",
		Description: "LLVM C++ compiler",
		Generator:   spec.FileGenerator(".cpp", ".cc", ".cxx", ".h", ".hpp"),
		Options:     cOptions,
	})
}
