// Package humanizelint detects hand-rolled reimplementations of
// github.com/dustin/go-humanize and suggests the library function instead.
//
// The linter scans Go source files using AST pattern matching. Each rule looks
// for a cluster of signals that, together, strongly indicate a manual
// reimplementation of a specific go-humanize feature. Single weak signals are
// never enough — the rules require multiple corroborating signals within the
// same function to keep false positives near zero.
//
// # Rules
//
//	H001  manual-bytes-format     → humanize.Bytes / humanize.IBytes
//	H002  manual-comma-format     → humanize.Comma / humanize.Commaf
//	H003  manual-reltime-format   → humanize.RelTime / humanize.Time
//	H004  manual-plural           → humanize.Plural / humanize.PluralWord
//	H005  manual-si-format        → humanize.SI / humanize.SIWithDigits
//	H006  manual-ftoa             → humanize.Ftoa / humanize.FtoaWithDigits
//
// Usage as a library:
//
//	import (
//	    humanizelint "github.com/larsartmann/go-humanize-linter"
//	    "github.com/larsartmann/go-linter-sdk"
//	)
//
//	reg := humanizelint.DefaultRegistry()
//	report, err := reg.Run(ctx, ".")
//	fmt.Println(linter.ExitCodeFromReport(report))
package humanizelint
