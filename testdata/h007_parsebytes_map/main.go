package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Pattern: map[string]int64 multiplier lookup
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	multipliers := map[string]int64{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
	for suffix, mult := range multipliers {
		if before, ok := strings.CutSuffix(s, suffix); ok {
			num, err := strconv.ParseInt(strings.TrimSpace(before), 10, 64)
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
