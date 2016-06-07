// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	"gopl.io/ch9/bank1"
)

func TestBankDeposit(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		fmt.Println("=", bank.Balance())
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		fmt.Println("=", bank.Balance())
		bank.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	// Alice
	go func() {
		bank.Deposit(300)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.WithDraw(100)
		done <- struct{}{}
	}()

	go func() {
		bank.WithDraw(75)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done
	<-done

	fmt.Println("=", bank.Balance())

	if got, want := bank.Balance(), 425; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

func TestBankWithDrawl(t *testing.T) {
	// done := make(chan struct{})

}
