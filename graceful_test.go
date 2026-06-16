package shutdown_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/bsm/shutdown"
)

func ExampleGraceful() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.FileServer(http.Dir("/usr/share/doc")),
	}

	if err := shutdown.Graceful(srv.ListenAndServe, srv.Shutdown, http.ErrServerClosed); err != nil {
		log.Fatalln("Server error", err)
	}
}

func TestGracefulContext_start_error(t *testing.T) {
	boom := errors.New("boom")

	// Unexpected start errors are surfaced and shutdown is not invoked.
	called := false
	err := shutdown.GracefulContext(context.Background(),
		func() error { return boom },
		func(context.Context) error { called = true; return nil },
	)
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
	if called {
		t.Fatal("shutdown should not run when start fails unexpectedly")
	}

	// Expected start errors are swallowed.
	err = shutdown.GracefulContext(context.Background(),
		func() error { return boom },
		func(context.Context) error { return nil },
		boom,
	)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

type ctxKey string

func TestGracefulContext_shutdown_detached(t *testing.T) {
	// A cancelled parent triggers shutdown; the shutdown context must remain
	// live (fresh timeout) yet still carry the parent's values.
	parent, cancel := context.WithCancel(context.WithValue(context.Background(), ctxKey("k"), "v"))
	cancel()

	err := shutdown.GracefulContext(parent,
		func() error { <-parent.Done(); return nil },
		func(ctx context.Context) error {
			if ctx.Err() != nil {
				t.Errorf("shutdown context already cancelled: %v", ctx.Err())
			}
			if got := ctx.Value(ctxKey("k")); got != "v" {
				t.Errorf("parent value not preserved, got %v", got)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
