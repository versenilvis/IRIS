package root

import (
	"os"
	"testing"

	"go.uber.org/goleak"
	"github.com/versenilvis/iris/internal/scoring"
)

func TestMain(m *testing.M) {
	code := m.Run()
	scoring.CloseGlobalFrecencyStore()
	if code == 0 {
		if err := goleak.Find(); err != nil {
			os.Stderr.WriteString("goleak: " + err.Error() + "\n")
			os.Exit(1)
		}
	}
	os.Exit(code)
}
