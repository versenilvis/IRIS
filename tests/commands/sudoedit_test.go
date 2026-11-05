package tests

import (
	"testing"

	_ "github.com/versenilvis/iris/commands/sys"
	"github.com/versenilvis/iris/spec"
)

func TestSudoeditSpec(t *testing.T) {
	s := spec.Registry["sudoedit"]
	if s == nil {
		t.Fatal("expected sudoedit spec to be registered")
	}

	foundAskpass := false
	for _, opt := range s.Options {
		if opt.Name == "--askpass" {
			foundAskpass = true
			break
		}
	}
	if !foundAskpass {
		t.Error("expected sudoedit to have --askpass option")
	}
}
