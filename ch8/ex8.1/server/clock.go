// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 222.

// Clock is a TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn, port string) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05")+" port:"+port+"\n")
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to pass port number as first argument to specify which ports to serve on")
	}
	port := os.Args[1]
	listener, err := net.Listen("tcp", "localhost:"+port)

	if err != nil {
		log.Fatal(err)
	}
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, port) // handle connections concurrently
	}
	//!-
}
