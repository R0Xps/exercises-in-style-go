package main

import "slices"

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
