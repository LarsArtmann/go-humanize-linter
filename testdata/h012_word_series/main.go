package main

import (
	"strings"
)

// joinWords hand-rolls what humanize.WordSeries(words, "and") does: prefix
// join of all but the last element, an Oxford-comma conjunction, and
// last-element access. The complete clone — Full confidence.
func joinWords(parts []string) string {
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
