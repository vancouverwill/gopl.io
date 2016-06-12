// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// exercise : following the approach of mirrored query in 8.4.4 implement a variant of fetch that requests several URLS concurrently.
// As soon as the first response arrives cancel the other ones

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	// resp, err := http.Get(url)
	req, _ := http.NewRequest("GET", url, nil)
	tr := &http.Transport{} // TODO: copy defaults from http.DefaultTransport
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	go func() {
		select {
		case <-done:
			tr.CancelRequest(req)
		}
	}()
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	local = strings.Replace(url, "http://", "", 2) + "." + local
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	close(done)
	return local, n, err
}

//!-

var done = make(chan struct{})

func main() {
	for _, url := range os.Args[1:] {
		go func(url string) {
			local, n, err := fetch(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
				return
			}
			fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
		}(url)
	}

	<-done
}
