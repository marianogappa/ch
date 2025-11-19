package parser

import (
	"testing"

	"github.com/marianogappa/ch/pkg/ch"
)

func TestCSVParser(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      rune
		df       string
		expected []ch.Row
	}{
		{
			name:  "Simple CSV",
			input: []string{"1,hello", "2,world"},
			sep:   ',',
			expected: []ch.Row{
				{Floats: []float64{1}, Strings: []string{"hello"}},
				{Floats: []float64{2}, Strings: []string{"world"}},
			},
		},
		{
			name:  "With Date",
			input: []string{"2021-01-01,10.5"},
			sep:   ',',
			df:    "2006-01-02",
			expected: []ch.Row{
				{DateTimes: []string{"2021-01-01"}, Floats: []float64{10.5}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCSVParser(tt.sep, tt.df)
			in := make(chan []byte, len(tt.input))
			for _, l := range tt.input {
				in <- []byte(l)
			}
			close(in)

			out, err := p.Parse(in)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			var rows []ch.Row
			for r := range out {
				rows = append(rows, r)
			}

			if len(rows) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(rows))
			}

			// Deep check first row if exists
			if len(rows) > 0 {
				r := rows[0]
				e := tt.expected[0]
				if len(r.Floats) != len(e.Floats) {
					t.Errorf("Floats mismatch")
				}
				if len(r.Strings) != len(e.Strings) {
					t.Errorf("Strings mismatch")
				}
			}
		})
	}
}

func TestInferLineFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sep      rune
		df       string
		expected string
	}{
		{
			name:     "Floats",
			input:    "1.0,2.5,3",
			sep:      ',',
			expected: "fff",
		},
		{
			name:     "Strings",
			input:    "hello,world",
			sep:      ',',
			expected: "ss",
		},
		{
			name:     "Mixed",
			input:    "1.0,hello,2021-01-01",
			sep:      ',',
			df:       "2006-01-02",
			expected: "fsd",
		},
		{
			name:     "Mixed with spaces",
			input:    " 1.0 , hello , 2021-01-01 ",
			sep:      ',',
			df:       "2006-01-02",
			expected: "fsd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InferLineFormat(tt.input, tt.sep, tt.df)
			if got != tt.expected {
				t.Errorf("InferLineFormat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLineFormat_ParseLine(t *testing.T) {
	tests := []struct {
		name          string
		format        string
		sep           rune
		df            string
		input         string
		wantErr       bool
		expectFloats  []float64
		expectStrings []string
		expectDates   []string
	}{
		{
			name:         "Simple Floats",
			format:       "ff",
			sep:          ',',
			input:        "1.0,2.0",
			expectFloats: []float64{1.0, 2.0},
		},
		{
			name:          "Simple Strings",
			format:        "ss",
			sep:           ',',
			input:         "hello,world",
			expectStrings: []string{"hello", "world"},
		},
		{
			name:        "Simple Date",
			format:      "d",
			sep:         ',',
			df:          "2006-01-02",
			input:       "2021-01-01",
			expectDates: []string{"2021-01-01"},
		},
		{
			name:    "Invalid Float",
			format:  "f",
			sep:     ',',
			input:   "not_a_float",
			wantErr: true,
		},
		{
			name:    "Invalid Date",
			format:  "d",
			sep:     ',',
			df:      "2006-01-02",
			input:   "not_a_date",
			wantErr: true,
		},
		{
			name:    "Format Mismatch Length",
			format:  "ff",
			sep:     ',',
			input:   "1.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf, err := NewLineFormat(tt.format, tt.sep, tt.df)
			if err != nil {
				t.Fatalf("NewLineFormat error: %v", err)
			}

			fs, ss, ds, err := lf.ParseLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(fs) != len(tt.expectFloats) {
				t.Errorf("Floats length mismatch")
			}
			if len(ss) != len(tt.expectStrings) {
				t.Errorf("Strings length mismatch")
			}
			if len(ds) != len(tt.expectDates) {
				t.Errorf("Dates length mismatch")
			}
		})
	}
}
