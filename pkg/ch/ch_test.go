package ch

import (
	"flag"
	"testing"
)

type mockOutput struct {
	name string
}

func (m *mockOutput) Name() string                             { return m.name }
func (m *mockOutput) RegisterFlags(fs *flag.FlagSet) any       { return nil }
func (m *mockOutput) Render(rows <-chan Row, config any) error { return nil }
func (m *mockOutput) Capabilities() Capabilities               { return Capabilities{} }

func TestRegistry(t *testing.T) {
	// Note: This test runs in the same process as other tests, so the registry might already be populated.
	// We should test adding a new one.

	name := "test_output"
	m := &mockOutput{name: name}

	// Register
	RegisterOutput(m)

	// Get
	got, err := GetOutput(name)
	if err != nil {
		t.Fatalf("GetOutput failed: %v", err)
	}
	if got != m {
		t.Errorf("Expected %v, got %v", m, got)
	}

	// List
	list := Outputs()
	found := false
	for _, n := range list {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Output %q not found in list: %v", name, list)
	}

	// Duplicate panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on duplicate registration")
		}
	}()
	RegisterOutput(m)
}
