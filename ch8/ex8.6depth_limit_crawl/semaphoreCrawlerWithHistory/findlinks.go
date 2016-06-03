// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 6)

func crawl(url string, depth int, path []string) []string {
	fmt.Println(depth, url, path)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

//!+
func main() {
	type depthList struct {
		depth int
		list  []string
		path  []string
	}
	worklist := make(chan depthList)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() {
		dl := depthList{0, os.Args[1:], nil}
		worklist <- dl
	}()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		dlist := <-worklist
		for _, link := range dlist.list {
			if !seen[link] {
				seen[link] = true

				if dlist.depth > 2 {
					continue
				}
				n++
				go func(link string, depth int, previousList depthList) {
					dl := depthList{depth + 1, crawl(link, depth, previousList.path), append(previousList.path, link)}
					worklist <- dl
				}(link, dlist.depth, dlist)
			}
		}
	}
}

//!-
