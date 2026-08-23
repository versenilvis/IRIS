package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDirPrecedence(t *testing.T) {
	custom := t.TempDir()
	xdg := t.TempDir()
	home := t.TempDir()

	t.Setenv("HOME", home)

	t.Run("override wins over XDG", func(t *testing.T) {
		t.Setenv(ConfigDirEnv, custom)
		t.Setenv("XDG_CONFIG_HOME", xdg)

		got, err := ConfigPath()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(custom, "config.toml"); got != want {
			t.Errorf("ConfigPath() = %q, want %q", got, want)
		}
		theme, err := ThemePath()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(custom, "theme.toml"); theme != want {
			t.Errorf("ThemePath() = %q, want %q", theme, want)
		}
	})

	t.Run("XDG wins over home", func(t *testing.T) {
		t.Setenv(ConfigDirEnv, "")
		t.Setenv("XDG_CONFIG_HOME", xdg)

		got, err := ConfigPath()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(xdg, "iris", "config.toml"); got != want {
			t.Errorf("ConfigPath() = %q, want %q", got, want)
		}
	})

	t.Run("home is the fallback", func(t *testing.T) {
		t.Setenv(ConfigDirEnv, "")
		t.Setenv("XDG_CONFIG_HOME", "")

		got, err := ConfigPath()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(home, ".config", "iris", "config.toml"); got != want {
			t.Errorf("ConfigPath() = %q, want %q", got, want)
		}
	})
}

// an override that points nowhere is a typo or a broken generated path;
// silently using the default location hides it
func TestConfigDirRejectsUnusableOverride(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "not-a-dir")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		dir    string
		reason string
	}{
		{"missing", filepath.Join(dir, "does-not-exist"), "does not exist"},
		{"not a directory", file, "is not a directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(ConfigDirEnv, tt.dir)

			_, err := ConfigPath()
			var badDir *ErrConfigDir
			if !errors.As(err, &badDir) {
				t.Fatalf("ConfigPath() error = %v, want *ErrConfigDir", err)
			}
			if badDir.Reason != tt.reason {
				t.Errorf("reason = %q, want %q", badDir.Reason, tt.reason)
			}
			if _, err := ThemePath(); !errors.As(err, &badDir) {
				t.Errorf("ThemePath() error = %v, want *ErrConfigDir", err)
			}
		})
	}
}

func TestConfigDirAcceptsRelativeOverride(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "iris-conf")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	t.Setenv(ConfigDirEnv, "iris-conf")

	got, err := ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("ConfigPath() = %q, want an absolute path", got)
	}
	if filepath.Base(filepath.Dir(got)) != "iris-conf" {
		t.Errorf("ConfigPath() = %q, want it under iris-conf", got)
	}
}

// an existing override directory with no files in it is normal: config init
// has to be able to create them, and an absent theme means the built-in one
func TestConfigDirAllowsEmptyOverrideDirectory(t *testing.T) {
	t.Setenv(ConfigDirEnv, t.TempDir())

	if _, err := ConfigPath(); err != nil {
		t.Errorf("ConfigPath() on an empty override dir = %v, want nil", err)
	}
	if _, err := ThemePath(); err != nil {
		t.Errorf("ThemePath() on an empty override dir = %v, want nil", err)
	}
}
