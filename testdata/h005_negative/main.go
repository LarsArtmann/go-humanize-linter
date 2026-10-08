package main

import "strconv"

// This function divides by 1000 but formats milliseconds, not SI prefixes.
func formatMillis(ms float64) string {
	return strconv.FormatFloat(ms/1000, 'f', 2, 64) + "s"
}

// Regression: a benchmark table row divides by 1e3 for the millisecond
// latency columns, and the ALL-CAPS header ends in 'T' ("THROUGHPUT").
// Prose must never satisfy the embedded-suffix check — only format strings
// carrying a % verb (like "%.1fM") may.
func formatBenchmarkRow(label string, p50us float64) string {
	header := "PLAN P95 THROUGHPUT"

	return header + " " + label + " " + strconv.FormatFloat(p50us/1e3, 'f', 2, 64) + "ms"
}

func main() {
	println(formatMillis(1500))
	println(formatBenchmarkRow("sqlite", 1234))
}
