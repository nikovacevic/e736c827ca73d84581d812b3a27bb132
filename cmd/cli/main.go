package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"sort"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	// TODO replace with txt source
	urls := []string{
		"http://i.imgur.com/FApqk3D.jpg",
		"http://i.imgur.com/TKLs9lo.jpg",
		"https://i.redd.it/d8021b5i2moy.jpg",
	}

	readCh := make(chan string)
	reduceCh := make(chan Image)
	writeCh := make(chan string)

	endCh := make(chan bool)

	// TODO make this a fan-out pattern across a pool of goroutines
	go Read(readCh, reduceCh)
	// TODO make this a fan-out pattern across a pool of goroutines
	go Reduce(reduceCh, writeCh)
	// TODO confirm this should be singular
	go Write(writeCh, endCh)

	for _, url := range urls {
		readCh <- url
	}
	close(readCh)

	<-endCh
}

// Read ...
func Read(in chan string, out chan Image) {
	defer close(out)
	for u := range in {
		// GET resource at given URL
		res, err := http.Get(u)
		if err != nil {
			// TODO replace this with better tracing
			fmt.Println(err)
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

// Image ...TODO
type Image struct {
	Image *image.Image
	URL   string
}

// HexCounter ...TODO
// Note: Do not track top values. Expected case, that would be worse runtime
// than O(N*T) runtime of on-demand solution.
type HexCounter struct {
	Values map[string]uint32
}

// NewHexCounter ...TODO
func NewHexCounter() *HexCounter {
	return &HexCounter{Values: map[string]uint32{}}
}

// Count ...TODO
func (hc *HexCounter) Count(hex string) {
	hc.Values[hex] = hc.Values[hex] + 1
}

// Slice ...TODO
func (hc *HexCounter) Slice() []HexPair {
	hps := []HexPair{}
	for hex, count := range hc.Values {
		hps = append(hps, HexPair{Hex: hex, Count: count})
	}
	return hps
}

// String ...TODO
func (hc *HexCounter) Top(n int) []string {
	pairs := hc.Slice()
	sort.Sort(sort.Reverse(ByCount(pairs)))
	strs := []string{}
	for _, hp := range pairs {
		strs = append(strs, hp.Hex)
	}
	return strs[0:n]
}

// HexPair ...TODO
type HexPair struct {
	Hex   string
	Count uint32
}

// ByCount implements sort.Interface for HexPairs based on HexPair.Count
type ByCount []HexPair

func (hps ByCount) Len() int           { return len(hps) }
func (hps ByCount) Swap(i, j int)      { hps[i], hps[j] = hps[j], hps[i] }
func (hps ByCount) Less(i, j int) bool { return hps[i].Count < hps[j].Count }

// Reduce ...TODO
func Reduce(in chan Image, out chan string) {
	defer close(out)
	for img := range in {
		// Count hex values
		hexCounter := NewHexCounter()

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
				hexCounter.Count(hex)
			}
		}

		out <- img.URL + "," + strings.Join(hexCounter.Top(3), ",")
	}
}

// Write ...TODO
func Write(in chan string, end chan bool) {
	for str := range in {
		fmt.Println(str)
	}
	end <- true
}
