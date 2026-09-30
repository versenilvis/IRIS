package ctxcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNode_NoPackageJson(t *testing.T) {
	tmpDir := t.TempDir()
	if v := ValidateNode([]string{"npm", "run", "dev"}, tmpDir); v != Invalid {
		t.Fatalf("expected Invalid when package.json missing, got %v", v)
	}
	if v := ValidateNode([]string{"pnpm", "dev"}, tmpDir); v != Invalid {
		t.Fatalf("expected Invalid for pnpm dev when package.json missing, got %v", v)
	}
	// built-ins are Free even without package.json
	if v := ValidateNode([]string{"npm", "install"}, tmpDir); v != Free {
		t.Fatalf("expected Free for npm install, got %v", v)
	}
	if v := ValidateNode([]string{"npx", "create-react-app"}, tmpDir); v != Free {
		t.Fatalf("expected Free for npx, got %v", v)
	}
}

func TestNode_ScriptsAndSubcommands(t *testing.T) {
	tmpDir := t.TempDir()
	pkgContent := `{
  "name": "test-pkg",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "test": "vitest"
  }
}`
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgContent), 0644)

	cases := []struct {
		cmd      []string
		expected Verdict
	}{
		{[]string{"npm", "run", "dev"}, Valid},
		{[]string{"npm", "run-script", "build"}, Valid},
		{[]string{"npm", "test"}, Valid},
		{[]string{"npm", "run", "missing"}, Invalid},
		{[]string{"npm", "install"}, Free},
		{[]string{"npm", "add", "react"}, Free},
		{[]string{"pnpm", "dev"}, Valid},
		{[]string{"pnpm", "build"}, Valid},
		{[]string{"pnpm", "unknown_cmd"}, Unknown}, // bare word in pnpm/yarn/bun -> Unknown
		{[]string{"yarn", "dev"}, Valid},
		{[]string{"yarn", "build"}, Valid},
		{[]string{"yarn", "add", "lodash"}, Free},
		{[]string{"bun", "run", "dev"}, Valid},
		{[]string{"bun", "test"}, Valid},
		{[]string{"bun", "install"}, Free},
		{[]string{"bunx", "prisma", "generate"}, Free},
		{[]string{"pnpm", "dlx", "prisma"}, Free},
	}

	for _, tc := range cases {
		v := ValidateNode(tc.cmd, tmpDir)
		if v != tc.expected {
			t.Errorf("cmd %v: expected %v, got %v", tc.cmd, tc.expected, v)
		}
	}
}

func TestNode_NearestPackageJsonOnly(t *testing.T) {
	tmpDir := t.TempDir()
	// parent has script 'parent_script'
	parentPkg := `{ "scripts": { "parent_script": "echo parent" } }`
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(parentPkg), 0644)

	subDir := filepath.Join(tmpDir, "packages", "sub")
	_ = os.MkdirAll(subDir, 0755)
	// sub has script 'sub_script', but lacks 'parent_script'
	subPkg := `{ "scripts": { "sub_script": "echo sub" } }`
	_ = os.WriteFile(filepath.Join(subDir, "package.json"), []byte(subPkg), 0644)

	// from sub: sub_script is Valid, parent_script is Invalid
	if v := ValidateNode([]string{"npm", "run", "sub_script"}, subDir); v != Valid {
		t.Errorf("expected sub_script to be Valid, got %v", v)
	}
	if v := ValidateNode([]string{"npm", "run", "parent_script"}, subDir); v != Invalid {
		t.Errorf("expected parent_script to be Invalid from sub package, got %v", v)
	}
}

func TestNode_Flags(t *testing.T) {
	tmpDir := t.TempDir()
	pkgContent := `{ "scripts": { "dev": "vite" } }`
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgContent), 0644)

	// workspace flags -> Unknown
	if v := ValidateNode([]string{"npm", "run", "dev", "-w", "backend"}, tmpDir); v != Unknown {
		t.Errorf("expected -w flag to be Unknown, got %v", v)
	}
	if v := ValidateNode([]string{"pnpm", "--filter", "backend", "dev"}, tmpDir); v != Unknown {
		t.Errorf("expected --filter to be Unknown, got %v", v)
	}

	// args after -- ignored
	if v := ValidateNode([]string{"npm", "run", "dev", "--", "--port", "3000"}, tmpDir); v != Valid {
		t.Errorf("expected args after -- to be ignored, got %v", v)
	}
}
