package main

import "fmt"

type random struct {
	size  int
	count int
}

func main() {
	fmt.Println("hello world")

	localArray()
	arrayAsParameter()
}

func localArray() {
	var randoms []random

	randoms = append(randoms, random{})

	for i := 0; i < 10; i++ {
		randoms[0].size += i
		randoms[0].count++
	}

	for _, r := range randoms {
		fmt.Printf("localArray: %v stores %v\n", r.count, r.size)
	}
}

func arrayAsParameter() {
	var randoms []random

	randoms = append(randoms, random{})

	addToArray(randoms)

	for _, r := range randoms {
		fmt.Printf("arrayAsParameter: %v stores %v\n", r.count, r.size)
	}
}

func addToArray(randoms []random) {
	for r := range randoms {
		for i := 0; i < 10; i++ {
			randoms[r].size += i
			randoms[r].count++
		}
	}
}
