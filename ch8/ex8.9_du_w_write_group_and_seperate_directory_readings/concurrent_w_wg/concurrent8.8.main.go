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

var vFlag = flag.Bool("v", false, "show verbose progress messages")

//!+
func main() {
	// ...determine roots...

	//!-
	flag.Parse()

	// Determine the initial directories.
	// roots := flag.Args()
	// if len(roots) == 0 {
	// 	roots = []string{"."}
	// }

	root := flag.Arg(0)

	if root == "" {
		root = "~/Google_Drive/books/"
	}

	fmt.Println("root dir is ", root)

	type dirDetails struct {
		size    int64
		c       chan int64
		n       sync.WaitGroup
		nSize   int
		dirName string
	}
	// var allFileSizes []chan int64
	var allFileSizes []dirDetails
	for _, directoryContent := range dirents(root) {
		if directoryContent.IsDir() {
			tmp := dirDetails{c: make(chan int64), dirName: directoryContent.Name()}
			allFileSizes = append(allFileSizes, tmp)
			// allFileSizes := append(allFileSizes, make(chan int64))
		}
	}

	fmt.Println("allFileSizes length ", len(allFileSizes))
	fmt.Println("allFileSizes ", allFileSizes)

	//!+
	// Traverse each root of the file tree in parallel.
	// fileSizes := make(chan int64)
	// var n sync.WaitGroup
	// for _, root := range roots {
	for _, allFileSize := range allFileSizes {
		allFileSize.n.Add(1)
		allFileSize.nSize++
		subdir := filepath.Join(root, allFileSize.dirName)
		go walkDir(subdir, &allFileSize.n, allFileSize.c, allFileSize.nSize)
		go func(allFileSize dirDetails) {
			fmt.Println("waiting for waiting group")
			allFileSize.n.Wait()
			fmt.Println("closing ", subdir)
			close(allFileSize.c)
		}(allFileSize)
	}

	// go func() {
	// 	n.Wait()
	// 	close(fileSizes)
	// }()

	// go func() {
	// 	allFileSizes[0].n.Wait()
	// 	close(allFileSizes[0].c)
	// }()

	// for _, allFileSize := range allFileSizes {
	// 	go func(allFileSize dirDetails) {
	// 		allFileSize.n.Wait()
	// 		close(allFileSize.c)
	// 	}(allFileSize)
	// }
	//!-

	// / Print the results periodically.
	var tick <-chan time.Time
	// if *vFlag {
	tick = time.Tick(500 * time.Millisecond)

	// var chan(bool )
	collectedChan := make(chan int64)

	go func(allFileSizes []dirDetails) {
	joiningLoop:
		// for {
		// 	select {
		// case size, ok := <-allFileSizes[0].c
		for {
			select {
			case size, ok := <-allFileSizes[0].c:
				if !ok {
					fmt.Println("allsizes finished")
					break joiningLoop
				}
				fmt.Println("first channel:", allFileSizes[0].nSize)
				// allFileSizes[0].n.
				collectedChan <- size
			default:
				// fmt.Println("no activity")
			}

		}
	}(allFileSizes)

	////////////////// ************ NOT GETTING TO HERE CURRENTLY AS WAIT GROUP NOT ENDING

	// }
	var totalNfiles, totalNbytes int64
	// for _, allFileSize := range allFileSizes {
	var nfiles, nbytes int64
printingLoop:
	for {
		select {
		case size, ok := <-collectedChan:
			if !ok {
				fmt.Printf("dir finished:%s\n", allFileSizes[0].dirName)
				printDiskUsage(nfiles, nbytes)
				break printingLoop // fileSizes was closed
			}
			fmt.Println("second channel")
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("totals")
	printDiskUsage(totalNfiles, totalNbytes) // final totals
	//!+
	// ...select loop...
}

//!-

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+walkDir
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64, waitGroupSize int) {
	defer n.Done()
	defer func(waitGroupSize int) { waitGroupSize = waitGroupSize - 1 }(waitGroupSize)
	fmt.Println("walkDir()", dir)
	// fmt.Println("dirents()", dirents(dir))
	for _, entry := range dirents(dir) {

		if entry.IsDir() {
			n.Add(1)
			waitGroupSize++
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes, waitGroupSize)
		} else {
			fileSizes <- entry.Size()
		}
	}
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

	fmt.Println("func dirent()", dir)
	// os.Exit(0)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
