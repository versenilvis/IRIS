package root

import (
	"reflect"
	"testing"
)

func TestShellArgs(t *testing.T) {
	tests := []struct {
		name  string
		login bool
		want  []string
	}{
		{name: "login", login: true, want: []string{"--login"}},
		{name: "non-login", login: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellArgs(tt.login); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("shellArgs(%v) = %v, want %v", tt.login, got, tt.want)
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
