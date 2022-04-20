package signal

import (
	"os"
	"os/signal"
)

// DetectSignal detects signals and sets func to process for a specific signal.
// syscall.SIGTERM or syscall.SIGKILL or os.Interrupt to describe what to do when the target container or OS stops.
func DetectSignal(f func(sig os.Signal), sig ...os.Signal) {
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, sig...)
		defer signal.Stop(s)
		sig := <-s
		f(sig)
	}()
}
