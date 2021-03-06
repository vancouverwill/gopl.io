// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 250.

// The du3 command computes the disk usage of the files in a directory.
package main

// PLAN
// 1. make array of objects each holding root dir, channel, waitgroup, their size and their

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
	count   int
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
	roots := flag.Args()
	// if len(roots) == 0 {
	// 	roots = []string{"."}
	// }

	// root := flag.Arg(0)

	// fmt.Println(root)

	//!+
	// Traverse each root of the file tree in parallel.
	// fileSizes := make(chan int64)
	// var n sync.WaitGroup
	var allFileSizes []dirDetails
	for _, root := range roots {
		tmp := dirDetails{c: make(chan int64), dirName: root}
		allFileSizes = append(allFileSizes, tmp)
	}

	for r := range roots {
		allFileSizes[r].n.Add(1)
		go walkDir(allFileSizes[r].dirName, &allFileSizes[r].n, allFileSizes[r].c)
		// }
		go func(r int) {
			allFileSizes[r].n.Wait()
			close(allFileSizes[r].c)
		}(r)
	}

	// combinedFileSizes := combine(allFileSizes)
	// combine(allFileSizes)

	var tick <-chan time.Time

	done := make(chan bool)

	var numDirectoriesWorking sync.WaitGroup

	numDirectoriesWorking.Add(len(allFileSizes))

	var nfiles, nbytes int64

	for r := range roots {
		go func(input *dirDetails) {
			defer func() {
				// done <- true
				numDirectoriesWorking.Done()
			}()
			for inputI := range input.c {
				input.size += inputI
				input.count++

				nfiles++
				nbytes += inputI
				// combined <- inputI
			}
		}(&allFileSizes[r])
	}

	go func() {
		numDirectoriesWorking.Wait()
		done <- true
	}()

	// Print the results periodically.

	if *vFlag {
		tick = time.Tick(1 * time.Millisecond)
	}
	// var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done:
			break loop // fileSizes was closed

		case <-tick:
			printDiskUsageCombined(allFileSizes)
		}
	}

	println()
	printDiskUsage(nfiles, nbytes) // final totals
	printDiskUsageCombined(allFileSizes)

	end := time.Now()
	diff := end.Sub(begin)
	fmt.Println("script took", diff.String())
}

//!-

func printDiskUsageCombined(dataSet []dirDetails) {
	fmt.Printf("Combined ")
	for _, data := range dataSet {
		fmt.Printf("%s has %d files  %.1f MB", data.dirName, data.count, float64(data.size)/1e6)
	}
	fmt.Printf("\n ")
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("printDiskUsage TOTAL - %d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
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
	go func(input *dirDetails) {
		defer close(combined)
		for inputI := range input.c {
			input.size += inputI
			input.count++
			combined <- inputI
		}
	}(&input[0])
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
