package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf(" parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return words, images, err
}

func countWordsAndImages(n *html.Node) (words, images int) {
	/* ... */
	words, images = 0, 0

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extraWords, extraImages := countWordsAndImages(c)
		words += extraWords
		images += extraImages
	}

	if n.Type == html.TextNode && len(strings.TrimSpace(n.Data)) > 0 {
		content := strings.TrimSpace(n.Data)
		contentWords := strings.Split(content, " ")
		words += len(contentWords)
	}

	if n.Data == "img" || n.Data == "script" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				images++
			}
		}
	}

	return words, images
}

func main() {
	for _, url := range os.Args[1:] {
		words, images, err := CountWordsAndImages(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CountWordsAndImages: %v\n", err)
			continue
		}
		fmt.Println(url, words, images)
	}
}
