package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

// WordCount counts the occurence of each words in 's' and
// returns a map containing the words as key and corrensponding count as element
func WordCount(s string) map[string]int {

	m := make(map[string]int)

	wordsArray := strings.Split(s, " ")

	for _, v := range wordsArray {
		_, ok := m[v]
		if !ok {
			m[v] = 1
		} else {
			m[v]++
		}
	}

	return m
}

func main() {

	wc.Test(WordCount)
}
