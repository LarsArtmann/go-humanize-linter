package main

import (
	"fmt"
	"strconv"
	"time"
)

// Pattern: time.Since + "ago" + time threshold comparisons
func relativeTime(t time.Time) string {
	elapsed := time.Since(t)
	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		return strconv.Itoa(int(elapsed.Minutes())) + "m ago"
	case elapsed < 24*time.Hour:
		return strconv.Itoa(int(elapsed.Hours())) + "h ago"
	default:
		return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
	}
}

func main() {
	println(relativeTime(time.Now()))
}
