package main

import (
	"strings"
)

// joinWords reimplements humanize.WordSeries(words, "and").
func joinWords(parts []string) string { // want "H012"
	if len(parts) == 1 {
		return parts[0]
	}

	if len(parts) == 2 {
		return parts[0] + " and " + parts[1]
	}

	return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
}

func main() {
	println(joinWords([]string{"a", "b", "c"}))
}
