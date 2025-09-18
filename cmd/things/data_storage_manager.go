package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
