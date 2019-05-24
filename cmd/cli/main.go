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
	if len(os.Args) < 3 {
		log.Fatal("Please provide two file paths: one for reading and one for writing")
	}

	// Open read-file
	inPath := os.Args[1]
	inFile, err := os.Open(inPath)
	if err != nil {
		log.Fatal("Failed to open read file")
	}
	defer inFile.Close()

	// Open write-file
	outPath := os.Args[2]
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatal("Failed to open write file")
	}
	defer outFile.Close()

	// Open error-file
	errFile, err := os.Create("resources/error.log")
	if err != nil {
		log.Fatal("Failed to open error file")
	}
	defer errFile.Close()

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

	// data channels
	fetchCh := make(chan string)
	reduceCh := make(chan app.Image)
	writeCh := make(chan string)
	// logging channels
	resultCh := make(chan string, 100)
	errorCh := make(chan error, 100)
	// channel to flag completion
	doneCh := make(chan bool)

	// Set up fetch worker pool
	var fetchWG sync.WaitGroup
	for w := 0; w < FetchWorkers; w++ {
		fetchWG.Add(1)
		go app.Fetch(fetchCh, reduceCh, &fetchWG, errorCh)
	}

	// Set up reduce worker pool
	var reduceWG sync.WaitGroup
	for w := 0; w < ReduceWorkers; w++ {
		reduceWG.Add(1)
		// Count frequency of hex values as the reduce action
		go app.Reduce(reduceCh, writeCh, app.CountHexValues, &reduceWG, errorCh)
	}

	// Set up write worker
	go app.Write(writeCh, resultCh, outFile, errorCh)

	// Set up logging
	go app.LogResults(resultCh, doneCh)
	go app.LogErrors(errorCh, errFile)

	start := time.Now()

	// Read all URLs from file, then tear-down close fetch input
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		fetchCh <- scanner.Text()
	}
	close(fetchCh)

	// Close reduce input after read workers finish
	fetchWG.Wait()
	close(reduceCh)

	// Close write input after reduce workers finish
	reduceWG.Wait()
	close(writeCh)

	<-doneCh
	end := time.Now().Sub(start)
	log.Printf("Finished in %v\n", end)
}
