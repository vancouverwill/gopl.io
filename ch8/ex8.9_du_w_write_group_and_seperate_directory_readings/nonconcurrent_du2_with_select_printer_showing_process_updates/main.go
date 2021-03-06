// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 247.

//!+main

// The du1 command computes the disk usage of the files in a directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Determine the initial directories.
	flag.Parse()
	// roots := flag.Args()
	// if len(roots) == 0 {
	// 	roots = []string{"."}
	// }
	root := flag.Arg(0)

	fmt.Println("root", root)

	// Traverse the file tree.

	type dirDetails struct {
		size    int64
		channel chan int64
		dirName string
	}
	// var allFileSizes []chan int64
	var allFileSizes []dirDetails
	// for _, directoryContent := range dirents(root) {
	// 	if directoryContent.IsDir() {
	// tmp := dirDetails{channel: make(chan int64), dirName: directoryContent.Name()}
	tmp := dirDetails{channel: make(chan int64), dirName: root}
	allFileSizes = append(allFileSizes, tmp)
	// allFileSizes := append(allFileSizes, make(chan int64))
	// 	}
	// }

	fmt.Println("allFileSizes", len(allFileSizes))

	// fileSizes := make(chan int64)
	go func(allFileSizes []dirDetails) {
		// for i, directoryContent := range dirents(root) {
		// 	if directoryContent.IsDir() {
		// for r := range allFileSizes {
		// subdir := filepath.Join(root, allFileSize.dirName)
		// subdir := filepath.Join(root, directoryContent.Name())
		// walkDir(subdir, allFileSize.channel)
		walkDir(allFileSizes[0].dirName, allFileSizes[0].channel)
		defer close(allFileSizes[0].channel)
		// }
		// 	}
		// }

	}(allFileSizes)

	var tick <-chan time.Time
	tick = time.Tick(100 * time.Millisecond)

	// Print the results.
	var totalNfiles, totalNbytes int64
	// for r := range allFileSizes {
	var nfiles, nbytes int64
	// for size := range allFileSize.channel {
	// 	nfiles++
	// 	nbytes += size
	// }
innerloop:
	for {
		select {
		case size, ok := <-allFileSizes[0].channel:
			if !ok {
				fmt.Println("break!")
				break innerloop // fileSizes was closed
				// break // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes, allFileSizes[0].dirName)
		}
	}
	printDiskUsage(nfiles, nbytes, allFileSizes[0].dirName)
	totalNfiles += nfiles
	totalNbytes += nbytes
	// }
	fmt.Println("total")
	printDiskUsage(totalNfiles, totalNbytes, root)
	fmt.Println()
	fmt.Println()
}

func printDiskUsage(nfiles, nbytes int64, dirName string) {
	fmt.Printf("%d files  %.1f MB location:%s\n", nfiles, float64(nbytes)/1e6, dirName)
}

//!-main

//!+walkDir

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

//!-walkDir

// The du1 variant uses two goroutines and
// prints the total after every file is found.
