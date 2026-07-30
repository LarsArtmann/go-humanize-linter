package main

import (
	"fmt"
)

// Pattern C: []string unit slice
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}

func main() {
	println(formatSize(1073741824))
}
