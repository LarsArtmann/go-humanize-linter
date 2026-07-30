package main

import (
	"fmt"
	"strconv"
)

// parseSize reimplements humanize.ParseBytes using a map[string]int64
// multiplier. The analysistest framework asserts an H007 diagnostic on the
// func line.
//
//nolint:gohumanize:H001 // H001 is suppressed because the map keys (KB, MB, GB, TB) legitimately match the byte-unit detector; we only want H007.
func parseSize(s string) (int64, error) { // want "H007"
	multipliers := map[string]int64{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
	for suffix, mult := range multipliers {
		if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
			num, err := strconv.ParseInt(s[:len(s)-len(suffix)], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid: %s", s)
			}

			return num * mult, nil
		}
	}

	return strconv.ParseInt(s, 10, 64)
}

func main() {
	v, _ := parseSize("10MB")
	println(v)
}