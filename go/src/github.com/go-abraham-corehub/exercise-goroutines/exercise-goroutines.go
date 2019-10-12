package main

import (
	"fmt"
	"time"
)

func rndint(x int) {
	for i := 0; i < 5; i++ {
		time.Sleep(1000 * time.Millisecond)
		fmt.Print(x)
	}
}

func main() {
	go rndint(0)
	rndint(1)
}
