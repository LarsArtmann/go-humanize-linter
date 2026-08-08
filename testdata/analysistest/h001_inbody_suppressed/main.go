package main

import "fmt"

// formatBytesInBody reimplements humanize.IBytes but the //nolint directive
// is placed INSIDE the function body (on the return statement line).
// The plugin path must honour this as a suppression — no diagnostic expected.
func formatBytesInBody(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp]) //nolint:gohumanize // in-body suppression
}

func main() {
	println(formatBytesInBody(1073741824))
}
