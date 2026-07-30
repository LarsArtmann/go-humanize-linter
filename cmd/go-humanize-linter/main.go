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

// versionDev is the placeholder used when the binary is built without
// ldflags-injected version metadata. It triggers a stderr warning at runtime
// so users notice they are running an unversioned build.
const versionDev = "dev"

// version is the CLI version. It is overridden at build time via -ldflags,
// e.g. -ldflags "-X main.version=v0.1.0".
var version = versionDev

// Output format identifiers. Single source of truth so the flag default, the
// switch cases, and the typed-error Format field cannot drift.
const (
	formatText  = "text"
	formatJSON  = "json"
	formatSARIF = "sarif"
)

// Failure-stage identifiers carried by *OutputError.Stage.
const (
	stageRender = "render"
	stageWrite  = "write"
)

func main() {
	var (
		enableIDs   stringList
		disableIDs  stringList
		format      string
		quiet       bool
		showVersion bool
		showRules   bool
		listFiles   bool
		explain     string
	)

	flag.Var(&enableIDs, "enable", "enable specific rule ID (repeatable, default: all)")
	flag.Var(&disableIDs, "disable", "disable specific rule ID (repeatable)")
	flag.StringVar(&format, "format", formatText, "output format: text, json, sarif")
	flag.BoolVar(&quiet, "quiet", false, "suppress summary line")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&showVersion, "v", false, "shorthand for --version")
	flag.BoolVar(&showRules, "rules", false, "list all rules with descriptions and exit")
	flag.BoolVar(&listFiles, "list-files", false, "list every Go file the walker would scan, then exit")
	flag.StringVar(&explain, "explain", "", "print the rationale for a rule (e.g. --explain H001) and exit")

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
		fmt.Printf("go-humanize-linter %s\n", version) //nolint:forbidigo // CLI stdout for --version

		if version == versionDev {
			fmt.Fprintln(os.Stderr, "warning: version is unset — built without -ldflags \"-X main.version=...\"")
		}

		return
	}

	if showRules {
		printRules()

		return
	}

	if listFiles {
		printScannedFiles(flag.Args())

		return
	}

	if explain != "" {
		printExplanation(explain)

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

	if err := output(os.Stdout, report, format, quiet); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	os.Exit(linter.ExitCodeFromReport(report))
}

// printRules writes a table of every rule (ID, name, severity, description) to
// stdout, then returns. Used by the --rules flag. Output goes to stdout so it
// pipes cleanly into grep / less / other tools.
func printRules() {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%-6s %-26s %-8s %s\n", "ID", "NAME", "SEV", "DESCRIPTION")

	for _, rule := range humanizelint.AllRules() {
		fmt.Fprintf(&builder, "%-6s %-26s %-8s %s\n",
			rule.Meta.ID, rule.Meta.Name, rule.Meta.Sev, rule.Meta.Description)
	}

	fmt.Print(builder.String()) //nolint:forbidigo // CLI stdout for --rules
}

// printScannedFiles lists every Go file the walker would scan, then returns.
// Used by the --list-files flag — useful for debugging "why is this file
// being scanned?" or "which files are excluded?".
func printScannedFiles(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: --list-files requires a directory argument")
		os.Exit(2)
	}

	files, err := humanizelint.WalkGoDir(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	for _, file := range files {
		fmt.Println(file.Path) //nolint:forbidigo // CLI stdout for --list-files
	}
}

// Rule ID constants. Duplicated from the humanizelint package (where they
// are unexported) because the CLI is the right place for the user-facing
// rationale text. Keep these in sync with humanizelint.AllRules().
const (
	h001 = "H001"
	h002 = "H002"
	h003 = "H003"
	h004 = "H004"
	h005 = "H005"
	h006 = "H006"
	h007 = "H007"
	h008 = "H008"
	h009 = "H009"
)

// ruleExplanations maps a rule ID to a longer rationale describing what
// it detects and what to use instead. The intent is to be useful when a
// user sees a finding and wants to understand the rule, not just the
// detection. These stay brief on purpose — they should fit on one screen.
var ruleExplanations = map[string]string{ //nolint:gochecknoglobals // CLI lookup table
	h001: "Manual byte-size formatting. The classic 'KMGTPE' index trick, " +
		"[]string of unit names, or repeated division by 1024 are signals that " +
		"the function is reimplementing humanize.Bytes / humanize.IBytes. " +
		"Both SI (KB/MB) and IEC (KiB/MiB) variants are covered.",
	h002: "Manual comma/thousands-separator insertion via mod-3 indexing, " +
		"step-by-3 loops, or for-loop digit grouping followed by a WriteString " +
		"',' is what humanize.Comma / humanize.Commaf do for you.",
	h003: "Manual relative-time formatting using time.Since/Sub + 'ago' " +
		"strings + time threshold comparisons. humanize.RelTime / " +
		"humanize.Time handle singular/plural forms and locale strings.",
	h004: "English pluralization via 'if n == 1' switches or singular/plural " +
		"parameter pairs. humanize.Plural / humanize.PluralWord cover 100+ " +
		"locales.",
	h005: "Manual SI-prefix formatting (1.5K, 2.3M) via division by 1000 " +
		"plus 'K'/'M' suffix strings. humanize.SI is a drop-in replacement.",
	h006: "Manual trailing-zero stripping via nested " +
		"strings.TrimRight(x, '0') + strings.TrimRight(..., '.'). " +
		"humanize.Ftoa handles the edge cases (e.g. '0' instead of '').",
	h007: "Manual byte-size string parsing via repeated HasSuffix checks " +
		"on KB/MB/GB suffixes or a map[string]int64 multiplier. " +
		"humanize.ParseBytes handles every SI/IEC unit and negative numbers.",
	h008: "Manual ordinal-number formatting via 'switch n%10' (or n%100) " +
		"with st/nd/rd/th cases. humanize.Ordinal covers the same ground plus " +
		"edge cases (11th, 12th, 13th) and multiple languages.",
	h009: "Manual float-with-thousands-separator formatting via '%.Nf' " +
		"Sprintf combined with a manual comma-grouping loop. " +
		"humanize.Commaf returns the same result in one call.",
}

// printExplanation prints the rationale for a given rule ID and returns. Used
// by the --explain flag.
func printExplanation(ruleID string) {
	explanation, ok := ruleExplanations[ruleID]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown rule: %s\n", ruleID)
		os.Exit(2)
	}

	fmt.Printf("Rule %s: %s\n", ruleID, explanation) //nolint:forbidigo
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

// OutputError is returned by output() when rendering or writing a report
// fails. It carries the format ("json", "sarif", "text") and the stage that
// failed ("render", "write") so callers can produce actionable messages and
// branch on concrete failure modes.
type OutputError struct {
	Format string
	Stage  string // "render" or "write"
	Err    error
}

// Error implements the error interface.
func (e *OutputError) Error() string {
	return fmt.Sprintf("render %s: %s: %v", e.Format, e.Stage, e.Err)
}

// Unwrap returns the underlying error so errors.Is / errors.AsType work
// through the chain.
func (e *OutputError) Unwrap() error {
	return e.Err
}

// output renders report to writer in the requested format. Returns an
// *OutputError if rendering or writing fails; the caller exits with status 2
// on error.
func output(writer io.Writer, report *finding.Report, format string, quiet bool) error {
	switch format {
	case formatJSON:
		data, err := report.JSON()
		if err != nil {
			return &OutputError{Format: formatJSON, Stage: stageRender, Err: err}
		}

		if _, err := fmt.Fprintln(writer, data); err != nil {
			return &OutputError{Format: formatJSON, Stage: stageWrite, Err: err}
		}

	case formatSARIF:
		if err := report.WriteSARIF(context.Background(), writer); err != nil {
			return &OutputError{Format: formatSARIF, Stage: stageRender, Err: err}
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

	return nil
}
