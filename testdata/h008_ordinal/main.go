package main

// ordinal reimplements humanize.Ordinal using a switch on n%10. The
// analysistest framework asserts an H008 diagnostic on the func line.
func ordinal(n int) string { // want "H008"
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func main() {
	println(ordinal(1), ordinal(2), ordinal(3), ordinal(11))
}
