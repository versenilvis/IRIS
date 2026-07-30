package root

import (
	"reflect"
	"testing"
)

func TestShellArgs(t *testing.T) {
	tests := []struct {
		name      string
		shellName string
		login     bool
		want      []string
	}{
		{name: "bash login", shellName: "bash", login: true, want: []string{"--login"}},
		{name: "bash non-login", shellName: "bash", login: false},
		{name: "zsh unchanged", shellName: "zsh", login: true},
		{name: "fish unchanged", shellName: "fish", login: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellArgs(tt.shellName, tt.login); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("shellArgs(%q, %v) = %v, want %v", tt.shellName, tt.login, got, tt.want)
			}
		})
	}
}

func TestLoginFlagIsRegistered(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("login")
	if flag == nil {
		t.Fatal("expected --login flag to be registered")
	}
	if flag.DefValue != "false" {
		t.Fatalf("expected --login to default to false, got %q", flag.DefValue)
	}
}
