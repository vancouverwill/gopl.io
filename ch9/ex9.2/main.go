// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.

/*
The population count of a bitstring is the number of set bits (1-bits) in the string.
For instance, the population count of the number 23, which is represented in binary as 101112, is 4.
The population count is used in cryptography and error-correcting codes, among other topics in computer science;
some people use it as an interview question. The population count is also known as Hamming weight.
*/

// (Package doc comment intentionally malformed to demonstrate golint.)
//!+
package popcount

import "sync"

// pc[i] is the population count of i.
var pc [256]byte

var loadPopCount sync.Once

// func init() {
func lazyLoad() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	loadPopCount.Do(lazyLoad)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

//!-
