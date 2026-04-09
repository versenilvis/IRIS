package root

import (
	"os"
	"sync"
)

var termMutex sync.Mutex

func TermWrite(data []byte) {
	termMutex.Lock()
	defer termMutex.Unlock()
	os.Stdout.Write(data)
}
