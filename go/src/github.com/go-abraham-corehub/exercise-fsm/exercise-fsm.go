package main

import (
	"fmt"
)

func main() {
	str := "0101011110010101011101101001011100101"
	parse(str)
}

func parse(str string) {
	fmt.Println(detectPatternSequence(str))
}

func detectPatternSequence(strm string) (string, error) {

	cS := "0"
	out := make([]byte, 0)
	var err error

	sT := []StateTable{
		{"1", byte('0'), "0"},
		{"11", byte('0'), "0"},
		{"0", byte('1'), "1"},
		{"1", byte('1'), "11"},
	}

	for _, cH := range strm {
		for _, rowST := range sT {
			if rowST.in == byte(cH) && cS == rowST.cS {
				cS = rowST.nS
				break
			}
		}
		if cS == "11" {
			out = append(out, byte('1'))
		} else {
			out = append(out, byte('0'))
		}
	}
	//fmt.Println("")

	return string(out), err
}

// StateTable is the State Transition Table of the FSM
type StateTable struct {
	cS string
	in byte
	nS string
}
