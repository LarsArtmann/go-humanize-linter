package main

import (
	"fmt"
)

// Pattern with a scoped H002 directive: this MUST NOT suppress H001.
// The directive suppresses only H002, so H001 should still fire.
//
//nolint:gohumanize:H002
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
