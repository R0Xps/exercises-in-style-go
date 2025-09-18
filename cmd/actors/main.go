package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	// Check for the required arguments
	if len(os.Args) != 3 {
		log.Fatal("required arguments: <stop_words_file> <input_file>")
	}

	// sync.WaitGroup is used to ensure all goroutines are done before exiting the program
	wg := new(sync.WaitGroup)

	// Create the needed actors, start their goroutines, and send their initialization messages
	wfm := NewWordFrequencyManager()
	wg.Go(wfm.Start)

	swm := NewStopWordManager()
	wg.Go(swm.Start)
	swm.Send([]any{"init", os.Args[1], wfm})

	dsm := NewDataStorageManager()
	wg.Go(dsm.Start)
	dsm.Send([]any{"init", os.Args[2], swm})

	wfc := NewWordFrequencyController()
	wg.Go(wfc.Start)
	wfc.Send([]any{"run", dsm})

	// This blocks until all goroutines are done
	wg.Wait()
}
