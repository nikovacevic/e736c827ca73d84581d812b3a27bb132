package app

import "image"

// Image represents a decoded image and its URL
type Image struct {
	Image *image.Image
	URL   string
}
