package main

import "strconv"

// This function divides by 1000 but formats milliseconds, not SI prefixes.
func formatMillis(ms float64) string {
	return strconv.FormatFloat(ms/1000, 'f', 2, 64) + "s"
}

func main() {
	println(formatMillis(1500))
}
