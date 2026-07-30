package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Pattern: multiple HasSuffix checks for byte-unit suffixes + multiplier
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if strings.HasSuffix(s, "GB") {
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "GB"), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid: %s", s)
		}
		return num * 1073741824, nil
	}
	if strings.HasSuffix(s, "MB") {
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "MB"), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid: %s", s)
		}
		return num * 1048576, nil
	}
	if strings.HasSuffix(s, "KB") {
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "KB"), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid: %s", s)
		}
		return num * 1024, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func main() {
	v, _ := parseSize("10MB")
	println(v)
}
