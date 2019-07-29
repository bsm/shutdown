package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/xerrors"
)

var defaultSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
}

// Context implements context.Context but optionally exposes wait for
type Context interface {
	context.Context

	// WaitFor waits for a blocking method to complete
	WaitFor(blocking func() error, signals ...os.Signal) error
}

// WithContext inits a new Context
func WithContext(parent context.Context) Context {
	ctx, cancel := context.WithCancel(parent)
	return &shutdown{Context: ctx, cancel: cancel}
}

type shutdown struct {
	context.Context
	cancel func()
}

// WaitFor accepts a blocking callback function and waits
// for the callback to return or a signal to trigger
func (s *shutdown) WaitFor(blocking func() error, signals ...os.Signal) error {
	defer s.cancel()

	errs := make(chan error, 1)
	go func() {
		errs <- blocking()
	}()

	sigs := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = defaultSignals
	}
	signal.Notify(sigs, signals...)

	select {
	case err := <-errs:
		return err
	case <-sigs:
	case <-s.Done():
	}

	if err := s.Err(); err != nil && !xerrors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
