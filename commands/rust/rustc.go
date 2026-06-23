package rust

import (
	"github.com/versenilvis/iris/commands/core"
)

func init() {
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
}
