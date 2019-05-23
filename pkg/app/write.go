package app

import (
	"bufio"
	"fmt"
	"os"
)

// Fetch ranges over the input channel of URLs, attempting to fetch them over
// HTTP, and sends the fetched image out over the out channel. Errors are
// sent to the error channel. Calling Done on the given wait group allows
// all workers in the fetch worker group to complete before closing channels.

// Write ranges over the input channel of strings (CSV lines), writing them
// to the given file, and passes the string along to be logged for completion.
// Errors are sent to the error channel. Calling Done on the given wait group
// allows all workers in the worker group to complete before closing channels.
func Write(in <-chan string, out chan<- string, file *os.File, errorCh chan<- error) {
	defer close(out)
	w := bufio.NewWriter(file)
	for str := range in {
		_, err := w.WriteString(str + "\n")
		if err != nil {
			errorCh <- fmt.Errorf("error writing to file: %v", err)
			continue
		}
		w.Flush()
		out <- str
	}
}
