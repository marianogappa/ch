package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
	"github.com/marianogappa/ch/pkg/input"
	"github.com/marianogappa/ch/pkg/output"
	_ "github.com/marianogappa/ch/pkg/output/drivers/chartjs"
	_ "github.com/marianogappa/ch/pkg/output/drivers/d3"
	_ "github.com/marianogappa/ch/pkg/output/drivers/debug"
	"github.com/marianogappa/ch/pkg/parser"
)

// TestOutputIntegration tests all chart types for all output drivers
// It generates output files with deterministic names for visual inspection
func TestOutputIntegration(t *testing.T) {
	// Define test cases: driver name, chart type, testdata file
	testCases := []struct {
		driverName string
		chartType  output.ChartType
		testData   string
		skip       bool // Skip if chart type not supported
	}{
		// D3 driver tests
		{"d3", output.ChartTypeLine, "line.csv", false},
		{"d3", output.ChartTypeBar, "bar.csv", false},
		{"d3", output.ChartTypePie, "pie.csv", false},
		{"d3", output.ChartTypeScatter, "scatter.csv", false},
		{"d3", output.ChartTypeHistogram, "freq.csv", false},
		{"d3", output.ChartTypeDoughnut, "", true},  // Not supported
		{"d3", output.ChartTypeBubble, "", true},    // Not supported
		{"d3", output.ChartTypeRadar, "", true},     // Not supported
		{"d3", output.ChartTypePolarArea, "", true}, // Not supported

		// ChartJS driver tests
		{"chartjs", output.ChartTypeLine, "line.csv", false},
		{"chartjs", output.ChartTypeBar, "bar.csv", false},
		{"chartjs", output.ChartTypePie, "pie.csv", false},
		{"chartjs", output.ChartTypeDoughnut, "pie.csv", false},
		{"chartjs", output.ChartTypeScatter, "scatter.csv", false},
		{"chartjs", output.ChartTypeBubble, "scatter.csv", false}, // Reuse scatter data
		{"chartjs", output.ChartTypeRadar, "scatter-rad.csv", false},
		{"chartjs", output.ChartTypePolarArea, "pie.csv", false},
		{"chartjs", output.ChartTypeHistogram, "", true}, // Not supported

		// Debug driver tests (Debug doesn't have chart types, but we test it with different data)
		{"debug", output.ChartTypeLine, "line.csv", false}, // ChartType is ignored for Debug
	}

	// Create test output directory
	testOutputDir := "test-output"
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		t.Fatalf("Failed to create test output directory: %v", err)
	}

	for _, tc := range testCases {
		if tc.skip {
			continue
		}

		t.Run(tc.driverName+"-"+string(tc.chartType), func(t *testing.T) {
			// Get the output driver
			driver, err := ch.GetOutput(tc.driverName)
			if err != nil {
				t.Fatalf("Failed to get output driver %s: %v", tc.driverName, err)
			}

			// Check if driver implements ChartOutput
			chartOutput, isChartOutput := driver.(output.ChartOutput)
			if !isChartOutput {
				t.Fatalf("Driver %s does not implement ChartOutput", tc.driverName)
			}

			// Check if chart type is supported (Debug ignores chart types)
			// For now, we'll let the driver validate this during rendering
			// If a chart type is not supported, the render will fail with an error

			// Determine output file path
			var outputPath string
			if tc.driverName == "debug" {
				// For Debug, use testdata filename instead of chart type
				outputPath = filepath.Join(testOutputDir, tc.driverName+"-"+tc.testData)
				// Remove .csv and add .json
				outputPath = outputPath[:len(outputPath)-4] + ".json"
			} else {
				outputPath = filepath.Join(testOutputDir, tc.driverName+"-"+string(tc.chartType)+".html")
			}

			// Setup input from testdata
			testDataPath := filepath.Join("testdata", tc.testData)
			fileInput := input.NewFileInput(testDataPath)
			stream, err := fileInput.Stream()
			if err != nil {
				t.Fatalf("Failed to create input stream: %v", err)
			}

			// Setup parser
			p := parser.NewCSVParser('\t', "")
			rows, err := p.Parse(stream)
			if err != nil {
				t.Fatalf("Failed to parse input: %v", err)
			}

			// Create chart config
			config := &output.ChartConfig{
				ChartType:  tc.chartType,
				Title:      "Test " + string(tc.chartType) + " Chart",
				XLabel:     "X Axis",
				YLabel:     "Y Axis",
				OutputPath: outputPath,
			}

			// Render chart
			if err := chartOutput.RenderChart(rows, config); err != nil {
				t.Fatalf("Failed to render chart: %v", err)
			}

			// Verify output file was created
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Output file was not created: %s", outputPath)
			} else if err != nil {
				t.Errorf("Error checking output file: %v", err)
			}
		})
	}
}

// TestLLMBlobOutput tests that the LLM blob can be generated and written to test-output
func TestLLMBlobOutput(t *testing.T) {
	// Generate LLM blob using the new ChartConfig documentation format
	blob, err := output.GenerateChartConfigDocumentation()
	if err != nil {
		t.Fatalf("Failed to generate LLM blob: %v", err)
	}

	// Create test output directory
	testOutputDir := "test-output"
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		t.Fatalf("Failed to create test output directory: %v", err)
	}

	// Write blob to test-output
	outputPath := filepath.Join(testOutputDir, "llm-blob.json")
	if err := os.WriteFile(outputPath, []byte(blob), 0644); err != nil {
		t.Fatalf("Failed to write LLM blob to file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("LLM blob file was not created: %s", outputPath)
	} else if err != nil {
		t.Errorf("Error checking LLM blob file: %v", err)
	}

	// Verify blob is valid JSON (basic check)
	if len(blob) == 0 {
		t.Error("LLM blob is empty")
	}

	// Check that blob contains expected keys for ChartConfig documentation
	if !contains(blob, "description") {
		t.Error("LLM blob missing 'description'")
	}
	if !contains(blob, "example") {
		t.Error("LLM blob missing 'example'")
	}
	if !contains(blob, "fields") {
		t.Error("LLM blob missing 'fields'")
	}
	if !contains(blob, "frontendSettings") {
		t.Error("LLM blob missing 'frontendSettings'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
