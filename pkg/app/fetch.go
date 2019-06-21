package app

import (
	"fmt"
	"net/http"
	"sync"

	// Support all image decoder/encoders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Fetch ranges over the input channel of URLs, attempting to fetch them over
// HTTP, and sends the fetched image out over the out channel. Errors are
// sent to the error channel. Calling Done on the given wait group allows
// all workers in the worker group to complete before closing channels.
func Fetch(in <-chan string, out chan<- Resource, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	for u := range in {
		// GET resource at given URL
		res, err := http.Get(u)
		if err != nil {
			errorCh <- fmt.Errorf("error fetching %s: %v", u, err)
			continue
		}
		// defer res.Body.Close()

		if res.StatusCode >= 400 {
			errorCh <- fmt.Errorf("invalid resource at %s: %v", u, res.StatusCode)
			continue
		}

		out <- Resource{res.Body, u}
	}
}
