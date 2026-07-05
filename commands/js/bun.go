package js

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "bun",
		Description: "bun js runtime",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "install packages", Options: []spec.Option{
				{Name: "--frozen-lockfile", Description: "no lockfile update"},
				{Name: "--production", Description: "production only"},
				{Name: "--no-save", Description: "don't update package.json"},
				{Name: "--dry-run", Description: "print what would be installed"},
			}},
			{Name: "add", Description: "add package", Options: []spec.Option{
				{Name: "-d", Description: "dev dependency"},
				{Name: "--dev", Description: "dev dependency (long form)"},
				{Name: "--optional", Description: "optional dependency"},
				{Name: "--peer", Description: "peer dependency"},
				{Name: "-g", Description: "global"},
				{Name: "--exact", Description: "exact version"},
				{Name: "--no-save", Description: "don't update package.json"},
			}},
			{Name: "remove", Description: "remove package", Options: []spec.Option{
				{Name: "-g", Description: "remove global package"},
			}},
			{Name: "run", Description: "run script", Generator: NpmScriptGenerator, Options: []spec.Option{
				{Name: "--silent", Description: "suppress output"},
				{Name: "--bun", Description: "force bun runtime"},
				{Name: "--watch", Description: "watch mode"},
				{Name: "--hot", Description: "hot reload"},
				{Name: "--smol", Description: "reduce memory usage"},
			}},
			{Name: "build", Description: "bundle files", Generator: spec.FileGenerator(".ts", ".tsx", ".js", ".jsx"), Options: []spec.Option{
				{Name: "--outdir", Description: "output directory"},
				{Name: "--outfile", Description: "output file"},
				{Name: "--minify", Description: "enable all minification"},
				{Name: "--minify-syntax", Description: "minify syntax only"},
				{Name: "--minify-whitespace", Description: "minify whitespace only"},
				{Name: "--minify-identifiers", Description: "minify identifiers only"},
				{Name: "--sourcemap", Description: "generate sourcemap (none/inline/external)"},
				{Name: "--target", Description: "target environment (browser/bun/node)"},
				{Name: "--format", Description: "output format (esm/cjs/iife)"},
				{Name: "--splitting", Description: "enable code splitting"},
				{Name: "--entry-naming", Description: "entry file naming pattern"},
				{Name: "--chunk-naming", Description: "chunk file naming pattern"},
				{Name: "--asset-naming", Description: "asset file naming pattern"},
				{Name: "--watch", Description: "watch mode"},
				{Name: "--external", Description: "mark package as external"},
				{Name: "--define", Description: "replace global identifiers"},
				{Name: "--loader", Description: "set file loader"},
				{Name: "--public-path", Description: "prefix for public assets"},
			}},
			{Name: "test", Description: "run tests", Options: []spec.Option{
				{Name: "--watch", Description: "watch mode"},
				{Name: "-t", Description: "filter by test name"},
				{Name: "--coverage", Description: "collect coverage"},
				{Name: "--bail", Description: "stop after first failure"},
				{Name: "--timeout", Description: "test timeout in ms"},
				{Name: "--rerun-each", Description: "rerun each test n times"},
				{Name: "--only", Description: "run test.only tests"},
				{Name: "--todo", Description: "show todo tests"},
				{Name: "--reporter", Description: "test reporter (default/junit/tap)"},
			}},
			{Name: "x", Description: "execute package (bunx)", Options: []spec.Option{
				{Name: "--bun", Description: "force bun runtime"},
				{Name: "--silent", Description: "suppress output"},
			}},
			{Name: "outdated", Description: "check outdated packages"},
			{Name: "patch", Description: "patch a package", Options: []spec.Option{
				{Name: "--commit", Description: "commit the patch"},
			}},
			{Name: "publish", Description: "publish package to npm", Options: []spec.Option{
				{Name: "--dry-run", Description: "simulate publish"},
				{Name: "--tag", Description: "dist-tag"},
				{Name: "--access", Description: "public or restricted"},
			}},
			{Name: "init", Description: "create new project", Options: []spec.Option{
				{Name: "-y", Description: "skip prompts"},
			}},
			{Name: "create", Description: "create from template"},
			{Name: "repl", Description: "launch bun repl"},
			{Name: "upgrade", Description: "upgrade bun"},
			{Name: "link", Description: "link package"},
			{Name: "unlink", Description: "unlink package"},
			{Name: "pm", Description: "package manager helpers", Subcommands: []spec.Subcommand{
				{Name: "cache", Description: "show cache path"},
				{Name: "ls", Description: "list packages"},
				{Name: "hash", Description: "print lockfile hash"},
				{Name: "hash-print", Description: "print resolved lockfile hash"},
				{Name: "trust", Description: "trust package scripts"},
				{Name: "untrusted", Description: "show untrusted packages"},
			}},
		},
		Options: []spec.Option{
			{Name: "--watch", Description: "watch for changes"},
			{Name: "--hot", Description: "hot module reload"},
			{Name: "--smol", Description: "reduce memory usage"},
			{Name: "--silent", Description: "suppress output"},
			{Name: "--bun", Description: "force bun runtime"},
			{Name: "--version", Description: "print version"},
		},
	})

	spec.Register(&spec.Spec{
		Name:        "bunx",
		Description: "execute package (bun x)",
		Options: []spec.Option{
			{Name: "--bun", Description: "force bun runtime"},
			{Name: "--silent", Description: "suppress output"},
		},
	})
}
