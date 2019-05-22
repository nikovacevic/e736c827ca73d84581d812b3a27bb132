package app

import (
	"image"
	"log"
	"net/http"
	"sync"
)

// Fetch ...TODO
func Fetch(in <-chan string, out chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	for u := range in {
		// GET resource at given URL
		res, err := http.Get(u)
		if err != nil {
			// TODO error handling
			// TODO replace this with better tracing
			continue
		}
		defer res.Body.Close()
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
