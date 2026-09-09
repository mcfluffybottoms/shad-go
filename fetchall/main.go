//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func logResponse(url string, resp *http.Response, e error, elapsed float64) string {
	if e != nil {
		return e.Error()
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("(%s) %fs \t%d: %s\n", resp.Status, elapsed, len(body), url)
	}
}

func main() {
	urls := os.Args[1:]
	outChannel := make(chan string, len(urls))
	for _, url := range urls {
		go func(url string) {
			t := time.Now()
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
			}
			elapsed := time.Since(t).Seconds()
			message := logResponse(url, resp, err, elapsed)
			outChannel <- message
		}(url)

	}
	for range urls {
		fmt.Print(<-outChannel)
	}
}
