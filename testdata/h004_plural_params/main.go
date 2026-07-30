package main

import "fmt"

// Pattern: function with singular/plural params + if n == 1
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

func main() {
	println(pluralize(1, "file", "files"))
}
