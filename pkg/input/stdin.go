package input

import (
	"bufio"
	"io"
	"os"
)

type StdinInput struct {
	reader io.Reader
}

func NewStdinInput() *StdinInput {
	return &StdinInput{reader: os.Stdin}
}

func NewReaderInput(r io.Reader) *StdinInput {
	return &StdinInput{reader: r}
}

func (s *StdinInput) Stream() (<-chan []byte, error) {
	out := make(chan []byte)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			// We need to copy the bytes because scanner.Bytes() is reused
			b := scanner.Bytes()
			c := make([]byte, len(b))
			copy(c, b)
			out <- c
		}
	}()
	return out, nil
}

// Ensure StdinInput satisfies the interface (implicitly, but good for checking)
// var _ ch.Input = (*StdinInput)(nil)
// Note: avoiding circular dependency or just not importing ch here if not strictly needed for the struct definition,
// but the method signature matches.
