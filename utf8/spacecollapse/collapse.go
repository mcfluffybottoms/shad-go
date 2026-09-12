//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	prevSpace := false
	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				builder.WriteRune(' ')
			}
		} else {
			builder.WriteRune(r)
		}
		prevSpace = unicode.IsSpace(r)
	}
	return builder.String()
}
