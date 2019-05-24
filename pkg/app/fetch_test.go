package app

import (
	"sync"
	"testing"
)

func TestFetch(t *testing.T) {
	inCh := make(chan string)
	outCh := make(chan Image)
	errorCh := make(chan error)
	var wg sync.WaitGroup

	go Fetch(inCh, outCh, &wg, errorCh)

	var img Image
	var err error

	// Verify that valid images can be fetched and result in the proper output
	url := "http://i.imgur.com/FApqk3D.jpg"
	exp := "http://i.imgur.com/FApqk3D.jpg,#ffffff,#000000,#f3c300"
	inCh <- url
	select {
	case err = <-errorCh:
		t.Errorf("Error fetching %s: %v", url, err)
	case img = <-outCh:
		act, err := CountHexValues(img)
		if err != nil {
			t.Errorf(err.Error())
		}
		if act != exp {
			t.Errorf("Expected %s to yield %s, actually yielded %s", url, exp, act)
		}
	}

	// Verify that images that do not exist pipe errors properly
	url = "https://nikovacevic.io/img/123"
	exp = "invalid resource at https://nikovacevic.io/img/123: 404"
	inCh <- url
	select {
	case err = <-errorCh:
		if err.Error() != exp {
			t.Errorf("Expected error: %s, actually got: %s", exp, err.Error())
		}
	case img = <-outCh:
		t.Errorf("Expected error fetching %s, but got image", url)
	}
}
