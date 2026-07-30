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
	r := linter.NewRegistry()

	for _, rule := range AllRules() {
		r.Register(rule)
	}

	return r
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

// DetectFuncDecl runs all enabled rules against a single [*ast.FuncDecl] and
// returns any findings. This is the per-function entry point used by the
// golangci-lint plugin wrapper (plugin/plugin.go) which iterates over
// pass.Files instead of walking a directory.
func DetectFuncDecl(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding {
	detectors := []func(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding{
		detectBytesFormat,
		detectCommaFormat,
		detectRelTimeFormat,
		detectPlural,
		detectSIFormat,
		detectFtoa,
		detectParseBytes,
	}

	var all []finding.Finding

	for _, d := range detectors {
		all = append(all, d(fset, file, fn, filePath)...)
	}

	return all
}
