package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t != nil {
		Walk(t.Left, ch)
		ch <- t.Value
		Walk(t.Right, ch)
	} else {
		return
	}
}

// Walker https://medium.com/@cooldeep25/solution-to-a-tour-of-go-exercise-equivalent-binary-trees-d1fff8d3cb6f
func Walker(t *tree.Tree, ch chan int) {
	Walk(t, ch)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walker(t1, ch1)
	go Walker(t2, ch2)

	for i := range ch1 {
		if i != <-ch2 {
			return false
		}
	}
	return true
}

func main() {
	for range []int{1, 2, 3} {
		fmt.Println(tree.New(1))
	}
	fmt.Println(tree.New(1))
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
