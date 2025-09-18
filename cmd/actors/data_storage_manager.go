package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DataStorageManager handles everything related to the input file
type DataStorageManager struct {
	messages        chan []any
	stopWordManager *StopWordManager
	data            string
}

// Create and return a pointer to a new DataStoragaManager object (actor)
func NewDataStorageManager() *DataStorageManager {
	dsm := &DataStorageManager{
		messages: make(chan []any, 100),
	}
	return dsm
}

// Send a message to this actor
func (dsm *DataStorageManager) Send(message []any) {
	dsm.messages <- message
}

// This function runs in a goroutine and keeps running until a 'die' message is received
func (dsm *DataStorageManager) Start() {
	for msg := range dsm.messages {
		dsm.dispatch(msg)
		if msg[0] == "die" {
			close(dsm.messages)
		}
	}
}

// Handle received messages, if they are a known type, run their appropriate functions, otherwise forward them to stopWordManager
func (dsm *DataStorageManager) dispatch(message []any) {
	switch message[0] {
	case "init":
		dsm.init(message[1:])
	case "send_word_freqs":
		dsm.processWords(message[1:])
	default:
		dsm.stopWordManager.Send(message)
	}
}

// Initialize the DataStorageManager object with a StopWordManager that is received in the message, and a string that is read and filtered from a file in a path received in the message as well
func (dsm *DataStorageManager) init(message []any) {
	inputFilePath := message[0].(string)
	dsm.stopWordManager = message[1].(*StopWordManager)

	file, err := os.Open(filepath.Clean(inputFilePath))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	dsm.data = filterAndNormalize(bytes)
}

// Split the data string into words, then forward them all to stopWordManager to filter, and send another message of type "top25" to a WordFrequencyManager through stopWordManager
func (dsm *DataStorageManager) processWords(message []any) {
	recipient := message[0].(*WordFrequencyController)
	words := strings.Fields(dsm.data)

	for _, w := range words {
		dsm.stopWordManager.Send([]any{"filter", w})
	}
	dsm.stopWordManager.Send([]any{"top25", recipient})
}

// Replace all non-letter characters with spaces, and convert all uppercase letters to lowercase, then return the result as a string
func filterAndNormalize(data []byte) string {
	for i, b := range data {
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if b < 'a' || b > 'z' {
			b = ' '
		}
		data[i] = b
	}
	return string(data)
}
