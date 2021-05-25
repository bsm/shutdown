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

// Wait accepts a callback function and waits for the callback to return or for
// a signal to trigger. If nil is passed instead of a function, Wait will block
// until a signal triggers.
func Wait(blocking func() error, signals ...os.Signal) error {
	return WaitContext(context.Background(), blocking, signals...)
}

// WaitContext behaves like Wait but with a parent context, which may include a
// deadline or a custom cancellation.
func WaitContext(parent context.Context, blocking func() error, signals ...os.Signal) error {
	if len(signals) == 0 {
		signals = defaultSignals
	}

	ctx, stop := signal.NotifyContext(parent, signals...)
	defer stop()

	if blocking == nil {
		blocking = func() error { <-ctx.Done(); return nil }
	}

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
