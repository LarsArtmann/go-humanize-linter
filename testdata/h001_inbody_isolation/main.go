package main

import (
	"fmt"
)

// Function A: has H001 pattern + in-body //nolint. Must be suppressed.
func formatBytesA(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp]) //nolint:gohumanize
}

// Function B: has H001 pattern with NO directive. Must still be flagged.
func formatBytesB(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	println(formatBytesA(1073741824))
	println(formatBytesB(1073741824))
}
