package main

import (
	"strconv"
	. "strings"
)

// parseDotImport strips commas via a dot-imported strings.ReplaceAll (bare
// ReplaceAll call) and parses with strconv. Tests H010 on dot-imported
// strings.
func parseDotImport(s string) (float64, error) {
	cleaned := ReplaceAll(s, ",", "")

	return strconv.ParseFloat(cleaned, 64)
}

func main() {
	v, _ := parseDotImport("3,456.78")
	println(v)
}
