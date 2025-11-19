package output

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

func TestRegistry(t *testing.T) {
	outputs := ch.Outputs()
	if len(outputs) == 0 {
		t.Fatal("No outputs registered")
	}

	found := false
	for _, name := range outputs {
		if name == "json" {
			found = true
			break
		}
	}
	if !found {
		t.Error("json output not found in registry")
	}
}

func TestJSONOutput(t *testing.T) {
	o := NewJSONOutput()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := o.RegisterFlags(fs)

	// Test default config
	c, ok := cfg.(*JSONConfig)
	if !ok {
		t.Fatal("Config is not JSONConfig")
	}
	if c.Pretty {
		t.Error("Default Pretty should be false")
	}

	// We can't easily test Render to stdout without capturing it,
	// but we can ensure it doesn't panic on empty channel
	rows := make(chan ch.Row)
	close(rows)
	if err := o.Render(rows, cfg); err != nil {
		t.Errorf("Render failed: %v", err)
	}
}

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

func TestD3Output(t *testing.T) {
	o := NewD3Output()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := o.RegisterFlags(fs)

	// Test default config
	c, ok := cfg.(*D3Config)
	if !ok {
		t.Fatal("Config is not D3Config")
	}
	if c.Title != "D3 Chart" {
		t.Errorf("Default Title should be 'D3 Chart', got %s", c.Title)
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
	rows <- ch.Row{Floats: []float64{1.0}, Strings: []string{"a"}}
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
