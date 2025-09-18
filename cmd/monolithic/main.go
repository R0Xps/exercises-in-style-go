package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

func main() {
	// Get arguments for stopWordsPath and inputPath
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("required arguments: <stop_words_file> <input_file>")
	}
	stopWordsPath := args[0]
	inputPath := args[1]

	// Open the file located at stopWordsPath and read, then convert it to a slice of words
	stopWordsFile, err := os.Open(filepath.Clean(stopWordsPath))
	if err != nil {
		log.Fatal(err)
	}
	defer func(stopWordsFile *os.File) {
		err := stopWordsFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stopWordsFile)

	stopWordsBytes, err := io.ReadAll(stopWordsFile)
	if err != nil {
		log.Fatal(err)
	}

	stopWordsString := string(stopWordsBytes) + " "

	stopWords := make([]string, 0)

	start := -1
	// Iterate over characters in the stop words file
	for i, c := range stopWordsString {
		if start == -1 {
			// We're currently not in a word
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				// This means we found the start of a word
				start = i
			}
		} else {
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				// We're still inside a word
				continue
			}
			// When we reach this point, we're at the character immediately after a word

			// Copy the entire word and convert it to lowercase
			word := strings.ToLower(stopWordsString[start:i])
			// Add the word to the stopWords slice
			stopWords = append(stopWords, word)
			// After we're done with a word, we want to look for the next one, so we reset the start index to -1
			start = -1
		}
	}

	// Open and read the file located at inputPath
	inputFile, err := os.Open(filepath.Clean(inputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(inputFile)

	inputBytes, err := io.ReadAll(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// I add a space after the string for the last word to be counted correctly, instead of adding all the word check/count logic after the loop again
	inputString := string(inputBytes) + " "

	// This slice is used to store the words and their frequencies in descending order by frequency
	wordFreq := make([]wordFreqEntry, 0)

	start = -1
	// Iterate over characters in the input file
	for i, c := range inputString {
		if start == -1 {
			// We're currently not in a word
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				// This means we found the start of a word
				start = i
			}
		} else {
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				// We're still inside a word
				continue
			}
			// When we reach this point, we're at the character immediately after a word

			// Copy the entire word and convert it to lowercase
			word := strings.ToLower(inputString[start:i])

			// Look for the word in the stopWords slice
			isStopWord := false
			for _, stopWord := range stopWords {
				if word == stopWord {
					isStopWord = true
					break
				}
			}

			if !isStopWord {
				// If the word is not a stop word, find it in the wordFreq slice
				idx := -1
				for i, wf := range wordFreq {
					if wf.word == word {
						idx = i
						break
					}
				}

				if idx == -1 {
					// The word is not the wordFreq slice so we append it to the slice with a frequency of 1
					wordFreq = append(wordFreq, wordFreqEntry{word, 1})
				} else {
					// The word is already in the wordFreq slice, so we increment its frequency
					wordFreq[idx].freq++
					// Then move it up the list until it's in the correct position again
					for idx > 0 && wordFreq[idx].freq > wordFreq[idx-1].freq {
						wordFreq[idx], wordFreq[idx-1] = wordFreq[idx-1], wordFreq[idx]
						idx--
					}
				}
			}
			// After we're done with a word, we want to look for the next one, so we reset the start index to -1
			start = -1
		}
	}

	// N is the number of words to be printed from the list, which is 25 words at most
	N := min(25, len(wordFreq))
	for _, wf := range wordFreq[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}
}
