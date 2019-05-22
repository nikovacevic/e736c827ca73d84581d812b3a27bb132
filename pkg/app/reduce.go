package app

import (
	"fmt"
	"strings"
	"sync"
)

// Reduce ...TODO
func Reduce(in <-chan Image, out chan<- string, fn func(Image) (string, error), wg *sync.WaitGroup) {
	defer wg.Done()
	for img := range in {
		str, err := fn(img)
		if err != nil {
			// TODO error handling
		}
		out <- str
	}
}

// CountHexValues ...TODO
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
