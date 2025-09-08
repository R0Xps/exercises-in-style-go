package main

import (
	"fmt"
	"log"
)

// WordFrequencyController acts as the driver code for the term frequency task
type WordFrequencyController struct {
	messages           chan []any
	dataStorageManager *DataStorageManager
}

// Create and return a pointer to a new WordFrequencyController object (actor)
func NewWordFrequencyController() *WordFrequencyController {
	wfc := &WordFrequencyController{
		messages: make(chan []any, 100),
	}
	return wfc
}

// Send a message to this actor
func (wfc *WordFrequencyController) Send(message []any) {
	wfc.messages <- message
}

// This function runs in a goroutine and keeps running until either a 'die' message is received, or the final results are printed
func (wfc *WordFrequencyController) Start() {
	for msg := range wfc.messages {
		wfc.dispatch(msg)
		if msg[0] == "die" {
			close(wfc.messages)
		}
	}
}

// Handle received messages, if they are a known type, run their appropriate functions, otherwise print an error message
func (wfc *WordFrequencyController) dispatch(message []any) {
	if message[0] == "run" {
		wfc.run(message[1:])
	} else if message[0] == "top25" {
		wfc.display(message[1:])
	} else {
		log.Println("Unknown message type:", message[0])
	}
}

// Start the chain of messages leading to the execution of the term frequency task
func (wfc *WordFrequencyController) run(message []any) {
	wfc.dataStorageManager = message[0].(*DataStorageManager)
	wfc.dataStorageManager.Send([]any{"send_word_freqs", wfc})
}

// Print the top (25 max) words and their frequencies
func (wfc *WordFrequencyController) display(message []any) {
	wordFreq := message[0].([]wordFreqEntry)

	N := min(25, len(wordFreq))
	for _, wf := range wordFreq[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}

	wfc.dataStorageManager.Send([]any{"die"})
	close(wfc.messages)
}
