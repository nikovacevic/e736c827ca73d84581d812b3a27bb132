package app

import (
	"bufio"
	"log"
	"os"
)

// LogResults prints results to stdout and signals that the full set of jobs
// is complete when the channel closes.
func LogResults(in <-chan string, done chan<- bool) {
	for str := range in {
		log.Println(str)
	}
	done <- true
}

// LogErrors prints errors to stdout and writes them to an error.log file
func LogErrors(in <-chan error, file *os.File) {
	w := bufio.NewWriter(file)
	for e := range in {
		log.Println(e.Error())
		_, err := w.WriteString(e.Error() + "\n")
		if err != nil {
			log.Println(err)
			continue
		}
		w.Flush()
	}
}
