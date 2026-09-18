package main

import (
	"strconv"
	"strings"
)

// Pattern with a scoped H010 directive: this MUST be suppressed for H010.
//
//nolint:gohumanize:H010
func parseCount(s string) (int64, error) {
	return strconv.ParseInt(strings.ReplaceAll(s, ",", ""), 10, 64)
}

func main() {
	v, _ := parseCount("123,456")
	println(v)
}
