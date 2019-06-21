package app

import (
	"bytes"
	"sync"
	"testing"
)

func TestFetch(t *testing.T) {
	inCh := make(chan string)
	outCh := make(chan Resource)
	errorCh := make(chan error)
	var wg sync.WaitGroup

	go Fetch(inCh, outCh, &wg, errorCh)

	var resource Resource
	var err error

	// Verify that valid images can be fetched and result in the proper output
	url := "http://i.imgur.com/FApqk3D.jpg"
	inCh <- url
	select {
	case err = <-errorCh:
		t.Errorf("Error fetching %s: %v", url, err)
	case resource = <-outCh:
		buf := new(bytes.Buffer)
		buf.ReadFrom(resource.Body)
		if len(buf.String()) == 0 {
			t.Errorf("Error, bytes should not be empty")
		}
	}

	// Verify that images that do not exist pipe errors properly
	url = "https://nikovacevic.io/img/123"
	exp := "invalid resource at https://nikovacevic.io/img/123: 404"
	inCh <- url
	select {
	case err = <-errorCh:
		if err.Error() != exp {
			t.Errorf("Expected error: %s, actually got: %s", exp, err.Error())
		}
	case resource = <-outCh:
		t.Errorf("Expected error fetching %s, but got image", url)
	}
}
