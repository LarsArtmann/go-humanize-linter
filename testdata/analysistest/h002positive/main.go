package main

import "strings"

// formatNumber reimplements humanize.Comma using mod-3 grouping + WriteString
// comma insertion. The analysistest framework asserts an H002 diagnostic on
// the func line.
func formatNumber(n int) string { // want "H002"
	str := ""
	mod := 0
	for n > 0 {
		if mod != 0 && mod%3 == 0 {
			str = "," + str
		}

		str = string(rune('0'+n%10)) + str

		n /= 10
		mod++
	}

	var b strings.Builder
	b.WriteString(str)

	return b.String()
}

func main() {
	println(formatNumber(1234567))
}