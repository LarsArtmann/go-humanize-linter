package main

import (
	"strconv"
	"strings"
)

var byteMultipliers = map[string]int64{
	"KB": 1000,
	"MB": 1000000,
	"GB": 1000000000,
}

func parseByteSize(s string) int64 {
	for unit, mult := range byteMultipliers {
		if strings.HasSuffix(s, unit) {
			n, _ := strconv.Atoi(s[:len(s)-len(unit)])
			return int64(n) * mult
		}
	}

	return 0
}

func main() {}
