package root

import "testing"

func TestScanConfigDirFlag(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		want  string
		found bool
	}{
		{"separate value", []string{"--config-dir", "/etc/iris"}, "/etc/iris", true},
		{"equals form", []string{"--config-dir=/etc/iris"}, "/etc/iris", true},
		{"after other flags", []string{"-d", "--config-dir", "/etc/iris"}, "/etc/iris", true},
		{"before a subcommand", []string{"--config-dir=/etc/iris", "config", "show"}, "/etc/iris", true},
		{"absent", []string{"config", "show"}, "", false},
		{"dangling with no value", []string{"--config-dir"}, "", false},
		{"empty value is still explicit", []string{"--config-dir="}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := scanConfigDirFlag(tt.args)
			if found != tt.found {
				t.Fatalf("scanConfigDirFlag(%q) found = %v, want %v", tt.args, found, tt.found)
			}
			if got != tt.want {
				t.Errorf("scanConfigDirFlag(%q) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
