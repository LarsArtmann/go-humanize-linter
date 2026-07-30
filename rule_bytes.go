package humanizelint

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-linter-sdk"
)

// RuleBytes (H001) detects manual byte-size formatting that should use
// humanize.Bytes or humanize.IBytes.
//
// Triggers when a function contains any of these signal clusters:
//
//   - "KMGTPE"[exp] index expression (the classic prefix trick)
//   - A []string literal with 2+ byte-unit strings (KB, MB, KiB, MiB, etc.)
//   - 2+ distinct byte-unit strings in format strings + division by a power of 1024
//   - 3+ distinct byte-unit strings (very strong signal even without explicit division)
func RuleBytes() linter.RuleFunc {
	return linter.RuleFunc{
		Meta: linter.RuleMeta{
			ID:          "H001",
			Name:        "manual-bytes-format",
			Description: "Manual byte-size formatting — use humanize.Bytes or humanize.IBytes instead of dividing by 1024 and formatting unit strings",
			Cat:         linter.CategoryStyle,
			Sev:         finding.SeverityWarning,
		},
		Run: func(_ context.Context, dir string) ([]finding.Finding, error) {
			return checkFuncDecls(dir, detectBytesFormat)
		},
	}
}

func detectBytesFormat(
	fset *token.FileSet,
	file *ast.File,
	fn *ast.FuncDecl,
	filePath string,
) []finding.Finding {
	kmgtp := hasKMGTPEIndex(fn)
	unitSlice := hasByteUnitSlice(fn)
	units := countByteUnits(fn)
	div1024 := hasDivisionByPowerOf1024(fn) || hasConst1024(file, fn)
	unitCount := len(units)

	var confidence finding.Confidence

	switch {
	case kmgtp:
		confidence = finding.ConfidenceFull
	case unitSlice:
		confidence = finding.ConfidenceFull
	case unitCount >= 3:
		confidence = finding.ConfidenceHigh
	case unitCount >= 2 && div1024:
		confidence = finding.ConfidenceHigh
	case unitCount >= 2:
		confidence = finding.ConfidenceMedium
	default:
		return nil
	}

	line, col := posOf(fset, fn.Pos())

	var msg string
	switch {
	case kmgtp:
		msg = "manual byte-size formatting (KMGTPE index trick) — use humanize.Bytes or humanize.IBytes instead"
	case unitSlice:
		msg = fmt.Sprintf(
			"manual byte-size formatting (unit string slice, %d unit strings) — use humanize.Bytes or humanize.IBytes instead",
			unitCount,
		)
	case unitCount >= 3:
		msg = fmt.Sprintf(
			"manual byte-size formatting (%d unit strings, div1024=%v) — use humanize.Bytes or humanize.IBytes instead",
			unitCount,
			div1024,
		)
	case unitCount >= 2 && div1024:
		msg = fmt.Sprintf(
			"manual byte-size formatting (%d unit strings, div1024=%v) — use humanize.Bytes or humanize.IBytes instead",
			unitCount,
			div1024,
		)
	case unitCount >= 2:
		msg = fmt.Sprintf(
			"manual byte-size formatting (%d unit strings) — use humanize.Bytes or humanize.IBytes instead",
			unitCount,
		)
	default:
		return nil
	}

	return []finding.Finding{
		makeFindingWithConfidence(
			"H001",
			msg,
			"Replace with humanize.Bytes(uint64(n)) for SI (KB/MB) or humanize.IBytes(uint64(n)) for IEC (KiB/MiB).",
			line,
			col,
			filePath,
			confidence,
		),
	}
}
