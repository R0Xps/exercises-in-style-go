package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// StopWordsManager handles everything to do with stop words
type StopWordsManager struct {
	stopWords []string
}

// Create and return a pointer to a new StopWordsManager with its stopWords field initialized to the words in the file at stopWordsFilePath
func NewStopWordsManager(stopWordsFilePath string) *StopWordsManager {
	file, err := os.Open(filepath.Clean(stopWordsFilePath))
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	stopWordsStr := filterAndNormalize(data)

	return &StopWordsManager{
		stopWords: strings.Fields(stopWordsStr),
	}
}

// Check if the given word is in the stopWords slice
func (swm *StopWordsManager) IsStopWord(word string) bool {
	return slices.Contains(swm.stopWords, word)
}
