package terminator

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Timeout defines the amount of time to wait before exiting the application after a termination request
const Timeout = 3 * time.Second

var (
	// WaitGroup is the global semaphore used to indicate that a goroutine is running
	WaitGroup sync.WaitGroup

	// Signal is the global channel used to indicate an application termination request.
	// The channel is closed on the termination event.
	Signal = make(chan struct{})
)

// init calls a concurrent routine to respond to termination requests.
// In case of an event, the signal channel is closed and after the specified time,
// the application is terminated with exit code 143.
func init() {
	go func() {

		// Listen to interrupt and termination signals
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		<-signalCh
		close(Signal)

		// Guarantee termination after specified timeout
		<-time.After(Timeout)
		os.Exit(143)
	}()
}
