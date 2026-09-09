//go:build !solution

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args
	count := make(map[string]int)

	var content strings.Builder
	for i := 1; i < len(args); i++ {
		path := args[i]
		f, err := os.Open(path)
		check(err)
		defer f.Close()
		for {
			bytes := make([]byte, 4)
			n, err := f.Read(bytes)
			if err == io.EOF {
				break
			}
			check(err)
			content.WriteString(string(bytes[:n]))
		}
		content.WriteString("\n")
	}

	lines := strings.SplitSeq(strings.Trim(content.String(), "\r\n"), "\n")
	for line := range lines {
		trimmed := strings.Trim(line, "\r\n")
		count[trimmed]++
	}

	for k, v := range count {
		if v > 1 {
			fmt.Printf("%d\t%s\n", v, k)
		}
	}
}
