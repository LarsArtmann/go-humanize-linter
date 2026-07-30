package main

// dayOfWeek uses switch n%10 but returns weekday names, not ordinal suffixes.
// Should NOT trigger H008.
func dayOfWeek(n int) string {
	switch n % 10 {
	case 0:
		return "Sunday"
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	default:
		return "Other"
	}
}

// remainderClass switches on n%7 (not n%10 or n%100). Should NOT trigger H008.
func remainderClass(n int) string {
	switch n % 7 {
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

// onlyTwoSuffixes returns only 2 of the 4 ordinal suffixes (below the
// minOrdinalSuffixesHit threshold of 3). Should NOT trigger H008.
func onlyTwoSuffixes(n int) string {
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	default:
		return "other"
	}
}

// spelledOutOrdinals returns full words, not 2-char suffixes. Should NOT
// trigger H008.
func spelledOutOrdinals(n int) string {
	switch n % 10 {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	default:
		return "other"
	}
}

func main() {
	_ = dayOfWeek(3)
	_ = remainderClass(10)
	_ = onlyTwoSuffixes(1)
	_ = spelledOutOrdinals(2)
}
