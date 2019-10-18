package main

import (
	"golang.org/x/tour/reader"
)

// MyReader is a custom type to create a custom Reader
type MyReader struct{}

// Read implements the Read method of io.Reader interface
// supplies (streams) the ASCII value of 'A' indefinitely
func (m MyReader) Read(b []byte) (int, error) {

	for i := range b {
		b[i] = 'A'
	}

	return len(b), nil
}

func main() {
	reader.Validate(MyReader{})
}
