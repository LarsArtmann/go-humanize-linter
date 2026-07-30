package main

import "strings"

// formatNumber reimplements humanize.Comma using mod-3 grouping + WriteString
// comma insertion. The analysistest framework asserts an H002 diagnostic on
// the func line.
func formatNumber(n int) string { // want "H002"
	var b strings.Builder
	mod := 0
	for n > 0 {
		if mod != 0 && mod%3 == 0 {
			b.WriteString(",")
		}

		b.WriteRune(rune('0' + n%10))

		n /= 10
		mod++
	}

	return b.String()
}

func main() {
	println(formatNumber(1234567))
}
