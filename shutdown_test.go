package shutdown_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/bsm/shutdown"
)

func ExampleWait() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.FileServer(http.Dir("/usr/share/doc")),
	}

	// Wait for either SIGINT/SIGTERM or ListenAndServe to exit.
	// Handle errors.
	err := shutdown.Wait(srv.ListenAndServe)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Server error", err)
	}

	// Perform a graceful server shutdown.
	log.Println("Shutting down ...")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Println("Shutdown error", err)
	}
}

func TestWait_fails_immediately(t *testing.T) {
	sentinel := errors.New("doh!")
	err := shutdown.Wait(func() error { return sentinel })
	if err == nil {
		t.Fatalf("expected error, got nil")
	} else if !errors.Is(err, sentinel) {
		t.Fatalf("expected specific error, got %v", err)
	}
}

func TestWaitContext_nil_callback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := shutdown.WaitContext(ctx, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := ctx.Err(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected specific error, got %v", err)
	}
}
