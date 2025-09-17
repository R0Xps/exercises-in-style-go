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
	// Call functions in order. Each function is explained below
	printTop25(sort(frequencies(removeStopWords(os.Args[1])(split(filterAndNormalize(readInputFile(os.Args[2])))))))
}

// Read the input file from the given path and return its contents as a slice of bytes
func readInputFile(filePath string) []byte {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return fileBytes
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

// Split the string around spaces and return a slice of strings containing all the words in the given string
func split(str string) []string {
	return strings.Fields(str)
}

// Here I used currying to convert a function that takes multiple arguments removeStopWords(stopWordsPath string, allWords []string) []string to a sequence of 2 functions that take 1 argument each
func removeStopWords(stopWordsPath string) func([]string) []string {
	// Return a new slice of strings containing only words that should be counted (non-stop words)
	return func(allWords []string) []string {
		stopWordsBytes := readInputFile(stopWordsPath)
		stopWordsStr := filterAndNormalize(stopWordsBytes)
		stopWords := strings.Fields(stopWordsStr)

		words := make([]string, 0)
		for _, w := range allWords {
			if !slices.Contains(stopWords, w) {
				words = append(words, w)
			}
		}

		return words
	}
}

// Return a map where the words are the keys and the values are their frequencies in the given words slice
func frequencies(words []string) map[string]int {
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}
	return freq
}

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

// Return a slice of wordFreqEntry containing all entries from the given map, sorted by frequency in descending order
func sort(freq map[string]int) []wordFreqEntry {
	wordFreq := make([]wordFreqEntry, 0, len(freq))
	for k, v := range freq {
		wordFreq = append(wordFreq, wordFreqEntry{k, v})
	}

	slices.SortFunc(wordFreq, func(i, j wordFreqEntry) int {
		return j.freq - i.freq
	})

	return wordFreq
}

// Print the first 25 elements (or all elements if there are less than 25) of the given list
func printTop25(wordFreq []wordFreqEntry) {
	N := min(25, len(wordFreq))
	for _, wf := range wordFreq[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}
}
