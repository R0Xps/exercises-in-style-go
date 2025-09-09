package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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
					dbFile := fmt.Sprintf(".db%v%v%v.db", string(os.PathSeparator), inputFile, getRandomDBName())
					args = append(args, dbFile)
				}

				stdOut := bytes.NewBuffer(nil)

				cmd := exec.Command("go", args...)
				cmd.Dir, _ = os.Getwd()
				cmd.Stdout = stdOut
				cmd.Stderr = stdOut

				err := cmd.Run()
				if err != nil {
					t.Error(cmd.Dir)
					t.Errorf("Error running %v on %v: %v\n", item.Name(), inputFile, err)
				}

				stdOutBytes, err := io.ReadAll(stdOut)
				if err != nil {
					t.Errorf("Error reading stdout of %v on %v: %v\n", item.Name(), inputFile, err)
				}

				f, err := os.CreateTemp(os.TempDir(), "test_"+item.Name())
				if err != nil {
					t.Errorf("Error creating temporary file: %v", err)
				}
				f.Write(stdOutBytes)
				f.Close()

				sortCmd := exec.Command("sort", f.Name())
				sortCmd.Stdout = stdOut
				sortCmd.Stderr = stdOut
				sortCmd.Run()

				stdOutBytes, err = io.ReadAll(stdOut)
				if err != nil {
					t.Errorf("Error reading stdout of sort on %v-%v: %v\n", item.Name(), inputFile, err)
				}

				f2, err := os.CreateTemp(os.TempDir(), "sorted_test_"+item.Name())
				if err != nil {
					t.Errorf("Error creating temporary file: %v", err)
				}
				f2.Write(stdOutBytes)
				f2.Close()

				expectedOutputFile := fmt.Sprintf(".output%v%v", string(os.PathSeparator), inputFile)
				diffOut := bytes.NewBuffer(nil)

				diffCmd := exec.Command("diff", "-u", f2.Name(), expectedOutputFile)
				diffCmd.Stdout = diffOut
				diffCmd.Stderr = diffOut

				err = diffCmd.Run()
				if err != nil {
					t.Error(diffCmd.Stdout, diffCmd.Stderr)
				}

				os.Remove(f.Name())
				os.Remove(f2.Name())
			}
		}
	}

}

func getRandomDBName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return base64.RawURLEncoding.EncodeToString(randBytes)
}
