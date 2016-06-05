// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

//!+
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to pass port number as first argument to specify which ports to watch")
	}
	// for i := 1; i < len(os.Args); i++ {
	// go func(i int) {
	port := os.Args[1]
	readInput(port)
	// }(i)
	// }

	time.Sleep(30 * time.Second)
}

func readInput(port string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
