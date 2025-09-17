package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

// wordFreqEntry struct is used to store a word-frequency pair
type wordFreqEntry struct {
	word string
	freq int
}

func main() {
	// Check for the required arguments
	if len(os.Args) != 4 {
		log.Fatal("required arguments: <stop_words_file> <input_file> <database_file>")
	}

	stopWordsFile := os.Args[1]
	inputFile := os.Args[2]
	dbFile := os.Args[3]

	// Connect to sqlite database
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// If the database file doesn't exist, create the tables and insert the data into them (automatically creates the database file)
	exists, err := fileExists(dbFile)
	if !exists {
		if err != nil {
			log.Fatal(err)
		}

		createTables(db)
		insertStopWords(db, stopWordsFile)
		insertData(db, inputFile)
	}

	// Get all words and their frequencies (25 words max)
	rows, err := db.Query("SELECT word, COUNT(*) AS freq FROM words GROUP BY word ORDER BY freq DESC LIMIT 25")
	if err != nil {
		log.Fatal("Error retrieving words and their frequencies from database:", err)
	}

	wordFreq := make([]wordFreqEntry, 0, 25)
	for rows.Next() {
		wordFreqEntry := wordFreqEntry{}
		err = rows.Scan(&wordFreqEntry.word, &wordFreqEntry.freq)
		if err != nil {
			log.Fatal("Error retrieving words and their frequencies from database:", err)
		}
		wordFreq = append(wordFreq, wordFreqEntry)
	}

	// Print all words and their frequencies from the slice
	for _, wf := range wordFreq {
		fmt.Println(wf.word, "-", wf.freq)
	}
}

// Check if a file exists
func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Create the required tables in the database
func createTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE documents (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		log.Fatal("Error creating documents table:", err)
	}
	_, err = db.Exec("CREATE TABLE words (id INTEGER PRIMARY KEY, doc_id INTEGER, word TEXT, FOREIGN KEY(doc_id) REFERENCES documents(id))")
	if err != nil {
		log.Fatal("Error creating words table:", err)
	}
	_, err = db.Exec("CREATE TABLE stop_words (word TEXT PRIMARY KEY)")
	if err != nil {
		log.Fatal("Error creating stop words table:", err)
	}
}

// Insert the words from the stop words file into the stop_words table
func insertStopWords(db *sql.DB, stopWordsFile string) {
	file, err := os.Open(filepath.Clean(stopWordsFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	stopWords := strings.Fields(filterAndNormalize(bytes))
	for _, word := range stopWords {
		_, err = db.Exec("INSERT INTO stop_words (word) VALUES (?)", word)
		if err != nil {
			log.Fatal("Error inserting stop words into database:", err)
		}
	}
}

// Insert the words from the input file into the words table, along with a new entry in the documents table referring to the input file itself
func insertData(db *sql.DB, inputFile string) {
	file, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	filteredInput := filterAndNormalize(bytes)
	words := strings.Fields(filteredInput)

	_, err = db.Exec("INSERT INTO documents (name) VALUES (?)", inputFile)
	if err != nil {
		log.Fatal("Error inserting new document into database:", err)
	}
	var docId int
	err = db.QueryRow("SELECT MAX(id) FROM documents WHERE name=?", inputFile).Scan(&docId)
	if err != nil {
		log.Fatal("Error getting new document id:", err)
	}

	var stopWords []string
	rows, err := db.Query("SELECT word FROM stop_words")
	if err != nil {
		log.Fatal("Error retrieving stop words from database:", err)
	}

	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			log.Fatal("Error retrieving stop words from database:", err)
		}
		stopWords = append(stopWords, word)
	}

	var wordId int
	err = db.QueryRow("SELECT MAX(id) FROM words").Scan(&wordId)
	if err != nil {
		log.Fatal("Error getting new word id:", err)
	}
	wordId++
	for _, word := range words {
		if slices.Contains(stopWords, word) {
			continue
		}

		_, err = db.Exec("INSERT INTO words (id, doc_id, word) values (?, ?, ?)", wordId, docId, word)
		if err != nil {
			log.Fatal("Error inserting data:", err)
			return
		}
		wordId++
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
