package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/input"
	"github.com/marianogappa/ch/pkg/output"
	"github.com/marianogappa/ch/pkg/parser"
)

func TestRun_JSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore stdout
	defer func() {
		os.Stdout = oldStdout
		w.Close()
	}()

	// Input data
	inputData := "1.0,hello\n2.0,world\n"
	stdin := strings.NewReader(inputData)

	// Setup Input
	in := input.NewReaderInput(stdin)

	// Setup Parser
	p := parser.NewCSVParser(',', "")
	p.LineFormat = "fs"

	// Get output driver
	outDriver, err := ch.GetOutput("json")
	if err != nil {
		t.Fatalf("Failed to get json output driver: %v", err)
	}

	// Check if driver implements ChartOutput
	chartOutput, isChartOutput := outDriver.(output.ChartOutput)
	if !isChartOutput {
		t.Fatalf("JSON driver does not implement ChartOutput")
	}

	// Create chart config
	chartConfig := &output.ChartConfig{
		OutputPath: "-", // Use stdout
	}

	// Run
	stream, err := in.Stream()
	if err != nil {
		t.Fatalf("Failed to create input stream: %v", err)
	}

	rows, err := p.Parse(stream)
	if err != nil {
		t.Fatalf("Failed to parse input: %v", err)
	}

	// Render using ChartConfig
	if err := chartOutput.RenderChart(rows, chartConfig); err != nil {
		t.Fatalf("Failed to render output: %v", err)
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
	var decodedRows []ch.Row
	for {
		var row ch.Row
		if err := decoder.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("Failed to decode JSON output: %v", err)
		}
		decodedRows = append(decodedRows, row)
	}

	if len(decodedRows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(decodedRows))
	}

	if len(decodedRows) > 0 {
		if len(decodedRows[0].Floats) != 1 || decodedRows[0].Floats[0] != 1.0 {
			t.Errorf("Row 1 float mismatch")
		}
		if len(decodedRows[0].Strings) != 1 || decodedRows[0].Strings[0] != "hello" {
			t.Errorf("Row 1 string mismatch")
		}
	}
}
