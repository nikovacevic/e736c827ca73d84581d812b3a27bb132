package app

import (
	"bufio"
	"log"
	"os"
)

// Write ...TODO
func Write(in <-chan string, out chan<- string, file *os.File) {
	defer close(out)
	w := bufio.NewWriter(file)
	for str := range in {
		_, err := w.WriteString(str + "\n")
		if err != nil {
			log.Printf("Error writing to CSV: %v\n", err)
		}
		w.Flush()
		out <- str
	}
}
