package main

import (
	"fmt"
	"strings"
)

// listPlain is a simple comma join without any conjunction handling — a
// plain list, not a word series. Should NOT trigger H012.
func listPlain(names []string) string {
	return fmt.Sprintf("%d emails (%s)", len(names), strings.Join(names, ", "))
}

// sentenceWithAnd contains the word "and" in prose but no strings.Join at
// all. Should NOT trigger H012.
func sentenceWithAnd(a, b string) string {
	return a + " and " + b
}

// joinPath joins with a slash separator; not a comma join. Should NOT
// trigger H012.
func joinPath(parts []string) string {
	return strings.Join(parts, "/")
}

// summaryWithProseAnd pairs a comma join with the word "and" in prose (a
// Sprintf format string). The conjunction is never a Join separator or a
// concatenation operand, so the position filter must keep this clean —
// this is the bug-tracking-schema false-positive shape that motivated the
// filter. Should NOT trigger H012.
func summaryWithProseAnd(names []string) string {
	return fmt.Sprintf("open and closed issues: %s", strings.Join(names, ", "))
}

// headingWithProseConjunction has a comma join and a bare conjunction
// string literal in prose position (a variable), never in join position.
// Should NOT trigger H012.
func headingWithProseConjunction(names []string) string {
	title := "cause and effect"
	return title + ": " + strings.Join(names, ", ")
}

// prefixJoinNoConjunction has the WordSeries prefix shape
// (Join(x[:len(x)-1], ", ")) but never introduces a conjunction. A
// truncated list is not a word series. Should NOT trigger H012.
func prefixJoinNoConjunction(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts[:len(parts)-1], ", ") + "..."
}

func main() {
	println(listPlain([]string{"a", "b"}))
	println(sentenceWithAnd("a", "b"))
	println(joinPath([]string{"a", "b"}))
	println(summaryWithProseAnd([]string{"a", "b"}))
	println(headingWithProseConjunction([]string{"a", "b"}))
	println(prefixJoinNoConjunction([]string{"a", "b", "c"}))
}
