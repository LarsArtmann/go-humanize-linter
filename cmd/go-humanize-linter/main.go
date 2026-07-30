// Command go-humanize-linter scans Go source files for hand-rolled
// reimplementations of github.com/dustin/go-humanize and reports them as
// findings.
//
// Usage:
//
//	go-humanize-linter [flags] <path>
//
// Flags:
//
//	--enable <id>    Enable a specific rule (repeatable). Default: all enabled.
//	--disable <id>   Disable a specific rule (repeatable).
//	--format <type>  Output format: text (default), json, sarif.
//	--quiet          Suppress summary line.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/larsartmann/go-finding"
	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-linter-sdk"
)

func main() {
	var (
		enableIDs  stringList
		disableIDs stringList
		format     string
		quiet      bool
	)

	flag.Var(&enableIDs, "enable", "enable specific rule ID (repeatable, default: all)")
	flag.Var(&disableIDs, "disable", "disable specific rule ID (repeatable)")
	flag.StringVar(&format, "format", "text", "output format: text, json, sarif")
	flag.BoolVar(&quiet, "quiet", false, "suppress summary line")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <path>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Scans Go source for reimplementations of go-humanize.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nRules:\n")

		for _, rule := range humanizelint.AllRules() {
			fmt.Fprintf(os.Stderr, "  %s  %s  %s\n", rule.Meta.ID, rule.Meta.Name, rule.Meta.Description)
		}
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	dir := args[0]

	registry := buildRegistry(enableIDs, disableIDs)

	report, err := registry.Run(context.Background(), dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	output(report, format, quiet)

	os.Exit(linter.ExitCodeFromReport(report))
}

// stringList implements flag.Value for repeatable string flags.
type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
	*s = append(*s, v)

	return nil
}

func buildRegistry(enableIDs, disableIDs []string) *linter.Registry {
	disabled := make(map[string]bool, len(disableIDs))
	for _, id := range disableIDs {
		disabled[id] = true
	}

	enabledOnly := len(enableIDs) > 0

	enableSet := make(map[string]bool, len(enableIDs))
	for _, id := range enableIDs {
		enableSet[id] = true
	}

	registry := linter.NewRegistry()

	for _, rule := range humanizelint.AllRules() {
		ruleID := rule.Meta.ID

		if disabled[ruleID] {
			continue
		}

		if enabledOnly && !enableSet[ruleID] {
			continue
		}

		registry.Register(rule)
	}

	return registry
}

func output(report *finding.Report, format string, quiet bool) {
	switch format {
	case "json":
		data, err := report.JSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "json error: %v\n", err)

			return
		}

		fmt.Println(data) //nolint:forbidigo // CLI stdout output

	case "sarif":
		if err := report.WriteSARIF(context.Background(), os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "sarif error: %v\n", err)

			return
		}

	default:
		for f := range report.All() {
			loc := fmt.Sprintf("%s:%d", f.Position.File, f.Position.Line)
			if f.Position.Column > 0 {
				loc += fmt.Sprintf(":%d", f.Position.Column)
			}

			suggestion := ""
			if f.Suggestion != "" {
				suggestion = "\n    💡 " + f.Suggestion
			}

			fmt.Printf("%s [%s] %s%s\n", loc, f.Rule, f.Message, suggestion) //nolint:forbidigo // CLI stdout output
		}

		if !quiet {
			fmt.Printf("\n%d findings\n", report.Len()) //nolint:forbidigo // CLI stdout output
		}
	}
}
