package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

// Slice used to store stop words after reading them from a file
var stopWords []string

func main() {
	// Check for the required arguments
	if len(os.Args) != 3 {
		log.Fatal("required arguments: <stop_words_file> <input_file>")
	}

	stopWords = getStopWords(os.Args[1])

	data := readInputFile(os.Args[2])
	parts := lop.Map(partition(data, 200), splitWords)
	wfMap := lo.Reduce(parts, countWords, map[string]int{})

	wordFreq := sorted(wfMap)

	// Print the first (25 max) words and their frequencies
	N := min(25, len(wordFreq))
	for _, wf := range wordFreq[:N] {
		fmt.Println(wf.word, "-", wf.freq)
	}
}

func getStopWords(filename string) []string {
	rawStopWords := readInputFile(filename)
	filteredStopWords := filterAndNormalize([]byte(rawStopWords))
	return strings.Fields(filteredStopWords)
}

// Read the input file and return its content as a string
func readInputFile(filename string) string {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

// Split the data string into smaller strings of nLines lines each (the last string might have less than nLines lines)
func partition(data string, nLines int) []string {
	lines := strings.Split(data, "\n")

	parts := make([]string, 0, len(lines)/nLines+1)
	for i := 0; i < len(lines); i += nLines {
		partEnd := min(len(lines), i+nLines)
		parts = append(parts, strings.Join(lines[i:partEnd], "\n"))
	}
	return parts
}

// This is the 'map' function of this MapReduce job. It takes a string, cleans it by replacing all non-letter characters with spaces, and converts all uppercase letters to lowercase.
// Then it splits the resulting string, leaving only the words. And returns a slice of all non-stop words with a frequency of 1 for each of them (repeats allowed)
func splitWords(data string, _ int) []wordFreqEntry {
	cleanData := filterAndNormalize([]byte(data))
	words := strings.Fields(cleanData)
	wordFreq := make([]wordFreqEntry, 0)

	for _, word := range words {
		if !isStopWord(word) {
			wordFreq = append(wordFreq, wordFreqEntry{word, 1})
		}
	}

	return wordFreq
}

// Check if the given word is in the stopWords slice
func isStopWord(word string) bool {
	return slices.Contains(stopWords, word)
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

// This is the 'reduce' function of this 'MapReduce' job. It combines all words and frequencies from the item slice into the agg map and returns the map
func countWords(agg map[string]int, item []wordFreqEntry, _ int) map[string]int {
	for _, wf := range item {
		agg[wf.word] += wf.freq
	}
	return agg
}

// Return a sorted slice of all entries in wfMap
func sorted(wfMap map[string]int) []wordFreqEntry {
	wordFreq := make([]wordFreqEntry, 0)
	for k, v := range wfMap {
		wordFreq = append(wordFreq, wordFreqEntry{k, v})
	}

	slices.SortFunc(wordFreq, func(i, j wordFreqEntry) int {
		return j.freq - i.freq
	})

	return wordFreq
}
