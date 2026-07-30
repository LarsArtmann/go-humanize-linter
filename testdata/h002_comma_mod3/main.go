package main

import (
	"strconv"
	"strings"
)

// Pattern: (len-i)%3 == 0 + strings.Builder comma insertion
func formatNumber(n int) string {
	str := strconv.Itoa(n)
	if len(str) <= 3 {
		return str
	}
	var result strings.Builder
	for i, r := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(r)
	}
	return result.String()
}

func main() {
	println(formatNumber(1234567))
}
