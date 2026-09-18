package main

import (
	"strconv"
	"strings"
)

// parseAmount strips commas into an intermediate variable, then parses.
// Two corroborating signals in one function: comma strip + strconv parse.
func parseAmount(s string) (float64, error) {
	cleaned := strings.Replace(s, ",", "", -1)

	return strconv.ParseFloat(cleaned, 64)
}

func main() {
	v, _ := parseAmount("12,345.67")
	println(v)
}
