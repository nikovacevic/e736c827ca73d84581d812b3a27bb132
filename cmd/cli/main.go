package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nikovacevic/e736c827ca73d84581d812b3a27bb132/pkg/app"
)

// FetchWorkers determines the size of the worker pool fetching URLs
const FetchWorkers = 25

// ReduceWorkers determines the size of the worker pool reducing Images
const ReduceWorkers = 10

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a file path")
	}

	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer file.Close()

	// Set up concurrent pipeline so that down-stream processes are not
	// blocked by expensive up-stream processes (e.g. fetching a file should
	// not block reducing and writing the previous file) and so that a
	// system with more resources can act in parallel.
	//
	// 1. Read filename (1 worker)
	//    | | fetchCh
	// 2. Fetch filename (many workers)
	//    | | reduceCh
	// 3. Reduce image (many workers)
	//    | | writeCh
	// 4. Write to CSV (1 worker)
	//    | | resultCh
	// 5. Log results (1 worker)
	//    | | doneCh (signals completion)

	fetchCh := make(chan string)
	reduceCh := make(chan app.Image)
	writeCh := make(chan string)
	resultCh := make(chan string)
	doneCh := make(chan bool)

	// Set up fetch worker pool
	var readWG sync.WaitGroup
	for w := 0; w < FetchWorkers; w++ {
		readWG.Add(1)
		go app.Fetch(fetchCh, reduceCh, &readWG)
	}

	// Set up reduce worker pool
	var reduceWG sync.WaitGroup
	for w := 0; w < ReduceWorkers; w++ {
		reduceWG.Add(1)
		// Count frequency of hex values as the reduce action
		go app.Reduce(reduceCh, writeCh, app.CountHexValues, &reduceWG)
	}

	// Set up write worker
	go app.Write(writeCh, resultCh)

	// Set up result logger
	go app.LogResults(resultCh, doneCh)

	start := time.Now()

	// Read all URLs from file, then tear-down close fetch input
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fetchCh <- scanner.Text()
	}
	close(fetchCh)

	// Close reduce input after read workers finish
	readWG.Wait()
	close(reduceCh)

	// Close write input after reduce workers finish
	reduceWG.Wait()
	close(writeCh)

	<-doneCh
	end := time.Now().Sub(start)
	log.Printf("Finished in %v\n", end)
}
