package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	// Create a new quarantine object, bind all functions to it, then execute them in order
	NewQuarantine(getInput).Bind(extractWords).Bind(removeStopWords).Bind(frequencies).Bind(sort).Bind(top25).Execute()
}

type Quarantine struct {
	functions []func(any) any
}

// Create and return a pointer to a new Quaranitine object with f as the first entry in its functions slice
func NewQuarantine(f func(any) any) *Quarantine {
	return &Quarantine{
		functions: []func(any) any{f},
	}
}

// Add another function to q's functions slice
func (q *Quarantine) Bind(f func(any) any) *Quarantine {
	q.functions = append(q.functions, f)
	return q
}

// Execute the functions in q's functions slice in order, giving the output of each function as input to the next
func (q *Quarantine) Execute() {
	guardFunc := func(v any) any {
		f, ok := v.(func() any)
		if !ok {
			return v
		}
		return f()
	}
	var val any = nil
	for _, f := range q.functions {
		val = guardFunc(f(val))
	}
}

// Return a function that returns the path to the input file
func getInput(_ any) any {
	return func() any {
		// Check for the required arguments
		if len(os.Args) != 3 {
			log.Fatal("required arguments: <stop_words_file> <input_file>")
		}

		return os.Args[2]
	}
}

// Return a function that returns a slice of strings containing all words from the file at filePath
func extractWords(filePath any) any {
	return func() any {
		file, err := os.Open(filePath.(string))
		if err != nil {
			log.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		words := strings.Fields(filterAndNormalize(bytes))
		return words
	}
}

// Return a function that returns a slice of string containing all non-stop words from the given words slice
func removeStopWords(words any) any {
	return func() any {
		allWords := words.([]string)
		stopWords := extractWords(os.Args[1]).(func() any)().([]string)
		nonStopWords := make([]string, 0)
		for _, word := range allWords {
			if !slices.Contains(stopWords, word) {
				nonStopWords = append(nonStopWords, word)
			}
		}
		return nonStopWords
	}
}

// Return a map containing all words from the words slice with their frequencies
func frequencies(words any) any {
	wordsSlice := words.([]string)
	wfMap := make(map[string]int)
	for _, word := range wordsSlice {
		wfMap[word]++
	}
	return wfMap
}

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

// Return a sorted slice containing all entries from the wf map
func sort(wf any) any {
	wfMap := wf.(map[string]int)
	wordFreq := make([]wordFreqEntry, 0)
	for k, v := range wfMap {
		wordFreq = append(wordFreq, wordFreqEntry{k, v})
	}

	slices.SortFunc(wordFreq, func(i, j wordFreqEntry) int {
		return j.freq - i.freq
	})

	return wordFreq
}

// Print the first 25 (or less if the slice is shorter than 25) elements in the wordFreq slice
func top25(wordFreq any) any {
	wordFreqSlice := wordFreq.([]wordFreqEntry)
	N := min(25, len(wordFreqSlice))
	for _, wf := range wordFreqSlice[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}
	return nil
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
