package main

import (
	"fmt"
	"strings"
)

// listPlain is a simple comma join without any conjunction handling — a
// plain list, not a word series. Should NOT trigger H012.
func listPlain(names []string) string {
	return fmt.Sprintf("%d emails (%s)", len(names), strings.Join(names, ", "))
}

// sentenceWithAnd contains the word "and" in prose but no strings.Join at
// all. Should NOT trigger H012.
func sentenceWithAnd(a, b string) string {
	return a + " and " + b
}

// joinPath joins with a slash separator; not a comma join. Should NOT
// trigger H012.
func joinPath(parts []string) string {
	return strings.Join(parts, "/")
}

func main() {
	println(listPlain([]string{"a", "b"}))
}
