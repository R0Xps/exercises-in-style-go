package main

import (
	"log"
	"os"
)

func main() {
	// Check for the required arguments
	if len(os.Args) != 3 {
		log.Fatal("required arguments: <stop_words_file> <input_file>")
	}

	// Initialize an instance of WordFrequencyController with the arguments passed to the program
	wfc := NewWordFrequencyController(os.Args[1], os.Args[2])
	wfc.Run()
}
