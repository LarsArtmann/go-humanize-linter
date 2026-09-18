package main

import (
	"strconv"
	"strings"
)

// parseCount reimplements humanize.ParseComma: the comma strip is inlined
// directly in the strconv.ParseInt call.
func parseCount(s string) (int64, error) { // want "H010"
	return strconv.ParseInt(strings.ReplaceAll(s, ",", ""), 10, 64)
}

func main() {
	v, _ := parseCount("1,234,567")
	println(v)
}
