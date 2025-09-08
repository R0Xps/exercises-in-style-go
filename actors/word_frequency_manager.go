package main

import "slices"

// WordFrequencyManager handles counting and sorting the words based on their frequencies
type WordFrequencyManager struct {
	messages chan []any
	freq     map[string]int
}

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

// Create and return a pointer to a new WordFrequencyManager object (actor)
func NewWordFrequencyManager() *WordFrequencyManager {
	wfm := &WordFrequencyManager{
		messages: make(chan []any, 100),
		freq:     make(map[string]int),
	}
	return wfm
}

// Send a message to this actor
func (wfm *WordFrequencyManager) Send(message []any) {
	wfm.messages <- message
}

// This function runs in a goroutine and keeps running until a 'die' message is received
func (wfm *WordFrequencyManager) Start() {
	for msg := range wfm.messages {
		wfm.dispatch(msg)
		if msg[0] == "die" {
			close(wfm.messages)
		}
	}
}

// Handle received messages, if they are a known type, run their appropriate functions, otherwise ignore them
func (wfm *WordFrequencyManager) dispatch(message []any) {
	if message[0] == "word" {
		wfm.increment(message[1:])
	} else if message[0] == "top25" {
		wfm.top25(message[1:])
	}
}

// Increments the frequency of a word in the freq map
func (wfm *WordFrequencyManager) increment(message []any) {
	word := message[0].(string)
	wfm.freq[word]++
}

// Returns a slice of all words and their frequencies ordered by frequency in descending order
func (wfm *WordFrequencyManager) top25(message []any) {
	recipient := message[0].(*WordFrequencyController)
	wordFreq := make([]wordFreqEntry, 0)
	for k, v := range wfm.freq {
		wordFreq = append(wordFreq, wordFreqEntry{k, v})
	}

	slices.SortFunc(wordFreq, func(i, j wordFreqEntry) int {
		return j.freq - i.freq
	})

	recipient.Send([]any{"top25", wordFreq})
}
