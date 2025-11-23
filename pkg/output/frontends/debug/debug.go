package debug

import (
	"encoding/json"
	"io"

	"github.com/marianogappa/ch/pkg/ch"
)

// Debug is the main struct for debug output
type Debug struct {
	pretty bool
}

// New creates a new Debug instance
func New(pretty bool) *Debug {
	return &Debug{
		pretty: pretty,
	}
}

// Encode encodes a row to the writer
func (d *Debug) Encode(row ch.Row, w io.Writer) error {
	enc := json.NewEncoder(w)
	if d.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(row)
}

