package main

import "fmt"

// formatMetric reimplements humanize.SI using division by 1000 + K/M suffix.
// The analysistest framework asserts an H005 diagnostic on the func line.
func formatMetric(num float64) string { // want "H005"
	if num >= 1_000_000 {
		return fmt.Sprintf("%.1fM", num/1_000_000)
	}

	if num >= 1_000 {
		return fmt.Sprintf("%.1fK", num/1_000)
	}

	return fmt.Sprintf("%.0f", num)
}

func main() {
	println(formatMetric(1234567))
}