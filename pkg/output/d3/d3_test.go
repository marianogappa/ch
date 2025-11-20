package d3

import (
	"flag"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

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

	// Test Render
	rows := make(chan ch.Row, 1)
	rows <- ch.Row{Floats: []float64{1.0}, Strings: []string{"a"}}
	close(rows)

	if err := o.Render(rows, cfg); err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
