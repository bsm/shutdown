package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var defaultSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
}

type Server interface {
	ListenAndServe() error
}

// WaitFor accepts a blocking callback function and waits
// for the callback to return or a signal to trigger
func WaitFor(blocking func() error, signals ...os.Signal) error {
	serverErrs := make(chan error, 1)
	go func() {
		serverErrs <- blocking()
	}()

	termSignals := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = defaultSignals
	}
	signal.Notify(termSignals, signals...)

	select {
	case err := <-serverErrs:
		return err
	case <-termSignals:
	}
	return nil
}

// Wait waits for a Server instance to shut down or a signal to trigger
func Wait(srv Server, signals ...os.Signal) error {
	return WaitFor(srv.ListenAndServe, signals...)
}
