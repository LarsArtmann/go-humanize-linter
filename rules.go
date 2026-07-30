package humanizelint

import (
	"go/ast"
	"go/token"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// DefaultRegistry returns a Registry pre-loaded with all humanize-lint rules,
// all enabled by default.
func DefaultRegistry() *linter.Registry {
	registry := linter.NewRegistry()

	for _, rule := range AllRules() {
		registry.Register(rule)
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
	}
}

// ruleDetectors pairs each rule's detector with its ID so DetectFuncDecl can
// honour per-rule //nolint:gohumanize:Hxxx scoping.
type ruleDetectors struct {
	id       string
	detector func(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding
}

// allRuleDetectors is the canonical (id, detector) list — order matches
// AllRules(). Both //nolint suppression paths and the plugin entry point
// iterate over this.
func allRuleDetectors() []ruleDetectors {
	return []ruleDetectors{
		{"H001", detectBytesFormat},
		{"H002", detectCommaFormat},
		{"H003", detectRelTimeFormat},
		{"H004", detectPlural},
		{"H005", detectSIFormat},
		{"H006", detectFtoa},
		{"H007", detectParseBytes},
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
	if file == nil || fn == nil {
		return nil
	}

	detectors := allRuleDetectors()
	all := make([]finding.Finding, 0, len(detectors))

	for _, detector := range detectors {
		suppressions := funcSuppressions(fset, file, fn)
		if isSuppressedRule(suppressions, detector.id) {
			continue
		}

		all = append(all, detector.detector(fset, file, fn, filePath)...)
	}

	return all
}

// funcSuppressions collects every //nolint directive (parsed via
// suppressedRules) that applies to fn — either as the doc comment, a
// trailing comment, or a comment on the line before the declaration.
// Returns nil if no directive applies.
func funcSuppressions(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) []string {
	if file == nil || fn == nil {
		return nil
	}

	fnLine := fset.Position(fn.Pos()).Line

	var merged []string

	for _, group := range file.Comments {
		isDoc := group == fn.Doc

		for _, comment := range group.List {
			parts := suppressedRules(comment.Text)
			if parts == nil {
				continue
			}

			cmtLine := fset.Position(comment.Pos()).Line
			if !isDoc && cmtLine != fnLine && cmtLine != fnLine-1 {
				continue
			}

			merged = append(merged, parts...)
		}
	}

	return merged
}
