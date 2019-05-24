package main

import (
	"bufio"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	// Test files
	inPath := "test/input.txt"
	outPath := "test/output.txt"
	errPath := "test/error.log"
	// Run pipeline
	run(inPath, outPath, errPath)

	// Assert 3 successes, 1 failure
	expSuccess := map[string]bool{
		"http://i.imgur.com/TKLs9lo.jpg,#ffffff,#fefefe,#f7f7f7\n":     true,
		"http://i.imgur.com/FApqk3D.jpg,#ffffff,#000000,#f3c300\n":     true,
		"https://i.redd.it/d8021b5i2moy.jpg,#ffffff,#010304,#020405\n": true,
	}
	expFailure := map[string]bool{
		"invalid resource at https://nikovacevic.io/img/123: 404\n": true,
	}

	// Open write-file
	outFile, err := os.Open(outPath)
	if err != nil {
		t.Errorf("Failed to open write file")
	}
	defer outFile.Close()

	// Assert correctness of output file
	outBuf := bufio.NewReader(outFile)
	for o := 0; o < 3; o++ {
		line, err := outBuf.ReadString('\n')
		if err != nil {
			t.Errorf("Expected three lines in %s, got fewer: %v", outPath, err)
		}
		if _, ok := expSuccess[line]; !ok {
			t.Errorf("Unexpected line in output: %s", line)
		}
	}

	// Open error-file
	errFile, err := os.Open(errPath)
	if err != nil {
		t.Errorf("Failed to open error file")
	}
	defer errFile.Close()

	// Assert correctness of error file contents
	errBuf := bufio.NewReader(errFile)
	line, err := errBuf.ReadString('\n')
	if err != nil {
		t.Errorf("Expected one line in %s, got fewer: %v", errPath, err)
	}
	if _, ok := expFailure[line]; !ok {
		t.Errorf("Unexpected line in error log: %s", line)
	}
}
