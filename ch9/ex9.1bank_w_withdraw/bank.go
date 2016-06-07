// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int) // send amount to deposit
type withDrawl struct {
	amount int
	result chan bool
}

var withdrawls = make(chan withDrawl) // send amount to deposit
var balances = make(chan int)         // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func WithDraw(amount int) bool {
	r := make(chan bool)
	wd := withDrawl{amount, r}
	withdrawls <- wd
	return <-r
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case wD := <-withdrawls:
			if balance >= wD.amount {
				balance -= wD.amount
				wD.result <- true
			} else {
				wD.result <- false
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
