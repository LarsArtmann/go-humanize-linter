package humanizelint

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// RulePlural (H004) detects manual English pluralization that should use
// humanize.Plural or humanize.PluralWord.
//
// Triggers when a function:
//   - Has parameters named "singular" and "plural" (explicit reimplementation)
//   - OR contains an `if x == 1` / `if x != 1` conditional where the branches
//     return different string values (the classic plural switch)
func RulePlural() linter.RuleFunc {
	return linter.RuleFunc{
		Meta: linter.RuleMeta{
			ID:          "H004",
			Name:        "manual-plural",
			Description: "Manual pluralization — use humanize.Plural or humanize.PluralWord instead of if-n==1 switches",
			Cat:         linter.CategoryStyle,
			Sev:         finding.SeverityWarning,
		},
		Run: func(_ context.Context, dir string) ([]finding.Finding, error) {
			return checkFuncDecls(dir, detectPlural)
		},
	}
}

func detectPlural(fset *token.FileSet, _ *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding {
	namedParams := hasPluralNamedParams(fn)
	equalsOne := hasEqualsOneBranch(fn)

	// For the equalsOne path, require the function to return a string type.
	// Without this, every `if result != 1 { return fmt.Errorf(...) }` in SQL
	// health checks, validation, and code generation triggers a false positive.
	if equalsOne && !namedParams && !funcReturnsString(fn) {
		equalsOne = false
	}

	if !namedParams && !equalsOne {
		return nil
	}

	confidence := finding.ConfidenceHigh
	if namedParams && equalsOne {
		confidence = finding.ConfidenceFull
	}

	line, col := posOf(fset, fn.Pos())

	return []finding.Finding{
		makeFindingWithConfidence(
			"H004",
			fmt.Sprintf(
				"manual pluralization (namedParams=%v, equalsOne=%v) — use humanize.Plural instead",
				namedParams, equalsOne,
			),
			"Replace with humanize.Plural(n, \"item\", \"\") or humanize.PluralWord(n, singular, plural).",
			line, col, filePath, confidence,
		),
	}
}
