package main

import (
	"fmt"
	"strconv"
	. "time"
)

// relativeTimeDotImport uses a dot-imported time package, so the duration
// constants appear as bare identifiers (Minute, Hour, Day) rather than
// time.Minute, time.Hour, etc. H003 must still detect this.
func relativeTimeDotImport(t Time) string {
	elapsed := Since(t)
	switch {
	case elapsed < Minute:
		return "just now"
	case elapsed < Hour:
		return strconv.Itoa(int(elapsed.Minutes())) + "m ago"
	case elapsed < 24*Hour:
		return strconv.Itoa(int(elapsed.Hours())) + "h ago"
	default:
		return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
	}
}

func main() {
	println(relativeTimeDotImport(Now()))
}
