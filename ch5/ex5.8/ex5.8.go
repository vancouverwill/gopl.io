// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 133.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 4 || len(os.Args)%3 != 1 {
		fmt.Println("ex5.8 should be run with at least one url argument `go build && ./outline2 http://www.gopl.io/ attName attrVal`")
		return
	}
	for i := 1; i < len(os.Args); i += 3 {
		startNode, _ := getDoc(os.Args[i], os.Args[i+1], os.Args[i+2])

		fmt.Println(ElementByAttr(startNode, htmlHasAttr, htmlHasAttr, os.Args[i+1], os.Args[i+2]))
	}
}

func getDoc(url, attName, attrVal string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	startNode, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return startNode, nil
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func ElementByAttr(n *html.Node, pre, post func(n *html.Node, attName, attrVal string) bool, attName, attrVal string) (result *html.Node) {
	if pre != nil {
		isObject := pre(n, attName, attrVal)
		if isObject {
			fmt.Println("success")
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = ElementByAttr(c, pre, post, attName, attrVal)

		if result != nil {
			return result
		}
	}

	if post != nil {
		isObject := pre(n, attName, attrVal)
		if isObject {
			fmt.Println("success 2")
			return n
		}
	}

	return result
}

//!-forEachNode

func htmlHasAttr(n *html.Node, attName, attrVal string) bool {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == attName {
				if a.Val == attrVal {
					return true
				}
			}
		}
	}
	return false
}
