package main

import (
	"strconv"
	"strings"
)

const digitsPerGroup = 3

// Pattern: named constant step size + comma insertion + digit conversion.
// No literal 3 in the code, so the fallback path must catch this.
func formatInt(n int64) string {
	digits := strconv.FormatInt(n, 10)
	if len(digits) <= digitsPerGroup {
		return digits
	}
	var b strings.Builder
	pre := len(digits) % digitsPerGroup
	if pre > 0 {
		b.WriteString(digits[:pre])
	}
	for i := pre; i < len(digits); i += digitsPerGroup {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(digits[i : i+digitsPerGroup])
	}
	return b.String()
}

func main() {
	println(formatInt(1234567))
}
