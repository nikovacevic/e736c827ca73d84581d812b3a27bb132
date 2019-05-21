package main

func main() {
	// TODO replace with txt source
	urls := []string{
		"http://i.imgur.com/FApqk3D.jpg",
		"http://i.imgur.com/TKLs9lo.jpg",
		"https://i.redd.it/d8021b5i2moy.jpg",
	}

	readCh := make(chan string)
	reduceCh := make(chan [][]int)
	writeCh := make(chan string)

	// TODO make this a fan-out pattern across a pool of goroutines
	go Read(urls, readCh, reduceCh)
	// TODO make this a fan-out pattern across a pool of goroutines
	go Reduce(reduceCh, writeCh)
	// TODO confirm this should be singular
	go Write(writeCh)
}

// Read ...
func Read(urls []string, in chan string, out chan [][]int) {

}

// Reduce ...
func Reduce(in chan [][]int, out chan string) {

}

// Write ...
func Write(in chan string) {

}
