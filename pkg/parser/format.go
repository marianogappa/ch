package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LineFormat represents the format of a line of input
type LineFormat struct {
	ColTypes   []ColType
	Separator  rune
	DateFormat string

	HasFloats     bool
	HasStrings    bool
	HasDateTimes  bool
	FloatCount    int
	StringCount   int
	DateTimeCount int
}

// ColType represents the type of a column in a data point
type ColType int

const (
	String ColType = iota
	Float
	DateTime
)

func (c ColType) String() string {
	switch c {
	case String:
		return "s"
	case Float:
		return "f"
	case DateTime:
		return "d"
	default:
		return "?"
	}
}

func (l LineFormat) String() string {
	var bs = make([]byte, 0)
	for _, c := range l.ColTypes {
		bs = append(bs, c.String()...)
	}
	return string(bs)
}

// NewLineFormat creates a LineFormat from a string string
func NewLineFormat(lineFormat string, separator rune, dateFormat string) (LineFormat, error) {
	if ok, err := regexp.Match("[dfs ]*", []byte(lineFormat)); !ok || err != nil {
		return LineFormat{}, fmt.Errorf("format: supplied lineFormat doesn't match syntax `[dfs ]*`")
	}
	var lf = LineFormat{ColTypes: nil, Separator: separator, DateFormat: dateFormat}

	for _, b := range lineFormat {
		switch b {
		case 's':
			lf.ColTypes = append(lf.ColTypes, String)
			lf.StringCount++
		case 'f':
			lf.ColTypes = append(lf.ColTypes, Float)
			lf.FloatCount++
		case 'd':
			lf.ColTypes = append(lf.ColTypes, DateTime)
			lf.DateTimeCount++
		default:
		}
	}
	lf.HasStrings = lf.StringCount > 0
	lf.HasFloats = lf.FloatCount > 0
	lf.HasDateTimes = lf.DateTimeCount > 0
	return lf, nil
}

// ParseLine parses one line of input according to the given format
func (l LineFormat) ParseLine(line string) ([]float64, []string, []string, error) {
	// Note: Changing return type of DateTimes to []string to match ch.Row definition for now,
	// or we can parse to time.Time and convert later.
	// The interface said []string for DateTimes in Row, but maybe it should be time.Time?
	// The original code used time.Time.
	// Let's stick to string in Row for maximum flexibility in transport, but here we validate it parses.

	// Actually, for the Row struct in interfaces.go, I defined DateTimes as []string.
	// But to be useful for charting, we might want time.Time.
	// However, JSON marshaling time.Time is standard.
	// Let's keep it as string in Row but ensure it's a valid date here.

	line = string(regexp.MustCompile(string(l.Separator)+"{2,}").ReplaceAll([]byte(line), []byte(string(l.Separator))))
	sp := strings.Split(strings.TrimSpace(line), string(l.Separator))

	fs := []float64{}
	ss := []string{}
	ds := []string{}

	if len(sp) < len(l.ColTypes) {
		return fs, ss, ds, fmt.Errorf("Input line has invalid format length; expected %v vs found %v", len(l.ColTypes), len(sp))
	}

	for i, colType := range l.ColTypes {
		s := strings.TrimSpace(sp[i])
		switch colType {
		case String:
			ss = append(ss, s)
		case DateTime:
			_, err := time.Parse(l.DateFormat, s)
			if err != nil {
				return fs, ss, ds, fmt.Errorf("Couldn't convert %v to date given: %v", s, err)
			}
			ds = append(ds, s) // Keep as string
		case Float:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fs, ss, ds, fmt.Errorf("Couldn't convert %v to float given: %v", s, err)
			}
			fs = append(fs, f)
		}
	}

	return fs, ss, ds, nil
}

func InferLineFormat(s string, sep rune, df string) string {
	s = string(regexp.MustCompile(string(sep)+"{2,}").ReplaceAll([]byte(s), []byte(string(sep))))
	ss := strings.Split(strings.TrimSpace(s), string(sep))
	lf := ""
	for _, sc := range ss {
		sc = strings.TrimSpace(sc)
		if _, err := strconv.ParseFloat(sc, 64); err == nil {
			lf += "f"
		} else if _, err := time.Parse(df, sc); err == nil && sc != "" {
			lf += "d"
		} else {
			lf += "s"
		}
	}
	return lf
}
