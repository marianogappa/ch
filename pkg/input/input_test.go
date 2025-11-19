package input

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileInput(t *testing.T) {
	// Create temp file
	content := "line1\nline2\n"
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test
	fi := NewFileInput(tmpfile.Name())
	stream, err := fi.Stream()
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	var lines []string
	for b := range stream {
		lines = append(lines, string(b))
	}

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "line1" {
		t.Errorf("Expected line1, got %s", lines[0])
	}
}

func TestStdinInput(t *testing.T) {
	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Replace os.Stdin with the read end of the pipe
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	// Write to the write end of the pipe in a goroutine
	go func() {
		defer w.Close()
		w.Write([]byte("stdin_line1\nstdin_line2\n"))
	}()

	// Test
	si := NewStdinInput()
	stream, err := si.Stream()
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	var lines []string
	for b := range stream {
		lines = append(lines, string(b))
	}

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "stdin_line1" {
		t.Errorf("Expected stdin_line1, got %s", lines[0])
	}
}
