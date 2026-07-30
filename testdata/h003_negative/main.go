package main

import (
	"fmt"
	"time"
)

// This uses time.Since but not for "ago" formatting — it's a timing check.
func isStale(t time.Time) bool {
	return time.Since(t) > time.Hour
}

// This mentions "ago" in a comment but not in code.
func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func main() {
	println(isStale(time.Now()))
	println(formatDuration(time.Second))
}
