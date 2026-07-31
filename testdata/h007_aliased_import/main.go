package main

import str "strings"

func parseSize(s string) int64 {
	if str.HasSuffix(s, "KB") {
		return 1000
	}

	if str.HasSuffix(s, "MB") {
		return 1000000
	}

	if str.HasSuffix(s, "GB") {
		return 1000000000
	}

	return 0
}

func main() {}
