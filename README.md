# Shutdown

[![GoDoc](https://godoc.org/github.com/bsm/shutdown?status.png)](http://godoc.org/github.com/bsm/shutdown)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Wait for servers to terminate gracefully.

## Example:

```go
import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/bsm/shutdown"
)

func main() {
	srv := &http.Server{
		Addr:		":8080",
		Handler:	http.FileServer(http.Dir("/usr/share/doc")),
	}

	// Wait for either SIGINT/SIGTERM or ListenAndServe to exit.
	// Handle errors.
	err := shutdown.Wait(srv.ListenAndServe)
	if err != nil && err != http.ErrServerClosed {
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
```
