package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

// DataStorageManager stores the contents of the input file, and can return a slice of all words in that file on demand
type DataStorageManager struct {
	data string
}

// Create and return a pointer to a new DataStorageManager object with its data being the filtered and normalized version of the contents of the file at inputFilePath
func NewDataStorageManager(inputFilePath string) *DataStorageManager {
	file, err := os.Open(filepath.Clean(inputFilePath))
	if err != nil {
		log.Fatal(err)
	}

	rawData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	data := filterAndNormalize(rawData)

	return &DataStorageManager{
		data: data,
	}
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

// Return a slice containing the words of the data string in the DataStorageManager object
func (dsm *DataStorageManager) Words() []string {
	return strings.Fields(dsm.data)
}

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

// WordFrequencyManager keeps track of the frequency of words, and returns a sorted slice of words and their frequencies on demand
type WordFrequencyManager struct {
	freq map[string]int
}

// Create and return a pointer to a new WordFrequencyManager object, with an empty frequency map
func NewWordFrequencyManager() *WordFrequencyManager {
	freq := make(map[string]int)
	return &WordFrequencyManager{
		freq: freq,
	}
}

// Increment the frequency of the given word
func (wfm *WordFrequencyManager) Increment(word string) {
	wfm.freq[word]++
}

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

// Return a list of words and their frequencies sorted by frequency in descending order
func (wfm *WordFrequencyManager) Sorted() []wordFreqEntry {
	wordFreq := make([]wordFreqEntry, 0, len(wfm.freq))
	for k, v := range wfm.freq {
		wordFreq = append(wordFreq, wordFreqEntry{word: k, freq: v})
	}

	slices.SortFunc(wordFreq, func(i, j wordFreqEntry) int {
		return j.freq - i.freq
	})

	return wordFreq
}

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
