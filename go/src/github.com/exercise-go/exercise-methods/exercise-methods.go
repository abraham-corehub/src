package main

import (
	"fmt"
	"math"
)

// Vertex contains two float values
type Vertex struct {
	X, Y float64
}

// Abs returns the absolute value of the components
func (v Vertex) Abs() float64 {
	return math.Sqrt((v.X*v.X + v.Y*v.Y))
}

func main() {

	v := Vertex{3, 4}
	fmt.Println(v.Abs())
}
