// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 122.
//!+main

// Findlinks1 prints the links in an HTML document read from standard input.
package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	fmt.Printf("%v", doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	index := 0
	elementCount := make(map[string]int)
	visit(elementCount, doc)
	// for elementType, count := range visit(elementCount, doc) {
	// fmt.Println(index, elementType, count)
	index++
	// }
}

//!-main

//!+visit
// visit appends to links each link found in n and returns the result.
func visit(elementCount map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return elementCount
	}

	if n.Type == html.ElementNode && n.Data == "script" {
		return elementCount
	}
	if n.Type == html.ElementNode && n.Data == "style" {
		return elementCount
	}

	if n.Type == html.TextNode && len(strings.TrimSpace(n.Data)) > 0 {
		fmt.Printf("%d %s\n", len(strings.TrimSpace(n.Data)), strings.TrimSpace(n.Data))
		elementCount[n.Data]++
	}

	elementCount = visit(elementCount, n.FirstChild)
	elementCount = visit(elementCount, n.NextSibling)

	return elementCount
}

//!-visit

/*
//!+html
package html

type Node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func Parse(r io.Reader) (*Node, error)
//!-html
*/
