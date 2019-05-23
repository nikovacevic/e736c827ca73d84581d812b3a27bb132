package app

import (
	"fmt"
	"strings"
	"sync"
)

// ReduceFn defines a reduction of an Image to a string
type ReduceFn func(Image) (string, error)

// Reduce ranges over the input channel of Images, reducing them to a string
// by a ReduceFn, then sends the output on the out channel. Errors are
// sent to the error channel. Calling Done on the given wait group allows
// all workers in the worker group to complete before closing channels.
func Reduce(in <-chan Image, out chan<- string, fn ReduceFn, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	for img := range in {
		str, err := fn(img)
		if err != nil {
			errorCh <- fmt.Errorf("error reducing image: %v", err)
			continue
		}
		out <- str
	}
}

// CountHexValues counts the frequency of simple (non-alpha-premultiplied) hex
// values in an image and returns a comma-separated row consisting of:
// "url,hex,hex,hex" for the top three hex values in the image.
func CountHexValues(img Image) (string, error) {
	// Count hex values
	Counter := NewCounter()

	bounds := (*img.Image).Bounds()
	// from https://golang.org/pkg/image/
	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// from https://golang.org/pkg/image/color/#Color
			// RGBA returns the alpha-premultiplied red, green, blue and alpha values
			// for the color. Each value ranges within [0, 0xffff], but is represented
			// by a uint32 so that multiplying by a blend factor up to 0xffff will not
			// overflow.
			ra, ga, ba, _ := (*img.Image).At(x, y).RGBA()
			// Bit shift to convert 16-bit values to simple 8-bit
			// r, g, b values between 0x00 and 0xff
			// TODO does this make for the "most correct", least-lossy value?
			hex := fmt.Sprintf("#%02x%02x%02x", ra>>8, ga>>8, ba>>8)
			Counter.Count(hex)
		}
	}

	// Formatted for CSV: url,hex,hex,hex
	str := img.URL + "," + strings.Join(Counter.Top(3), ",")

	return str, nil
}
