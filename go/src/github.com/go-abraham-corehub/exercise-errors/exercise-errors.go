package main

import (
	"fmt"
)

// ErrNegativeSqrt is a custom type for handling Negative value errors for Sqrt function
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative numbers: %v", float64(e))
}

// Sqrt finds the square root of a number using Newton's Method
func Sqrt(x float64) (float64, error) {

	if x < 0 {
		return x, ErrNegativeSqrt(x)
	}
	z := 0.5 * x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return 0, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
