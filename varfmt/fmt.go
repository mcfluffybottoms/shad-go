//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var b strings.Builder
	var num strings.Builder
	var opened bool
	var argPos int
	for _, r := range format {
		if !opened && r == '{' {
			opened = true
			num.Reset()
		} else if opened && r == '}' {
			opened = false
			numStr := num.String()
			var i int
			if numStr == "" {
				i = argPos
			} else {
				ii, err := strconv.Atoi(numStr)
				if err != nil {
					panic(err)
				}
				i = ii
			}
			fmt.Fprint(&b, args[i])
			argPos++
		} else if opened {
			num.WriteRune(r)
		} else {
			b.WriteRune(r)
		}

	}
	return b.String()
}
