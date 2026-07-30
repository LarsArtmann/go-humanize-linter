package main

import "fmt"

// formatBytes reimplements humanize.Bytes using the KMGTPE index trick.
// The analysistest framework asserts an H001 diagnostic on the func line.
func formatBytes(bytes int64) string { // want "H001"
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
