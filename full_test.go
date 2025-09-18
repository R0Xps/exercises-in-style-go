package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOutputs(t *testing.T) {
	items, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Error reading root directory: %v", err)
	}

	inputFilesPath := filepath.Join("examples", "input")
	inputFiles, err := os.ReadDir(inputFilesPath)
	if err != nil {
		t.Fatalf("Error reading input directory: %v", err)
	}

	for _, inputFile := range inputFiles {
		inputFilePath := filepath.Join(inputFilesPath, inputFile.Name())
		for _, item := range items {
			if strings.HasPrefix(item.Name(), ".") || item.Name() == "examples" {
				continue
			}
			func() {
				if item.IsDir() {
					packagePath := "." + string(os.PathSeparator) + item.Name()
					args := []string{"run", packagePath, filepath.Join("examples", "stop_words.txt"), inputFilePath}
					if item.Name() == "persistent_tables" {
						dbFile := filepath.Join(os.TempDir(), inputFile.Name()+getRandomDBName()+".db")
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
						t.Fatalf("Error running %v on %v: %v\n", item.Name(), inputFile.Name(), err)
					}
					defer func(cmd *exec.Cmd) {
						err := cmd.Wait()
						if err != nil {
							t.Fatalf("Error running %v on %v: %v\n", item.Name(), inputFile.Name(), err)
						}
					}(cmd)

					stdOutBytes, err := io.ReadAll(stdoutPipe)
					if err != nil {
						t.Fatalf("Error reading stdout of %v on %v: %v\n", item.Name(), inputFile.Name(), err)
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
						t.Fatalf("Error reading stdout of sort on %v-%v: %v\n", item.Name(), inputFile.Name(), err)
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

					expectedOutputFile := filepath.Join("examples", "output", inputFile.Name())

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
