package main

import (
	"fmt"
	"strings"
)

// formatAmount reimplements humanize.Commaf using fmt.Sprintf("%.2f") + a
// manual separator loop. The analysistest framework asserts an H009
// diagnostic on the func line.
func formatAmount(f float64) string { // want "H009"
	s := fmt.Sprintf("%.2f", f)
	var b strings.Builder
	dot := strings.Index(s, ".")
	intPart := s[:dot]
	fracPart := s[dot:]

	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteString(",")
		}

		b.WriteRune(r)
	}

	b.WriteString(fracPart)

	return b.String()
}

func main() {
	println(formatAmount(1234567.89))
}