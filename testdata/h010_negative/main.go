package main

import (
	"errors"
	"strconv"
	"strings"
)

// csvField splits a CSV line and parses the first field. This is CSV field
// parsing, NOT comma-stripping: there is no strings.Join(y, "") rebuild and
// no ReplaceAll(x, ",", ""). Should NOT trigger H010.
func csvField(line string) (int64, error) {
	fields := strings.Split(line, ",")

	return strconv.ParseInt(fields[0], 10, 64)
}

// stripSpaces removes spaces (not commas) before parsing. Not a comma strip.
// Should NOT trigger H010.
func stripSpaces(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(s, " ", ""), 64)
}

// displayCleanup strips commas for display but never parses a number.
// Missing the corroborating strconv signal. Should NOT trigger H010.
func displayCleanup(s string) string {
	return strings.ReplaceAll(s, ",", "")
}

// replacerStrip strips commas via strings.NewReplacer. Not one of the
// detected strip forms (documented gap — see docs/rules/H010.md): stays
// clean for now. Should NOT trigger H010.
func replacerStrip(s string) (int64, error) {
	replacer := strings.NewReplacer(",", "")

	return strconv.ParseInt(replacer.Replace(s), 10, 64)
}

// rejectCommas validates that a string contains no commas, then parses it.
// The ',' comparison has no string rebuild — a validator, not a parser.
// Should NOT trigger H010.
func rejectCommas(s string) (int64, error) {
	for _, r := range s {
		if r == ',' {
			return 0, errors.New("commas not allowed")
		}
	}

	return strconv.ParseInt(s, 10, 64)
}

func main() {
	v, _ := csvField("123,456")
	println(v)
}
