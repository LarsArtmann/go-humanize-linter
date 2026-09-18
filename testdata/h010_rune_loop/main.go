package main

import (
	"strconv"
	"strings"
)

// parseRuneLoop strips commas with a manual rune loop that rebuilds the
// string via a Builder, then parses. The indirect strip form — Medium
// confidence. Tests the runeFilterLoop signal end to end.
func parseRuneLoop(s string) (int64, error) {
	var b strings.Builder
	for _, r := range s {
		if r != ',' {
			b.WriteRune(r)
		}
	}

	return strconv.ParseInt(b.String(), 10, 64)
}

func main() {
	v, _ := parseRuneLoop("9,876,543")
	println(v)
}
