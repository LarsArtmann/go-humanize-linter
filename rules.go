package humanizelint

import (
	"context"
	"go/ast"
	"go/token"
	"sync"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// Rule ID constants — the single source of truth for rule identifiers.
// All packages (CLI, plugin, detectors) must reference these instead of
// string literals to prevent drift when rules are added or renamed.
const (
	RuleIDH001 = "H001"
	RuleIDH002 = "H002"
	RuleIDH003 = "H003"
	RuleIDH004 = "H004"
	RuleIDH005 = "H005"
	RuleIDH006 = "H006"
	RuleIDH007 = "H007"
	RuleIDH008 = "H008"
	RuleIDH009 = "H009"
	RuleIDH010 = "H010"
	// RuleIDH012 detects manual word-series joining. RuleIDH011 is reserved
	// for manual-si-parse (deferred: zero corpus demand in the 2026-09-18
	// research; see ROADMAP "H011+") — do not reuse it for another rule.
	RuleIDH012 = "H012"
)

// RuleIDH0SUP is the pseudo-rule ID for suppression-verification diagnostics.
// It is NOT in AllRules() and cannot be enabled/disabled via --enable/--disable.
// It is produced only by VerifySuppressions and must bypass confidence filtering
// so that stale or misspelled //nolint directives are always surfaced.
const RuleIDH0SUP = "H0SUP"

// toolName is the finding tool name for all humanize-lint findings.
const toolName finding.ToolName = "go-humanize-linter"

// coverageRule wraps a walker-based RuleFunc so the parsed-file count of its
// most recent scan flows into the report as filesScanned evidence
// (linter.CoverageReporter). DefaultRegistry installs the wrapper around
// every rule; the count itself is recorded by checkFuncDecls through a
// per-Check context sink, so concurrent registry runs never share state.
// Note: AllRules() keeps returning plain RuleFunc values (public API —
// BuildFlow reads rule.Meta.ID off them), so cherry-pickers that build their
// own registry from AllRules get coverage-silent rules by construction.
type coverageRule struct {
	linter.RuleFunc

	scanned int
}

// Coverage implements linter.CoverageReporter with the count recorded by the
// most recent Check. Files the walker failed to read or parse are not counted
// — WalkGoDir skips them silently, and evidence must be honest.
func (r *coverageRule) Coverage(_ context.Context, _ string) (int, []string) {
	return r.scanned, nil
}

// Check runs the wrapped rule with a scan sink in the context so the walker
// records exactly how many files THIS rule's scan parsed.
func (r *coverageRule) Check(ctx context.Context, dir string) ([]finding.Finding, error) {
	return r.RuleFunc.Check(contextWithScanSink(ctx, &r.scanned), dir)
}

// DefaultRegistry returns a Registry pre-loaded with all humanize-lint rules,
// all enabled by default. Every rule is wrapped for coverage reporting: the
// registry's reports carry real filesScanned evidence instead of a
// meaningless zero (fleet dashboards cannot otherwise tell "clean tree"
// from "scanned nothing").
func DefaultRegistry() *linter.Registry {
	registry := linter.NewRegistry()

	for _, rule := range AllRules() {
		registry.Register(&coverageRule{RuleFunc: rule})
	}

	return registry
}

// AllRules returns every rule in this linter as a slice. Useful for consumers
// that want to cherry-pick rules into their own registry.
func AllRules() []linter.RuleFunc {
	return []linter.RuleFunc{
		RuleBytes(),
		RuleComma(),
		RuleRelTime(),
		RulePlural(),
		RuleSI(),
		RuleFtoa(),
		RuleParseBytes(),
		RuleOrdinal(),
		RuleCommaf(),
		RuleParseComma(),
		RuleWordSeries(),
	}
}

// HumanizeDetector is the facade that owns the per-function detection
// pipeline: it iterates over every detector, applies per-rule suppression, and
// returns the aggregated findings. Use it instead of DetectFuncDecl when you
// need to configure the detector set, plug a custom suppression resolver, or
// share state (e.g. metrics) across calls.
//
// Construction:
//
//	detector := humanizelint.NewHumanizeDetector()    // all 11 rules enabled
//	detector := humanizelint.NewHumanizeDetector(     // opt-in subset
//	    humanizelint.RuleBytes(),
//	    humanizelint.RuleComma(),
//	)
//
// Then run it on a single function:
//
//	findings := detector.Run(fset, file, fn, "demo.go")
type HumanizeDetector struct {
	detectors []ruleDetectors
}

// NewHumanizeDetector constructs a HumanizeDetector running the given rules
// in the given order. If no rules are passed, all 11 default rules are
// registered.
func NewHumanizeDetector(rules ...linter.RuleFunc) *HumanizeDetector {
	if len(rules) == 0 {
		rules = AllRules()
	}

	d := &HumanizeDetector{
		detectors: make([]ruleDetectors, 0, len(rules)),
	}

	byID := cachedDetectorByID()

	for _, rule := range rules {
		det, ok := byID[rule.Meta.ID]
		if !ok {
			continue
		}

		d.detectors = append(d.detectors, ruleDetectors{id: rule.Meta.ID, detector: det})
	}

	return d
}

// Run applies every registered detector to a single function, honouring
// per-rule //nolint:gohumanize[:Hxxx] suppression directives. Returns the
// aggregated findings (nil if none).
func (d *HumanizeDetector) Run(
	fset *token.FileSet,
	file *ast.File,
	fn *ast.FuncDecl,
	filePath string,
) []finding.Finding {
	if file == nil || fn == nil {
		return nil
	}

	suppressions := funcSuppressions(fset, file, fn)

	var all []finding.Finding

	for _, detector := range d.detectors {
		if isSuppressedRule(suppressions, detector.id) {
			continue
		}

		all = append(all, detector.detector(fset, file, fn, filePath)...)
	}

	return all
}

// ruleDetectors pairs each rule's detector with its ID so DetectFuncDecl can
// honour per-rule //nolint:gohumanize:Hxxx scoping.
type ruleDetectors struct {
	id       string
	detector detectorFn
}

// detectorFn is the signature every per-function detector implements.
type detectorFn = func(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding

// cachedDetectorByID builds the (rule-ID → detector) lookup map once and
// reuses it across all HumanizeDetector instances.
var cachedDetectorByID = sync.OnceValue(func() map[string]detectorFn { //nolint:gochecknoglobals // cached
	m := make(map[string]detectorFn)

	for _, det := range allRuleDetectors() {
		m[det.id] = det.detector
	}

	return m
})

// allRuleDetectors is the canonical (id, detector) list — order matches
// AllRules(). Both //nolint suppression paths and the plugin entry point
// iterate over this.
func allRuleDetectors() []ruleDetectors {
	return []ruleDetectors{
		{RuleIDH001, detectBytesFormat},
		{RuleIDH002, detectCommaFormat},
		{RuleIDH003, detectRelTimeFormat},
		{RuleIDH004, detectPlural},
		{RuleIDH005, detectSIFormat},
		{RuleIDH006, detectFtoa},
		{RuleIDH007, detectParseBytes},
		{RuleIDH008, detectOrdinal},
		{RuleIDH009, detectCommaf},
		{RuleIDH010, detectParseComma},
		{RuleIDH012, detectWordSeries},
	}
}

// DetectFuncDecl runs all enabled rules against a single [*ast.FuncDecl] and
// returns any findings. This is the per-function entry point used by the
// golangci-lint plugin wrapper (plugin/plugin.go) which iterates over
// pass.Files instead of walking a directory.
//
// Suppression via //nolint:gohumanize (or //nolint:all, or scoped
// //nolint:gohumanize:H001) is honoured per-rule. A bare //nolint suppresses
// every rule; //nolint:gohumanize suppresses only this linter.
func DetectFuncDecl(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding {
	return NewHumanizeDetector().Run(fset, file, fn, filePath)
}

// funcSuppressions collects every //nolint directive (parsed via
// suppressedRules) that applies to fn — either as the doc comment, a
// trailing comment, a comment on the line before the declaration, or
// a comment inside the function body. Returns nil if no directive applies.
func funcSuppressions(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []string {
	if file == nil || fn == nil {
		return nil
	}

	var merged []string

	for _, group := range file.Comments {
		isDoc := group == fn.Doc

		for _, comment := range group.List {
			parts := suppressedRules(comment.Text)
			if parts == nil {
				continue
			}

			cmtLine := fset.Position(comment.Pos()).Line
			if !commentAssociatedWithFunc(fset, fn, isDoc, cmtLine) {
				continue
			}

			merged = append(merged, parts...)
		}
	}

	return merged
}
