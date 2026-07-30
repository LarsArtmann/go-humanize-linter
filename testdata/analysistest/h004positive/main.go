package main

import "fmt"

// pluralItems reimplements humanize.Plural via an if-n==1 switch with string
// branches. The analysistest framework asserts an H004 diagnostic on the func
// line.
func pluralItems(n int) string { // want "H004"
	if n == 1 {
		return fmt.Sprintf("%d item", n)
	}

	return fmt.Sprintf("%d items", n)
}

func main() {
	println(pluralItems(5))
}
