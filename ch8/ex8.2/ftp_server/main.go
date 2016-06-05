// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 222.

// Clock is a TCP server that periodically writes the time.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func handleConn(c net.Conn, port string) {
	defer c.Close()
	input := bufio.NewScanner(c)

	_, err := io.WriteString(c, "successfully connected to "+port+"\nftp>")
	if err != nil {
		return // e.g., client disconnected
	}

	for input.Scan() {
		inputText, err := processInput(input.Text(), port)
		if err != nil {
			return // e.g., client disconnected
		}
		_, err = io.WriteString(c, inputText+"\nftp>")
		if err != nil {
			return // e.g., client disconnected
		}
		// time.Sleep(1 * time.Second)
	}
}

func processInput(input string, port string) (string, error) {
	switch input {
	case "ls":
		files, err := ioutil.ReadDir(".")
		if err != nil {
			err = fmt.Errorf(" error reading directoy: %s", err)
			return "", err
		}
		output := ""
		for _, file := range files {
			output += file.Name() + " "
		}
		return output, nil
	case "put":
		return "put file in directory", nil
	default:
		return input + " port:" + port, nil
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
