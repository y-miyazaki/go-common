package signal

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDetectSignal(t *testing.T) {
	var receivedSig os.Signal
	var wg sync.WaitGroup
	wg.Add(1)

	DetectSignal(func(sig os.Signal) {
		receivedSig = sig
		wg.Done()
	}, os.Interrupt)

	// Simulate sending signal
	go func() {
		time.Sleep(10 * time.Millisecond)
		pid := os.Getpid()
		process, err := os.FindProcess(pid)
		if err == nil {
			process.Signal(os.Interrupt)
		}
	}()

	wg.Wait()
	assert.Equal(t, os.Interrupt, receivedSig)
}
