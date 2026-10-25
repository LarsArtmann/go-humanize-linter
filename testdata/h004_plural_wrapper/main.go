package main

import (
	"fmt"

	"github.com/dustin/go-humanize/english"
)

// Pattern: signature-only matches (singular/plural params) whose bodies
// already delegate to english.Plural / english.PluralWord. H004 must NOT
// fire — these are legitimate wrappers, not reimplementations.

func countMeta(count int, singular, plural string) string {
	if count == 0 {
		return ""
	}

	return english.Plural(count, singular, plural)
}

func wordMeta(count int, singular, plural string) string {
	return english.PluralWord(count, singular, plural)
}

func main() {
	println(countMeta(0, "file", "files"))
	println(wordMeta(2, "entry", "entries"))
	println(fmt.Sprint("ok"))
}
