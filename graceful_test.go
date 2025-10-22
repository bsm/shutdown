package shutdown_test

import (
	"log"
	"net/http"

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
