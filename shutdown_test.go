package shutdown_test

import (
	"log"
	"net/http"
	"time"

	"github.com/bsm/shutdown"
)

func Example() {
	h := http.FileServer(http.Dir("/usr/share/doc"))
	s := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	err := shutdown.Wait(s)
	if err != nil {
		log.Fatal("FATAL ", err.Error())
	}
	s.SetKeepAlivesEnabled(false)
	time.Sleep(time.Second)
}
