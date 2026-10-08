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

// csvSampleRow formats an export row the way go-health-dashboard's
// ExportHandler does: strconv.FormatFloat with 'g' verb and -1 precision
// (shortest round-trip) beside WriteByte(',') CSV field delimiters. The
// negative precision parses as ast.UnaryExpr, which getBasicLit unwraps to
// "1" — that turned the old "-1" string exclusion into dead code and
// tripped H009 on every CSV writer. Neither signal is commaf: no fixed
// precision, and the commas are field delimiters, not digit grouping.
// Should NOT trigger H009.
func csvSampleRow(b *strings.Builder, at string, value float64, status string) {
	b.WriteString(at)
	b.WriteByte(',')
	b.WriteString(strconv.FormatFloat(value, 'g', -1, 64))
	b.WriteByte(',')
	b.WriteString(status)
}

func main() {
	_ = formatPlainFloat(3.14)
	_ = commaNoFloat(1000000)
	_ = formatPercent(1, 3)

	var b strings.Builder

	csvSampleRow(&b, "2026-10-08T12:00:00Z", 0.5, "warn")
}
