package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	testReadFileBuf()
}

func testReadFileBuf() {
	nameFileInput := "The Ugly.html"
	f, err := os.Open(nameFileInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	mFB := myFileBuf{bufio.NewReader(f), 0, nil}
	b, err := read(mFB)
	for err == nil {
		fmt.Print(string(b))
		b, err = read(mFB)
	}
}

func testReadStringBuf() {
	str := "The Ugly.html"

	mSB := myStringBuf{strings.NewReader(str), 0, nil}
	b, err := read(mSB)
	for err == nil {
		fmt.Print(string(b))
		b, err = read(mSB)
	}
}

type myBufInterface interface {
	Read() (rune, error)
	Write() error
}

type myFileBuf struct {
	myReader *bufio.Reader
	myRune   rune
	myErr    error
}

type myStringBuf struct {
	myReader *strings.Reader
	myRune   rune
	myErr    error
}

func (s myStringBuf) Read() (rune, error) {
	b := make([]byte, 1)
	_, err := s.myReader.Read(b)
	s.myRune = rune(b[0])
	s.myErr = err
	return s.myRune, s.myErr
}

func (f myFileBuf) Read() (rune, error) {
	f.myRune, _, f.myErr = f.myReader.ReadRune()
	return f.myRune, f.myErr
}

func (f myFileBuf) Write() error {

	return nil
}

func (s myStringBuf) Write() error {

	return nil
}

func read(mRI myBufInterface) (rune, error) {

	b, err := mRI.Read()
	return b, err
}

func write(mRI myBufInterface) (rune, error) {

	err := mRI.Write()
	return 0, err
}
