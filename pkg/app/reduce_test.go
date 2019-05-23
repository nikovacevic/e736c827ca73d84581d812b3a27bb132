package app

import (
	"image"
	"log"
	"os"
	"testing"

	_ "image/png"
)

func loadImage(path string) Image {
	f, _ := os.Open(path)
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return Image{Image: &img, URL: path}
}

func TestCountHexValues(t *testing.T) {
	// red only has one color. It should return just that one color.
	red := loadImage("test/red.png")
	exp := "test/red.png,#ff0000"
	act, err := CountHexValues(red)
	if err != nil {
		t.Errorf(err.Error())
	}
	if act != exp {
		t.Errorf("Expected red.png to yield %s, actually yielded %s", exp, act)
	}

	// wbr should return white, blue, red.
	wrb := loadImage("test/wrb.png")
	exp = "test/wrb.png,#ffffff,#ff0000,#0000ff"
	act, err = CountHexValues(wrb)
	if err != nil {
		t.Errorf(err.Error())
	}
	if act != exp {
		t.Errorf("Expected wrb.png to yield %s, actually yielded %s", exp, act)
	}

	// wgr should return white, green, red. Red is tied with blue for the
	// final color, but it's hex value is higher.
	wgr := loadImage("test/wgr.png")
	exp = "test/wgr.png,#ffffff,#00ff00,#ff0000"
	act, err = CountHexValues(wgr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if act != exp {
		t.Errorf("Expected wgr.png to yield %s, actually yielded %s", exp, act)
	}

}
