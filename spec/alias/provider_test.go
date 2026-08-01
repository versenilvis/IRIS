package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoProvider(t *testing.T) {
	tempDir := t.TempDir()
	cargoDir := filepath.Join(tempDir, ".cargo")
	if err := os.Mkdir(cargoDir, 0755); err != nil {
		t.Fatal(err)
	}

	configToml := `
[alias]
b = "build"
c = "check"
t = ["test", "--", "--nocapture"]
`
	if err := os.WriteFile(filepath.Join(cargoDir, "config.toml"), []byte(configToml), 0644); err != nil {
		t.Fatal(err)
	}

	p := &CargoProvider{}
	aliases := p.parse(tempDir)

	expected := map[string]string{
		"b": "build",
		"c": "check",
		"t": "test -- --nocapture",
	}

	if len(aliases) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(aliases))
	}

	for _, a := range aliases {
		if expected[a.Name] != a.Expansion {
			t.Errorf("expected %s to map to %s, got %s", a.Name, expected[a.Name], a.Expansion)
		}
	}
}
