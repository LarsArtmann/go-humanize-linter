package main

import (
	"fmt"
)

// Pattern A with a suppression directive: this MUST NOT be flagged.
//
//nolint:gohumanize // intentional self-rolled example, suppressed
func formatBytes(bytes int64) string {
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
	println(formatBytes(1073741824))
}
