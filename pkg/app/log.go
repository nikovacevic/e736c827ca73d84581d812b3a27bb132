package app

import "log"

// LogResults ...TODO
func LogResults(in <-chan string, done chan<- bool) {
	for str := range in {
		// TODO better logging?
		log.Println(str)
	}
	done <- true
}
