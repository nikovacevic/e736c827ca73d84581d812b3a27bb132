package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	red := color.RGBA{0xff, 0x00, 0x00, 0xff}
	green := color.RGBA{0x00, 0xff, 0x00, 0xff}
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}

	width := 10
	height := 10

	start := image.Point{0, 0}
	end := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{start, end})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x%3 == 0 {
				img.Set(x, y, white)
			} else if (x+y)%2 == 0 {
				img.Set(x, y, green)
			} else if x%2 == 0 {
				img.Set(x, y, red)
			} else {
				img.Set(x, y, blue)
			}
		}
	}

	f, _ := os.Create("test/image.png")
	png.Encode(f, img)
}
