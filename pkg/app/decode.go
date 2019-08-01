package app

import (
	"image"
	"sync"

	// Support all image decoder/encoders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Decode accepts a Resource, decodes its body into an Image, passes that Image
// to the output channel, and closes the body. Errors are sent to the error
// channel. Calling Done on the given wait group allows all workers in the
// worker group to complete before closing channels.
func Decode(in <-chan Resource, out chan<- Image, errorCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for resource := range in {
		defer resource.Body.Close()

		// TODO protect against huge image files that could
		// exhaust memory resources. At least recover from panic.
		// https://golang.org/pkg/bytes/#Buffer.ReadFrom

		img, _, err := image.Decode(resource.Body)
		if err != nil {
			errorCh <- err
			continue
		}
		out <- Image{Image: &img, URL: resource.URL}
	}
}
