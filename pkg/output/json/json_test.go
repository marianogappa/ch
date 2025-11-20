package json

import (
	"flag"
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

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
