package main

import (
	"bytes"
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
			if item.IsDir() {
				packagePath := "." + string(os.PathSeparator) + item.Name()
				args := []string{"run", packagePath, "stop_words.txt", inputFile}
				if item.Name() == "persistent_tables" {
					dbFile := fmt.Sprintf("%v%v%v%v.db", os.TempDir(), string(os.PathSeparator), inputFile, getRandomDBName())
					args = append(args, dbFile)
				}

				cmd := exec.Command("go", args...)
				stdoutPipe, err := cmd.StdoutPipe()
				if err != nil {
					t.Fatalf("Error opening stdout pipe: %v", err)
				}
				cmd.Dir, _ = os.Getwd()

				err = cmd.Start()
				if err != nil {
					t.Fatalf("Error running %v on %v: %v\n", item.Name(), inputFile, err)
				}

				stdOutBytes, err := io.ReadAll(stdoutPipe)
				if err != nil {
					t.Errorf("Error reading stdout of %v on %v: %v\n", item.Name(), inputFile, err)
				}

				err = cmd.Wait()
				if err != nil {
					t.Errorf("Error running command 'go run': %v\n", err)
				}

				f, err := os.CreateTemp(os.TempDir(), "test_"+item.Name())
				if err != nil {
					t.Errorf("Error creating temporary file: %v", err)
				}
				_, err = f.Write(stdOutBytes)
				if err != nil {
					t.Fatal(err)
					return
				}
				err = f.Close()
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
					t.Fatal(err)
					return
				}

				stdOutBytes, err = io.ReadAll(stdoutPipe)
				if err != nil {
					t.Errorf("Error reading stdout of sort on %v-%v: %v\n", item.Name(), inputFile, err)
				}

				err = sortCmd.Wait()
				if err != nil {
					t.Errorf("Error running command 'sort': %v", err)
				}

				f2, err := os.CreateTemp(os.TempDir(), "sorted_test_"+item.Name())
				if err != nil {
					t.Errorf("Error creating temporary file: %v", err)
				}
				_, err = f2.Write(stdOutBytes)
				if err != nil {
					t.Fatal(err)
					return
				}
				err = f2.Close()
				if err != nil {
					t.Fatal(err)
					return
				}

				expectedOutputFile := fmt.Sprintf(".output%v%v", string(os.PathSeparator), inputFile)
				diffOut := bytes.NewBuffer(nil)

				diffCmd := exec.Command("diff", "-u", f2.Name(), expectedOutputFile)
				diffCmd.Stdout = diffOut
				diffCmd.Stderr = diffOut

				err = diffCmd.Run()
				if err != nil {
					t.Error(diffCmd.Stdout, diffCmd.Stderr)
				}

				err = os.Remove(f.Name())
				if err != nil {
					t.Fatal(err)
					return
				}
				err = os.Remove(f2.Name())
				if err != nil {
					t.Fatal(err)
					return
				}
			}
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
