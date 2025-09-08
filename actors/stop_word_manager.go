package main

import (
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

// StopWordsManager handles everything about stop words, starting at reading them from a file, up to filtering words and only forwarding non-stop words
type StopWordManager struct {
	messages             chan []any
	wordFrequencyManager *WordFrequencyManager
	stopWords            []string
}

// Create and return a pointer to a new StopWordManager object (actor)
func NewStopWordManager() *StopWordManager {
	swm := &StopWordManager{
		messages: make(chan []any, 100),
	}
	return swm
}

// Send a message to this actor
func (swm *StopWordManager) Send(message []any) {
	swm.messages <- message
}

// This function runs in a goroutine and keeps running until a 'die' message is received
func (swm *StopWordManager) Start() {
	for msg := range swm.messages {
		swm.dispatch(msg)
		if msg[0] == "die" {
			close(swm.messages)
		}
	}
}

// Handle received messages, if they are a known type, run their appropriate functions, otherwise forward them to wordFrequencyManager
func (swm *StopWordManager) dispatch(message []any) {
	if message[0] == "init" {
		swm.init(message[1:])
	} else if message[0] == "filter" {
		swm.filter(message[1:])
	} else {
		swm.wordFrequencyManager.Send(message)
	}
}

// Initializes the StopWordManager object with a WordFrequencyManager that is received in the message, and a slice of stop words that are read from a file in a path received in the message as well
func (swm *StopWordManager) init(message []any) {
	stopWordsFilePath := message[0].(string)
	swm.wordFrequencyManager = message[1].(*WordFrequencyManager)

	file, err := os.Open(stopWordsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	swm.stopWords = strings.Split(string(bytes), ",")
}

// Filter received words and only forward non-stop words to wordFrequencyManager
func (swm *StopWordManager) filter(message []any) {
	word := message[0].(string)
	if !slices.Contains(swm.stopWords, word) {
		swm.wordFrequencyManager.Send([]any{"word", word})
	}
}
