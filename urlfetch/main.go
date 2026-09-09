//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	urls := os.Args
	for i, url := range urls {
		if i == 0 {
			continue
		}
		resp, err := http.Get(url)
		check(err)
		defer resp.Body.Close()

		contents, err := io.ReadAll(resp.Body)
		check(err)

		fmt.Println(string(contents))
	}
}
