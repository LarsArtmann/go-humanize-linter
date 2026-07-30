package main

import "fmt"

// Pattern: if n == 1 returning different strings (no singular/plural param names)
func pluralStars(n int) string {
	if n == 1 {
		return fmt.Sprintf("%d star", n)
	}
	return fmt.Sprintf("%d stars", n)
}

func main() {
	println(pluralStars(5))
}
