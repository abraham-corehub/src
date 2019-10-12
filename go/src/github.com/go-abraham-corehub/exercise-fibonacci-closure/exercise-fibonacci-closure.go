package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		y, c := a, a+b
		a, b = b, c
		return y
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		k := f()
		fmt.Print(k, ", ")
	}
}
