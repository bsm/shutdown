package shutdown_test

import (
	"context"
	"log"
	"net/http"

	"github.com/bsm/shutdown"
)

func Example() {
	ctx := shutdown.WithContext(context.Background())
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.FileServer(http.Dir("/usr/share/doc")),
	}
	defer srv.Shutdown(context.Background())

	err := ctx.WaitFor(srv.ListenAndServe)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("Server error", err)
	}
	log.Println("Shutting down ...")
}
