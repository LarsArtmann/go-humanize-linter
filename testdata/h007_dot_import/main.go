package main

import (
	. "strings"
)

// parseSizeDotImport reimplements humanize.ParseBytes using dot-imported
// strings functions (HasSuffix, TrimSuffix). Tests that H007 detection
// works when strings is dot-imported.
func parseSizeDotImport(s string) (uint64, bool) {
	multipliers := map[string]uint64{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}

	for suffix, mult := range multipliers {
		if HasSuffix(s, suffix) {
			return mult, true
		}
		if TrimSuffix(s, suffix) != s {
			return mult, true
		}
	}

	return 0, false
}
