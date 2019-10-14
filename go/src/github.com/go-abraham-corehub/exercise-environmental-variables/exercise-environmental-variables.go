//https://www.golangprograms.com/how-to-set-get-and-list-environment-variables.html

package main

import (
	"fmt"
	"os"
)

func main() {
	// Set custom env variable
	//os.Setenv("GOMAXPROCS", "100")
	os.Setenv("CUSTOM", "500")

	// fetcha all env variables
	for _, element := range os.Environ() {
		fmt.Println(element)
	}

	fmt.Println()
	// fetch specific env variables
	fmt.Println("GOMAXPROCS =>", os.Getenv("GOMAXPROCS"))
	fmt.Println("CUSTOM =>", os.Getenv("CUSTOM"))
	fmt.Println("GOROOT =>", os.Getenv("GOROOT"))
	fmt.Println("GOPATH =>", os.Getenv("GOPATH"))
	fmt.Println("GOHOME =>", os.Getenv("GOHOME"))
	fmt.Println("PATH =>", os.Getenv("PATH"))
}
