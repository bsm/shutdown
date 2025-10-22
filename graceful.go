package shutdown

import (
	"context"
	"errors"
	"time"
)

// DefaultShutdownTimeout defines how long to wait before forcibly shutting down.
var DefaultShutdownTimeout = 5 * time.Second

// Inspired by https://github.com/ory/graceful, just not net/http specific.
// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// StartFunc is the type of the function invoked by Graceful to start the server
type StartFunc func() error

// ShutdownFunc is the type of the function invoked by Graceful to shutdown the server
type ShutdownFunc func(context.Context) error

// Graceful is a short-hand for Wait with error handling and automatic shutdown.
func Graceful(start StartFunc, shutdown ShutdownFunc, expErrs ...error) error {
	return GracefulContext(context.Background(), start, shutdown, expErrs...)
}

// GracefulContext is a short-hand for WaitContext with error handling and automatic shutdown.
func GracefulContext(ctx context.Context, start StartFunc, shutdown ShutdownFunc, expErrs ...error) error {
	err := WaitContext(ctx, start)
	if err != nil {
		for _, exp := range expErrs {
			if errors.Is(err, exp) {
				return nil
			}
		}
		return err
	}

	timeout, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if err := shutdown(timeout); err != nil {
		return err
	}
	return nil
}
