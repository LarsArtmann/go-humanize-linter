package main

import (
	"fmt"
	"strconv"
	"strings"
)

// formatPlainFloat uses fmt.Sprintf("%.2f") but has no comma separator loop.
// Should NOT trigger H009.
func formatPlainFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// commaNoFloat has a comma separator loop but no %.Nf formatting. This is
// integer comma formatting (H002 territory), not float comma (H009).
// Should NOT trigger H009.
func commaNoFloat(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteRune(r)
	}
	return b.String()
}

// formatPercent uses %.1f but returns a percentage string, not a comma-formatted
// float. No separator loop. Should NOT trigger H009.
func formatPercent(part, total float64) string {
	return fmt.Sprintf("%.1f%%", part/total*100)
}

func main() {
	_ = formatPlainFloat(3.14)
	_ = commaNoFloat(1000000)
	_ = formatPercent(1, 3)
}
