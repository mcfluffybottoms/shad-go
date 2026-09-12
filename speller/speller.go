//go:build !solution

package speller

import "strings"

var nums = map[int64]string{
	0:  "",
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
	20: "twenty",
	30: "thirty",  // 30
	40: "forty",   // 40
	50: "fifty",   // 50
	60: "sixty",   // 60
	70: "seventy", // 70
	80: "eighty",  // 80
	90: "ninety",  // 90
}

func Spell(m int64) string {
	if m == 0 {
		return "zero"
	}

	var b strings.Builder
	var n int64
	if m < 0 {
		n = -m
		b.WriteString("minus ")
	} else {
		n = m
	}

	if n <= 20 {
		b.WriteString(nums[n])
		return b.String()
	}
	var counter = map[string]int64{
		"":          n % 1000,
		" thousand": n % 1000000 / 1000,
		" million":  n % 1000000000 / 1000000,
		" billion":  n / 1000000000,
	}

	var countUsed = false
	for _, key := range []string{" billion", " million", " thousand", ""} {
		count := counter[key]
		if count > 0 {
			if countUsed {
				b.WriteString(" ")
			}
			b.WriteString(spellHundred(count))
			b.WriteString(key)
			countUsed = true
		}
	}

	return b.String()
}

func spellHundred(n int64) string {
	var b strings.Builder
	b.WriteString(nums[n/100])
	if n/100 != 0 {
		b.WriteString(" hundred")
	}
	if n/100 != 0 && (n%100-n%10 != 0 || n%10 != 0) {
		b.WriteString(" ")
	}

	if n%100 < 20 {
		b.WriteString(nums[n%100])
	} else {
		b.WriteString(nums[n%100-n%10])
		if n%100-n%10 != 0 && n%10 != 0 {
			b.WriteString("-")
		}
		b.WriteString(nums[n%10])
	}
	return b.String()
}
