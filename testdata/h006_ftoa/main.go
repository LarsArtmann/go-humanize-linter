package main

import (
	"fmt"
	"strings"
)

// Pattern: nested strings.TrimRight — trailing zero stripping
func formatBucketVal(v float64) string {
	if v == 0 {
		return "0"
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", v), "0"), ".")
}

func main() {
	println(formatBucketVal(3.14))
}
