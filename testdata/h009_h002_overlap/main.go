package main

import (
	"fmt"
	"strconv"
	"strings"
)

// formatPrice reimplements humanize.Commaf: it formats a float with %.2f
// then manually inserts commas. This should trigger H009 ONLY, not H002.
func formatPrice(price float64) string {
	s := fmt.Sprintf("%.2f", price)
	parts := strings.Split(s, ".")
	intPart := parts[0]

	var b strings.Builder
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteByte(intPart[i])
	}

	return b.String() + "." + parts[1]
}

// formatPriceInt reimplements humanize.Comma: integer-only comma formatting.
// This should trigger H002 ONLY, not H009 (no float formatting).
func formatPriceInt(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
