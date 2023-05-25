package terminator

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const Timeout = 3 * time.Second

var (
	WaitGroup sync.WaitGroup
	Signal    = make(chan struct{})
)

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
