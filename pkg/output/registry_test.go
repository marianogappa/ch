package output

import (
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
	_ "github.com/marianogappa/ch/pkg/output/chartjs"
	_ "github.com/marianogappa/ch/pkg/output/d3"
	_ "github.com/marianogappa/ch/pkg/output/json"
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
