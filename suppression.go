package humanizelint

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/larsartmann/go-finding"
)

// SuppressionDirective records a //nolint-style directive that is attached
// to a function declaration. It is used by --verify-suppressions to detect
// stale or misspelled suppressions.
type SuppressionDirective struct {
	FilePath     string
	FunctionLine int
	Line         int
	Column       int
	RawText      string
	Suppressed   []string
}

// allRuleIDs returns the stable rule identifiers H001-H009. Kept as a function
// rather than a global so the set cannot drift from the canonical constants.
func allRuleIDs() []string {
	return []string{
		RuleIDH001,
		RuleIDH002,
		RuleIDH003,
		RuleIDH004,
		RuleIDH005,
		RuleIDH006,
		RuleIDH007,
		RuleIDH008,
		RuleIDH009,
	}
}

// suppressesGohumanize reports whether the parsed suppression list mentions our
// linter (gohumanize) or any of our rule IDs.
func suppressesGohumanize(suppressed []string) bool {
	if isSuppressedAll(suppressed) {
		return true
	}

	for _, ruleID := range allRuleIDs() {
		if isSuppressedRule(suppressed, ruleID) {
			return true
		}
	}

	return false
}

// hasUnknownHumanizeLinterName reports whether the suppression list contains a
// linter name that looks like ours but is not exactly "gohumanize". The common
// AI mistake is writing "go-humanize-linter" because that matches the module
// path. Such directives silently do nothing.
func hasUnknownHumanizeLinterName(suppressed []string) bool {
	for _, s := range suppressed {
		lower := strings.ToLower(s)

		if lower == "gohumanize" || lower == "all" {
			continue
		}

		if strings.Contains(lower, "humanize") {
			return true
		}
	}

	return false
}

// collectSuppressions walks dir and returns every //nolint-style directive that
// either suppresses our linter or misspells its name.
func collectSuppressions(dir string) ([]SuppressionDirective, error) {
	files, err := WalkGoDir(dir)
	if err != nil {
		return nil, err
	}

	var out []SuppressionDirective

	for _, pf := range files {
		for _, decl := range pf.File.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			out = append(out, extractFunctionSuppressions(pf.Fset, pf.File, fn, pf.Path)...)
		}
	}

	return out, nil
}

func extractFunctionSuppressions(
	fset *token.FileSet,
	file *ast.File,
	fn *ast.FuncDecl,
	filePath string,
) []SuppressionDirective {
	if file == nil || fn == nil {
		return nil
	}

	fnLine := fset.Position(fn.Pos()).Line
	var out []SuppressionDirective

	for _, group := range file.Comments {
		isDoc := group == fn.Doc

		for _, comment := range group.List {
			suppressed := suppressedRules(comment.Text)
			if suppressed == nil {
				continue
			}

			cmtLine := fset.Position(comment.Pos()).Line
			if !isDoc && cmtLine != fnLine && cmtLine != fnLine-1 {
				continue
			}

			pos := fset.Position(comment.Pos())

			out = append(out, SuppressionDirective{
				FilePath:     filePath,
				FunctionLine: fnLine,
				Line:         pos.Line,
				Column:       pos.Column,
				RawText:      comment.Text,
				Suppressed:   suppressed,
			})
		}
	}

	return out
}

// VerifySuppressions reports two classes of suppression problems:
//   1. Unknown linter names that look like "gohumanize" (e.g. "go-humanize-linter").
//   2. //nolint:gohumanize[:Hxxx] directives that did not suppress any finding.
//
// The returned findings use the pseudo-rule ID "H0SUP" so they are clearly
// verification diagnostics, not humanize reimplementation findings.
func VerifySuppressions(dir string, report *finding.Report) ([]finding.Finding, error) {
	directives, err := collectSuppressions(dir)
	if err != nil {
		return nil, err
	}

	findingsByPosition := findingsByFileLine(report)

	var out []finding.Finding

	for _, d := range directives {
		if hasUnknownHumanizeLinterName(d.Suppressed) {
			out = append(out, makeSuppressionVerificationFinding(
				d, "unknown linter name in //nolint directive (did you mean gohumanize?)",
			))

			continue
		}

		if !suppressesGohumanize(d.Suppressed) {
			continue
		}

		if !directiveMatchesFinding(d, d.Suppressed, findingsByPosition) {
			out = append(out, makeSuppressionVerificationFinding(
				d, "//nolint:gohumanize directive suppresses zero findings in this function",
			))
		}
	}

	return out, nil
}

// findingsByFileLine indexes the report's findings by file path and line
// number. This is the lookup structure used to test whether a suppression
// directive actually suppressed something.
func findingsByFileLine(report *finding.Report) map[string]map[int][]finding.Finding {
	out := make(map[string]map[int][]finding.Finding)

	if report == nil {
		return out
	}

	for f := range report.All() {
		file := string(f.Position.File)
		line := f.Position.Line

		if out[file] == nil {
			out[file] = make(map[int][]finding.Finding)
		}

		out[file][line] = append(out[file][line], f)
	}

	return out
}

// directiveMatchesFinding reports whether any finding in the same function
// (same file and function-start line) is covered by the suppression directive.
func directiveMatchesFinding(
	d SuppressionDirective,
	suppressed []string,
	findingsByPosition map[string]map[int][]finding.Finding,
) bool {
	fileFindings, ok := findingsByPosition[d.FilePath]
	if !ok {
		return false
	}

	for _, f := range fileFindings[d.FunctionLine] {
		if isSuppressedRule(suppressed, string(f.Rule)) {
			return true
		}
	}

	return false
}

const suppressionVerificationRuleID = "H0SUP"

func makeSuppressionVerificationFinding(d SuppressionDirective, message string) finding.Finding {
	return finding.NewBuilder(
		finding.RuleName(suppressionVerificationRuleID),
		finding.ToolName("go-humanize-linter"),
		message,
		finding.SeverityWarning,
		finding.Pos(finding.FilePath(d.FilePath), d.Line, d.Column),
	).
		WithCategory(finding.CategoryConfiguration).
		WithConfidence(finding.ConfidenceHigh).
		MustBuild()
}

// VerifySuppressionComment is a convenience helper for tests: it parses a
// single comment string and returns whether it suppresses our linter and
// whether it contains a misspelled linter name.
func VerifySuppressionComment(comment string) (suppressesOurs, hasTypo bool) {
	suppressed := suppressedRules(comment)
	if suppressed == nil {
		return false, false
	}

	return suppressesGohumanize(suppressed), hasUnknownHumanizeLinterName(suppressed)
}
