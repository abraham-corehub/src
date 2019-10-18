package main

import (
	"fmt"
)

// Fetcher is custome
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// CrawlRecursive uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func CrawlRecursive(url string, depth int, fetcher Fetcher, quit chan bool, urlsVisited []string) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	if depth <= 0 {
		quit <- true
		return
	}

	if !isPresent(urlsVisited, url) {
		urlsVisited = unique(append(urlsVisited, url))
	} else {
		quit <- true
		return
	}

	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		fmt.Println(err)
		quit <- true
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	childrenQuit := make(chan bool)
	for _, childrenURL := range urls {
		go CrawlRecursive(childrenURL, depth-1, fetcher, childrenQuit, urlsVisited)
		<-childrenQuit

	}
	quit <- true
	return
}

func isPresent(urls []string, url string) bool {
	var present bool
	for _, u := range urls {
		if u == url {
			present = true
			break
		} else {
			present = false
		}
	}
	return present
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// Crawl Web Crawler
func Crawl(url string, depth int, fetcher Fetcher) {
	quit := make(chan bool)

	var urlsVisited []string

	go CrawlRecursive(url, depth, fetcher, quit, urlsVisited)

	<-quit

}

// unique https://www.golangprograms.com/remove-duplicate-values-from-slice.html
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
