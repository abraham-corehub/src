package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r1 rot13Reader) Read(b []byte) (int, error) {

	a := make([]byte, len(b))
	n, err := r1.r.Read(a)

	for i := range b {

		b[i] = rot13(a[i])
	}
	return n, err
}

func rot13(a byte) byte {

	var b byte

	if a >= 'A' && a <= 'z' {
		b = a + 13
		if b > 'z' || (b > 'Z' && b < 'a') || (b <= 'm' && b >= 'a') {
			b -= 26
		}
	} else {
		b = a
	}
	return b
}

func main() {
	s := strings.NewReader("Lbh Penpxrq gur Pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
