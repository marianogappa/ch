package input

import (
	"bufio"
	"os"
)

type FileInput struct {
	path string
}

func NewFileInput(path string) *FileInput {
	return &FileInput{path: path}
}

func (f *FileInput) Stream() (<-chan []byte, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		defer close(out)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			b := scanner.Bytes()
			c := make([]byte, len(b))
			copy(c, b)
			out <- c
		}
	}()
	return out, nil
}
