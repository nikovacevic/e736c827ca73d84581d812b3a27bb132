package app

// Write ...TODO
func Write(in <-chan string, out chan<- string) {
	defer close(out)
	for str := range in {
		// TODO write str
		out <- str
	}
}
