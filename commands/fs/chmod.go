package fs

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "chmod",
		Description: "change file permissions",
		MaxArgs:     2, // 1 option + 1 file
		Subcommands: []spec.Subcommand{},
		Options: []spec.Option{
			{Name: "-R", Description: "apply recursively to directories"},
			{Name: "--recursive", Description: "apply recursively to directories"},
			{Name: "-v", Description: "verbose: show each file processed"},
			{Name: "--verbose", Description: "verbose: show each file processed"},
			{Name: "-c", Description: "like verbose but report only when a change is made"},
			{Name: "--changes", Description: "like verbose but report only when a change is made"},
			{Name: "-f", Description: "suppress most error messages"},
			{Name: "--silent", Description: "suppress most error messages"},
			{Name: "--quiet", Description: "suppress most error messages"},
			{Name: "--reference", Description: "use another file's mode as reference (--reference=FILE)"},
			{Name: "-H", Description: "follow symlinks on command line (use with -R)"},
			{Name: "-L", Description: "follow all symlinks (use with -R)"},
			{Name: "-P", Description: "do not follow any symlinks (use with -R)"},
		},
		Generator: modeGenerator(),
	})
}

// modeGenerator suggests common chmod mode strings as the first argument,
// then falls back to file suggestions for subsequent arguments
func modeGenerator() spec.GeneratorFunc {
	filegen := spec.FileGenerator(".sh", ".py", ".bin", ".run")

	commonModes := []struct {
		mode string
		desc string
	}{
		{"+x", "add execute for all"},
		{"-x", "remove execute for all"},
		{"u+x", "add execute for owner"},
		{"a+x", "add execute for all (explicit)"},

		{"u+rw", "owner read+write"},
		{"go-w", "remove write from group+others"},
		{"a-w", "remove write from everyone"},
		{"u+s", "set SUID bit"},
		{"g+s", "set SGID bit"},
		{"+t", "set sticky bit"},

		{"755", "rwxr-xr-x (executable/dir, world-readable)"},
		{"644", "rw-r--r-- (regular file, world-readable)"},
		{"600", "rw------- (private file)"},
		{"700", "rwx------ (private executable)"},
		{"777", "rwxrwxrwx (full permissions for all)"},
		{"666", "rw-rw-rw- (read+write for all)"},
		{"400", "r-------- (read-only for owner)"},
		{"444", "r--r--r-- (read-only for all)"},
		{"750", "rwxr-x--- (group can read+execute)"},
		{"640", "rw-r----- (group can read)"},
		{"664", "rw-rw-r-- (group can read+write)"},
	}

	return func(tokens []string, prefix string, partial string) []spec.Suggestion {
		argCount := 0
		for i := 1; i < len(tokens); i++ {
			t := tokens[i]
			if t != "" && t[0] != '-' {
				argCount++
			}
		}

		if argCount == 0 {
			var results []spec.Suggestion
			for _, m := range commonModes {
				if partial == "" || spec.HasPrefixCI(m.mode, partial) {
					results = append(results, spec.Suggestion{
						Cmd:  m.mode,
						Desc: m.desc,
					})
				}
			}
			return results
		}
		return filegen(tokens, prefix, partial)
	}
}
