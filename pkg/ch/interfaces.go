package ch

import (
	"flag"
	"fmt"
	"sort"
	"sync"
)

// Input represents a source of data.
// It streams data in chunks (e.g., lines or bytes).
type Input interface {
	Stream() (<-chan []byte, error)
}

// Row represents a single data point with mixed types.
// It corresponds to a parsed line of input.
type Row struct {
	Floats    []float64
	Strings   []string
	DateTimes []string
}

// Parser interprets the raw input stream into structured Rows.
type Parser interface {
	Parse(<-chan []byte) (<-chan Row, error)
}

// Capabilities defines what an Output can do.
type Capabilities struct {
	Streaming   bool
	Interactive bool
}

// Output represents a destination for the charted data.
type Output interface {
	Name() string
	// RegisterFlags registers the flags for this output on the given FlagSet.
	// It returns a pointer to the configuration struct that will be populated when flags are parsed.
	RegisterFlags(fs *flag.FlagSet) any
	// Render renders the data using the given configuration.
	// The config argument is the same pointer returned by RegisterFlags.
	Render(rows <-chan Row, config any) error
	Capabilities() Capabilities
}

var (
	outputsMu sync.RWMutex
	outputs   = make(map[string]Output)
)

// RegisterOutput registers an output driver.
func RegisterOutput(driver Output) {
	outputsMu.Lock()
	defer outputsMu.Unlock()
	if driver == nil {
		panic("ch: RegisterOutput driver is nil")
	}
	name := driver.Name()
	if _, dup := outputs[name]; dup {
		panic("ch: RegisterOutput called twice for driver " + name)
	}
	outputs[name] = driver
}

// GetOutput returns an output driver by name.
func GetOutput(name string) (Output, error) {
	outputsMu.RLock()
	defer outputsMu.RUnlock()
	driver, ok := outputs[name]
	if !ok {
		return nil, fmt.Errorf("ch: unknown output driver %q", name)
	}
	return driver, nil
}

// Outputs returns a sorted list of the names of the registered outputs.
func Outputs() []string {
	outputsMu.RLock()
	defer outputsMu.RUnlock()
	var list []string
	for name := range outputs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
