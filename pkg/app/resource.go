package app

import "io"

// Resource represents a fetche resource's URL and Body
type Resource struct {
	Body io.ReadCloser
	URL  string
}
