package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var defaultSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
}

// Wait accepts a blocking callback function and waits
// for the callback to return or a signal to trigger
func Wait(blocking func() error, signals ...os.Signal) error {
	return WaitContext(context.Background(), blocking, signals...)
}

// WaitContext behaves like Wait but with a custom parent context.
func WaitContext(parent context.Context, blocking func() error, signals ...os.Signal) error {
	if len(signals) == 0 {
		signals = defaultSignals
	}

	ctx, stop := signal.NotifyContext(parent, signals...)
	defer stop()

	errs := make(chan error, 1)
	go func() {
		errs <- blocking()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return nil
}
