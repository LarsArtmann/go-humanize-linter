package main

import (
	"strconv"
	"strings"
)

// Pattern: i += 3 stepping + strings.Join with comma
func formatNumber(n int64) string {
	digits := strconv.FormatInt(n, 10)
	if len(digits) <= 3 {
		return digits
	}
	var groups []string
	pre := len(digits) % 3
	if pre > 0 {
		groups = append(groups, digits[:pre])
	}
	for i := pre; i < len(digits); i += 3 {
		groups = append(groups, digits[i:i+3])
	}
	return strings.Join(groups, ",")
}

func main() {
	println(formatNumber(1234567))
}
