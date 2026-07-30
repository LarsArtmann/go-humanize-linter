package humanizelint

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// RuleParseBytes (H007) detects manual byte-size string parsing ("10MB" →
// bytes) that should use humanize.ParseBytes.
//
// Triggers when a function contains:
//   - 2+ strings.HasSuffix/CutSuffix/TrimSuffix calls checking byte-unit suffixes
//   - OR a map[string]int64 literal with 2+ byte-unit keys used as multipliers
func RuleParseBytes() linter.RuleFunc {
	return linter.RuleFunc{
		Meta: linter.RuleMeta{
			ID:          "H007",
			Name:        "manual-parse-bytes",
			Description: "Manual byte-size string parsing — use humanize.ParseBytes instead of suffix matching and multiplication",
			Cat:         linter.CategoryStyle,
			Sev:         finding.SeverityWarning,
		},
		Run: func(_ context.Context, dir string) ([]finding.Finding, error) {
			return checkFuncDecls(dir, detectParseBytes)
		},
	}
}

func detectParseBytes(fset *token.FileSet, _ *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding {
	suffixChecks := hasByteUnitSuffixChecks(fn)
	multMap := hasByteUnitMultiplierMap(fn)

	if suffixChecks < 2 && !multMap {
		return nil
	}

	var confidence finding.Confidence

	switch {
	case suffixChecks >= 3 || multMap:
		confidence = finding.ConfidenceHigh
	case suffixChecks >= 2:
		confidence = finding.ConfidenceMedium
	default:
		return nil
	}

	line, col := posOf(fset, fn.Pos())

	return []finding.Finding{
		makeFindingWithConfidence(
			"H007",
			fmt.Sprintf(
				"manual byte-size parsing (suffixChecks=%d, multMap=%v) — use humanize.ParseBytes instead",
				suffixChecks, multMap,
			),
			"Replace with humanize.ParseBytes(s) which handles all SI and IEC suffixes.",
			line, col, filePath, confidence,
		),
	}
}
