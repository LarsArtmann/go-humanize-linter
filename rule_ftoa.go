package humanizelint

import (
	"context"
	"go/ast"
	"go/token"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// RuleFtoa (H006) detects manual float-to-string conversion with trailing-zero
// stripping that should use humanize.Ftoa.
//
// Triggers when a function contains:
//
//	strings.TrimRight(strings.TrimRight(x, "0"), ".")
//
// This nested-TrimRight pattern is the textbook reimplementation of
// humanize.Ftoa's trailing-zero removal.
func RuleFtoa() linter.RuleFunc {
	return linter.RuleFunc{
		Meta: linter.RuleMeta{
			ID:          "H006",
			Name:        "manual-ftoa",
			Description: "Manual float formatting with trailing-zero stripping — use humanize.Ftoa instead of nested strings.TrimRight",
			Cat:         linter.CategoryStyle,
			Sev:         finding.SeverityWarning,
		},
		Run: func(_ context.Context, dir string) ([]finding.Finding, error) {
			return checkFuncDecls(dir, detectFtoa)
		},
	}
}

func detectFtoa(fset *token.FileSet, _ *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding {
	if !hasNestedTrimRight(fn) {
		return nil
	}

	line, col := posOf(fset, fn.Pos())

	return []finding.Finding{
		makeFindingWithConfidence(
			"H006",
			"manual float trailing-zero stripping (nested TrimRight) — use humanize.Ftoa instead",
			"Replace with humanize.Ftoa(f) or humanize.FtoaWithDigits(f, 1).",
			line, col, filePath, finding.ConfidenceFull,
		),
	}
}
