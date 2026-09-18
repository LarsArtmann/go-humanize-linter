package main

import (
	"strconv"
	"strings"
)

// parseTotal strips commas via Split + Join, then parses. Indirect strip
// form — Medium confidence.
func parseTotal(s string) (int64, error) {
	parts := strings.Split(s, ",")
	joined := strings.Join(parts, "")

	return strconv.ParseInt(joined, 10, 64)
}

func main() {
	v, _ := parseTotal("9,876,543")
	println(v)
}
