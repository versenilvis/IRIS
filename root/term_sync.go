package root

import (
	"os"
	"sync"
)

var termMutex sync.Mutex

func TermWrite(data []byte) {
	termMutex.Lock()
	defer termMutex.Unlock()
	_, _ = os.Stdout.Write(data)
}
