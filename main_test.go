package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

func TestRun_JSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	// Input data
	inputData := "1.0,hello\n2.0,world\n"
	stdin := strings.NewReader(inputData)

	// Args
	args := []string{"ch", "--output", "json", "--format", "fs", "--separator", ","}

	// Run
	err := Run(args, stdin)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close write end of pipe to read from it
	w.Close()

	// Read stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify output
	// The output should be a sequence of JSON objects
	decoder := json.NewDecoder(strings.NewReader(output))
	var rows []ch.Row
	for {
		var row ch.Row
		if err := decoder.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("Failed to decode JSON output: %v", err)
		}
		rows = append(rows, row)
	}

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	if len(rows) > 0 {
		if len(rows[0].Floats) != 1 || rows[0].Floats[0] != 1.0 {
			t.Errorf("Row 1 float mismatch")
		}
		if len(rows[0].Strings) != 1 || rows[0].Strings[0] != "hello" {
			t.Errorf("Row 1 string mismatch")
		}
	}
}
