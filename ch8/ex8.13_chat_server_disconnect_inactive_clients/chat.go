// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

//!+broadcaster
type client struct {
	c    chan<- string // an outgoing message channel
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[string]client) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for _, cli := range clients {
				cli.c <- msg
			}

		case cli := <-entering:
			clients[cli.name] = cli
			messageToNewClient := "Already here: "
			for _, cli := range clients {
				messageToNewClient += cli.name + ", "
			}
			cli.c <- messageToNewClient

		case cli := <-leaving:
			if _, ok := clients[cli.name]; ok == true {
				delete(clients, cli.name)
				close(cli.c)
			}
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{ch, who}

	reset := make(chan struct{})

	go func() {
		tick := time.Tick(1 * time.Second)
		for countdown := 5; countdown > 0; countdown-- {
			// fmt.Println(countdown)
			select {
			case <-tick:
				// Do nothing.
				// case <-abort:
				// 	fmt.Println("Launch aborted!")
				// 	return
				// }
				ch <- fmt.Sprintf("%d remaining", countdown)
			case <-reset:
				countdown = 10
				ch <- "you have reset!"
			}
		}

		messages <- who + " run out of time and has left"
		leaving <- client{ch, who}
		conn.Close()
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		reset <- struct{}{}
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client{ch, who}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
