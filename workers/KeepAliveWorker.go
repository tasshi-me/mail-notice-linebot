package workers

import (
	"log"
	"net/http"
	"time"
)

// KeepAliveWorker ..
func KeepAliveWorker(url string) {
	interval := 20 * time.Minute
	tic := time.NewTicker(interval)
	for {
		select {
		case <-tic.C:
			res, _ := http.Get(url)
			defer res.Body.Close()
			log.Println("Keep-Alive: ", url)
		}
	}
}
