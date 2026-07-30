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
//	--rules          List all rules with descriptions and exit.
//	--version, -v    Print version and exit.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/larsartmann/go-finding"
	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-linter-sdk"
)

// version is the CLI version. It is overridden at build time via -ldflags,
// e.g. -ldflags "-X main.version=v0.1.0".
var version = "dev"

func main() {
	var (
		enableIDs   stringList
		disableIDs  stringList
		format      string
		quiet       bool
		showVersion bool
		showRules   bool
	)

	flag.Var(&enableIDs, "enable", "enable specific rule ID (repeatable, default: all)")
	flag.Var(&disableIDs, "disable", "disable specific rule ID (repeatable)")
	flag.StringVar(&format, "format", "text", "output format: text, json, sarif")
	flag.BoolVar(&quiet, "quiet", false, "suppress summary line")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&showVersion, "v", false, "shorthand for --version")
	flag.BoolVar(&showRules, "rules", false, "list all rules with descriptions and exit")

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

	if showVersion {
		fmt.Printf("go-humanize-linter %s\n", version) //nolint:forbidigo // CLI stdout output

		return
	}

	if showRules {
		printRules()

		return
	}

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

	output(os.Stdout, report, format, quiet)

	os.Exit(linter.ExitCodeFromReport(report))
}

// printRules writes a table of every rule (ID, name, severity, description) to
// stdout, then returns. Used by the --rules flag.
func printRules() {
	fmt.Fprintf(os.Stderr, "%-6s %-26s %-8s %s\n", "ID", "NAME", "SEV", "DESCRIPTION")

	for _, rule := range humanizelint.AllRules() {
		fmt.Fprintf(os.Stderr, "%-6s %-26s %-8s %s\n",
			rule.Meta.ID, rule.Meta.Name, rule.Meta.Sev, rule.Meta.Description)
	}
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

func output(writer io.Writer, report *finding.Report, format string, quiet bool) {
	switch format {
	case "json":
		data, err := report.JSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "json error: %v\n", err)

			return
		}

		fmt.Fprintln(writer, data)

	case "sarif":
		if err := report.WriteSARIF(context.Background(), writer); err != nil {
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

			fmt.Fprintf(writer, "%s [%s] %s%s\n", loc, f.Rule, f.Message, suggestion)
		}

		if !quiet {
			fmt.Fprintf(writer, "\n%d findings\n", report.Len())
		}
	}
}
