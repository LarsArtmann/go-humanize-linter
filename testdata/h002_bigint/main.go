package main

import (
	"math/big"
	"strings"
)

// bigComma groups the digits of a big.Int with commas — the hand-rolled
// equivalent of humanize.BigComma. H002 fires and the suggestion must point
// at humanize.BigComma because the value is a *big.Int.
func bigComma(b *big.Int) string {
	s := b.String()

	var out strings.Builder

	start := 0
	if len(s)%3 != 0 {
		start = len(s) % 3
		out.WriteString(s[:start])
	}

	for i := start; i < len(s); i += 3 {
		if i > 0 {
			out.WriteString(",")
		}

		out.WriteString(s[i : i+3])
	}

	return out.String()
}

func main() {
	println(bigComma(big.NewInt(1234567)))
}
