package chartjs

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

func TestChartJSOutput(t *testing.T) {
	o := NewChartJSOutput()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := o.RegisterFlags(fs)

	// Test default config
	c, ok := cfg.(*ChartJSConfig)
	if !ok {
		t.Fatal("Config is not ChartJSConfig")
	}
	if c.ChartType != "line" {
		t.Errorf("Default ChartType should be line, got %s", c.ChartType)
	}

	// Mock openBrowser
	opened := ""
	oldOpenBrowser := openBrowser
	defer func() { openBrowser = oldOpenBrowser }()
	openBrowser = func(url string) error {
		opened = url
		return nil
	}

	// Test Render
	rows := make(chan ch.Row, 1)
	rows <- ch.Row{Floats: []float64{1.0}, Strings: []string{"a"}, DateTimes: []string{"2021-01-01"}}
	close(rows)

	if err := o.Render(rows, cfg); err != nil {
		t.Errorf("Render failed: %v", err)
	}

	if opened == "" {
		t.Error("Expected browser to be opened")
	}
}

func TestChartJSOutput_Frequency(t *testing.T) {
	o := NewChartJSOutput()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := o.RegisterFlags(fs)

	// Mock openBrowser
	opened := ""
	oldOpenBrowser := openBrowser
	defer func() { openBrowser = oldOpenBrowser }()
	openBrowser = func(url string) error {
		opened = url
		return nil
	}

	// Test Render with only strings
	rows := make(chan ch.Row, 3)
	rows <- ch.Row{Strings: []string{"apple"}}
	rows <- ch.Row{Strings: []string{"banana"}}
	rows <- ch.Row{Strings: []string{"apple"}}
	close(rows)

	if err := o.Render(rows, cfg); err != nil {
		t.Errorf("Render failed: %v", err)
	}

	if opened == "" {
		t.Fatal("Expected browser to be opened")
	}

	// Read the generated file
	content, err := os.ReadFile(opened)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	html := string(content)

	// Check for expected data
	// We expect "apple" count 2, "banana" count 1
	if !strings.Contains(html, "apple") {
		t.Error("Expected 'apple' in output")
	}
	if !strings.Contains(html, "banana") {
		t.Error("Expected 'banana' in output")
	}
	// Check for counts (floats)
	if !strings.Contains(html, "2") {
		t.Error("Expected count 2 in output")
	}
}
