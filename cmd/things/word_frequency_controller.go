package main

import "fmt"

// WordFrequencyController holds objects of DataStorageManager, StopWordsManager, and WordFrequencyManager, and uses them together to complete the term frequency task and print its output
type WordFrequencyController struct {
	dataStorageManager   *DataStorageManager
	stopWordsManager     *StopWordsManager
	wordFrequencyManager *WordFrequencyManager
}

// Create and return a pointer to a new WordFrequencyController object, with objects of DataStorageManager, StopWordsManager, WordFrequencyManager all initialized with the appropriate values
func NewWordFrequencyController(stopWordsFilePath, inputFilePath string) *WordFrequencyController {
	return &WordFrequencyController{
		dataStorageManager:   NewDataStorageManager(inputFilePath),
		stopWordsManager:     NewStopWordsManager(stopWordsFilePath),
		wordFrequencyManager: NewWordFrequencyManager(),
	}
}

// Run the controller and use the 3 separate objects together to get the desired output and print it
func (wfc *WordFrequencyController) Run() {
	words := wfc.dataStorageManager.Words()
	for _, word := range words {
		if !wfc.stopWordsManager.IsStopWord(word) {
			wfc.wordFrequencyManager.Increment(word)
		}
	}

	wordFreq := wfc.wordFrequencyManager.Sorted()

	// Print the first 25 elements (or all elements if there are less than 25) of the wordFreq slice
	N := min(25, len(wordFreq))
	for _, wf := range wordFreq[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}
}
