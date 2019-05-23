package app

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"sync"
)

// Fetch ranges over the input channel of URLs, attempting to fetch them over
// HTTP, and sends the fetched image out over the out channel. Errors are
// sent to the error channel. Calling Done on the given wait group allows
// all workers in the worker group to complete before closing channels.
func Fetch(in <-chan string, out chan<- Image, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	for u := range in {
		// GET resource at given URL
		res, err := http.Get(u)
		if err != nil {
			errorCh <- fmt.Errorf("error fetching %s: %v", u, err)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			errorCh <- fmt.Errorf("invalid resource at %s: %v", u, res.StatusCode)
			continue
		}

		// TODO protect against huge image files that could
		// exhaust memory resources. At least recover from panic.
		// https://golang.org/pkg/bytes/#Buffer.ReadFrom

		img, _, err := image.Decode(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		out <- Image{Image: &img, URL: u}
	}
}
