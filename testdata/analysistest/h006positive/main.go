package main

import (
	"fmt"
	"strings"
)

// trimTrailing reimplements humanize.Ftoa using nested strings.TrimRight.
// The analysistest framework asserts an H006 diagnostic on the func line.
func trimTrailing(v float64) string { // want "H006"
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", v), "0"), ".")
}

func main() {
	println(trimTrailing(3.14))
}