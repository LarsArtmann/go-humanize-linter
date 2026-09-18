package main

import (
	num "strconv"
	str "strings"
)

// parseAliased strips commas via an aliased strings import and parses via an
// aliased strconv import. Tests that H010 resolves import aliases.
func parseAliased(s string) (int64, error) {
	cleaned := str.ReplaceAll(s, ",", "")

	return num.ParseInt(cleaned, 10, 64)
}

func main() {
	v, _ := parseAliased("4,567,890")
	println(v)
}
