package main

import (
	"fmt"
)

//SqrtA calculates the Square Root of a number using the Newton's Algorithm
func SqrtA(x float64) float64 {
	z := 0.5 * x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
		fmt.Println(i, z)
	}
	return z
}

//SqrtB calculates the Square Root of a number using the Newton's Algorithm
func SqrtB(x float64) float64 {
	z := 0.5 * x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
		fmt.Println(i, z)
	}
	return z
}

func main() {
	fmt.Println(SqrtA(32423424))
}
