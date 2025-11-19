package parser

import (
	"fmt"

	"github.com/marianogappa/ch/pkg/ch"
)

type CSVParser struct {
	Separator  rune
	DateFormat string
	// If empty, will try to infer
	LineFormat string
}

func NewCSVParser(separator rune, dateFormat string) *CSVParser {
	return &CSVParser{
		Separator:  separator,
		DateFormat: dateFormat,
	}
}

func (p *CSVParser) Parse(in <-chan []byte) (<-chan ch.Row, error) {
	out := make(chan ch.Row)

	go func() {
		defer close(out)

		var (
			buffer   []string
			lf       LineFormat
			err      error
			inferred bool
		)

		// If format is already known, use it.
		if p.LineFormat != "" {
			lf, err = NewLineFormat(p.LineFormat, p.Separator, p.DateFormat)
			if err != nil {
				// TODO: Handle error better?
				fmt.Printf("Error parsing format: %v\n", err)
				return
			}
			inferred = true
		}

		// Buffer for inference
		const inferenceLines = 5

		for lineBytes := range in {
			line := string(lineBytes)
			if !inferred {
				buffer = append(buffer, line)
				if len(buffer) >= inferenceLines {
					lf = p.infer(buffer)
					inferred = true
					// Process buffered lines
					for _, l := range buffer {
						p.emit(l, lf, out)
					}
					buffer = nil
				}
			} else {
				p.emit(line, lf, out)
			}
		}

		// If stream ended before inferenceLines, infer from what we have
		if !inferred && len(buffer) > 0 {
			lf = p.infer(buffer)
			for _, l := range buffer {
				p.emit(l, lf, out)
			}
		}
	}()

	return out, nil
}

func (p *CSVParser) infer(lines []string) LineFormat {
	counts := make(map[string]int)
	for _, l := range lines {
		fmtStr := InferLineFormat(l, p.Separator, p.DateFormat)
		counts[fmtStr]++
	}

	max := 0
	best := ""
	for k, v := range counts {
		if v > max {
			max = v
			best = k
		}
	}

	lf, _ := NewLineFormat(best, p.Separator, p.DateFormat)
	return lf
}

func (p *CSVParser) emit(line string, lf LineFormat, out chan<- ch.Row) {
	fs, ss, ds, err := lf.ParseLine(line)
	if err != nil {
		// Skip bad lines? Or log?
		return
	}
	out <- ch.Row{
		Floats:    fs,
		Strings:   ss,
		DateTimes: ds,
	}
}
