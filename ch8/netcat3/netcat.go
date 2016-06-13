// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

/*
`go build ../ch8/reverb2
go build ../ch8/netcat2
./reverb2 & ./netcat2
*/

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		src, err := io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		if err != nil {
			log.Fatalf("connection is closed:%v", err)
		}
		log.Println("src", src)
		log.Println("err", err)
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	log.Println("going to closed")
	mustCopy(conn, os.Stdin)
	conn.Close()
	log.Println("has closed")
	<-done // wait for background goroutine to finish
	log.Println("has done")
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("mustCopy %v", err)
	}
}
