package signal

import (
	"os"
	"os/signal"
)

// DetectSignal detects signals and sets func to process for a specific signal.
// "syscall.SIGTERM" and "syscall.SIGKILL" to describe what to do when the target container or OS stops.
func DetectSignal(f func(), sig ...os.Signal) {
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, sig...)
		<-s
		f()
	}()
}
