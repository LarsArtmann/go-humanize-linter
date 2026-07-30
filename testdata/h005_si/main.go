package main

import "strconv"

// Pattern: division by 1000 + "K"/"M" suffix
func formatNumber(num float64) string {
	if num >= 1000000 {
		return strconv.FormatFloat(num/1000000, 'f', 1, 64) + "M"
	}
	if num >= 1000 {
		return strconv.FormatFloat(num/1000, 'f', 1, 64) + "K"
	}
	return strconv.FormatFloat(num, 'f', 0, 64)
}

func main() {
	println(formatNumber(1234567))
}
