package main

import (
	"fmt"
	"time"
)

// relativeTime reimplements humanize.RelTime using time.Since + "ago" + time
// thresholds. The analysistest framework asserts an H003 diagnostic on the
// func line.
func relativeTime(t time.Time) string { // want "H003"
	elapsed := time.Since(t)
	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		return fmt.Sprintf("%dm ago", int(elapsed.Minutes()))
	case elapsed < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(elapsed.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
	}
}

func main() {
	println(relativeTime(time.Now()))
}