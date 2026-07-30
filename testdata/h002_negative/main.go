package main

import "strings"

// formatSSN formats a US Social Security Number for display. It writes
// separators via WriteString but does NOT have a mod-3 / step-by-3 grouping
// loop. Should NOT trigger H002.
func formatSSN(s string) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(r)
		if b.Len() == 3 || b.Len() == 6 {
			b.WriteString("-")
		}
	}

	return b.String()
}

func main() {
	println(formatSSN("123456789"))
}