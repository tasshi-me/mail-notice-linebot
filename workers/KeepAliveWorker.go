package workers

import "time"

// KeepAliveWorker ..
func KeepAliveWorker(url string) {
	interval := 10 * time.Minute
	tic := time.NewTicker(interval)
	for {
		select {
		case <-tic.C:
			//TODO: Curl itself
		}
	}
}
