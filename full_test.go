package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestOutputs(t *testing.T) {
	items, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Error reading directory: %v", err)
	}

	inputFiles := []string{"input1.txt", "input2.txt", "pride-and-prejudice.txt"}
	for _, inputFile := range inputFiles {
		for _, item := range items {
			if strings.HasPrefix(item.Name(), ".") {
				continue
			}
			func() {
				if item.IsDir() {
					packagePath := "." + string(os.PathSeparator) + item.Name()
					args := []string{"run", packagePath, "stop_words.txt", inputFile}
					if item.Name() == "persistent_tables" {
						dbFile := fmt.Sprintf("%v%v%v%v.db", os.TempDir(), string(os.PathSeparator), inputFile, getRandomDBName())
						args = append(args, dbFile)
						defer func(name string) {
							err := os.Remove(name)
							if err != nil {
								t.Fatalf("Error removing temp file: %v", err)
							}
						}(dbFile)
					}

					cmd := exec.Command("go", args...)
					stdoutPipe, err := cmd.StdoutPipe()
					if err != nil {
						t.Fatalf("Error opening stdout pipe: %v", err)
					}

					err = cmd.Start()
					if err != nil {
						t.Fatalf("Error running %v on %v: %v\n", item.Name(), inputFile, err)
					}
					defer func(cmd *exec.Cmd) {
						err := cmd.Wait()
						if err != nil {
							t.Fatalf("Error running %v on %v: %v\n", item.Name(), inputFile, err)
						}
					}(cmd)

					stdOutBytes, err := io.ReadAll(stdoutPipe)
					if err != nil {
						t.Fatalf("Error reading stdout of %v on %v: %v\n", item.Name(), inputFile, err)
					}

					f, err := os.CreateTemp(os.TempDir(), "test_"+item.Name())
					if err != nil {
						t.Fatalf("Error creating temporary file: %v", err)
					}
					defer func(name string) {
						err := os.Remove(name)
						if err != nil {
							t.Fatalf("Error removing temp file: %v", err)
						}
					}(f.Name())
					defer func(f *os.File) {
						err := f.Close()
						if err != nil {
							t.Fatalf("Error closing temporary file: %v", err)
						}
					}(f)
					_, err = f.Write(stdOutBytes)
					if err != nil {
						t.Fatal(err)
						return
					}

					sortCmd := exec.Command("sort", f.Name())
					stdoutPipe, err = sortCmd.StdoutPipe()
					if err != nil {
						t.Fatalf("Error opening stdout pipe: %v", err)
					}

					err = sortCmd.Start()
					if err != nil {
						t.Fatalf("Error running sort command: %v", err)
						return
					}
					defer func(sortCmd *exec.Cmd) {
						err := sortCmd.Wait()
						if err != nil {
							t.Fatalf("Error running sort command: %v", err)
						}
					}(sortCmd)

					stdOutBytes, err = io.ReadAll(stdoutPipe)
					if err != nil {
						t.Fatalf("Error reading stdout of sort on %v-%v: %v\n", item.Name(), inputFile, err)
					}

					f2, err := os.CreateTemp(os.TempDir(), "sorted_test_"+item.Name())
					if err != nil {
						t.Fatalf("Error creating temporary file: %v", err)
					}
					defer func(name string) {
						err := os.Remove(name)
						if err != nil {
							t.Fatalf("Error removing temp file: %v", err)
						}
					}(f2.Name())
					defer func(f2 *os.File) {
						err := f2.Close()
						if err != nil {
							t.Fatalf("Error closing temporary file: %v", err)
						}
					}(f2)

					_, err = f2.Write(stdOutBytes)
					if err != nil {
						t.Fatal(err)
						return
					}

					expectedOutputFile := fmt.Sprintf(".output%v%v", string(os.PathSeparator), inputFile)

					diffCmd := exec.Command("diff", "-u", f2.Name(), expectedOutputFile)
					stdoutPipe, err = diffCmd.StdoutPipe()
					if err != nil {
						t.Fatalf("Error opening stdout pipe: %v", err)
					}

					err = diffCmd.Start()
					if err != nil {
						t.Fatal(stdoutPipe)
					}
					defer func(diffCmd *exec.Cmd) {
						err := diffCmd.Wait()
						if err != nil {
							t.Fatalf("Error running diff command: %v", err)
						}
					}(diffCmd)
				}
			}()
		}
	}
}

func getRandomDBName() string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(randBytes)
}
