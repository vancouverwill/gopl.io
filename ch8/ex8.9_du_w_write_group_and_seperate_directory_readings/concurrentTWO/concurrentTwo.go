// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 250.

// The du3 command computes the disk usage of the files in a directory.
package main

// The du3 variant traverses all directories in parallel.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type dirDetails struct {
	size    int64
	c       chan int64
	n       sync.WaitGroup
	dirName string
}

var vFlag = flag.Bool("v", false, "show verbose progress messages")

//!+
func main() {
	// ...determine roots...

	begin := time.Now()

	//!-
	flag.Parse()

	// Determine the initial directories.
	// roots := flag.Args()
	// if len(roots) == 0 {
	// 	roots = []string{"."}
	// }

	root := flag.Arg(0)

	fmt.Println(root)

	//!+
	// Traverse each root of the file tree in parallel.
	// fileSizes := make(chan int64)
	// var n sync.WaitGroup
	var allFileSizes []dirDetails
	tmp := dirDetails{c: make(chan int64), dirName: root}
	allFileSizes = append(allFileSizes, tmp)

	// for _, root := range roots {
	allFileSizes[0].n.Add(1)
	go walkDir(allFileSizes[0].dirName, &allFileSizes[0].n, allFileSizes[0].c)
	// }
	go func() {
		allFileSizes[0].n.Wait()
		close(allFileSizes[0].c)
	}()
	//!-

	combinedFileSizes := combine(allFileSizes)

	// Print the results periodically.
	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-combinedFileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes) // final totals
	//!+
	// ...select loop...

	end := time.Now()
	diff := end.Sub(begin)
	fmt.Println("script took", diff.String())
}

//!-

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+walkDir
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func combine(input []dirDetails) chan int64 {
	combined := make(chan int64)
	go func(input chan int64) {
		defer close(combined)
		for inputI := range input {
			combined <- inputI
		}
	}(input[0].c)
	return combined
}

//!-walkDir

//!+sema
// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	// ...
	//!-sema

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
