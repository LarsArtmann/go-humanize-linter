// pattern_ordinal.go — H008 manual-ordinal detection.
//
// Ordinal detection: switch n%10 returning "st" / "nd" / "rd" / "th".
//
// Strategy: look for a switch on n%10 or n%100 with each case returning a
// different 2-character ordinal suffix string.

package humanizelint

import (
	"go/ast"
	"go/token"
)

// ordinalSuffixes is the set of valid English ordinal suffixes. The
// detector flags a function that returns ALL FOUR from a switch.
var ordinalSuffixes = map[string]bool{ //nolint:gochecknoglobals // package-level lookup
	"st": true,
	"nd": true,
	"rd": true,
	"th": true,
}

// hasOrdinalSwitch reports whether fn contains a switch statement on n%10 or
// n%100 with case branches returning at least 3 of the 4 ordinal suffixes
// ("st", "nd", "rd", "th"). Multi-signal heuristic — exact-match hand-rolled
// ordinal code is uncommon; usually 3+ cases is enough to suspect a
// reimplementation.
func hasOrdinalSwitch(fn *ast.FuncDecl) bool {
	found := map[string]bool{}

	ast.Inspect(fn, func(n ast.Node) bool {
		switchStmt, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}

		// Condition must be `n % 10` or `n % 100` (mod expression).
		be, ok := switchStmt.Tag.(*ast.BinaryExpr)
		if !ok || be.Op != token.REM {
			return true
		}

		if isLiteralInt(be.Y, 10) || isLiteralInt(be.Y, 100) {
			// Walk case branches.
			for _, c := range switchStmt.Body.List {
				caseClause, ok := c.(*ast.CaseClause)
				if !ok {
					continue
				}

				walkOrdinalReturns(caseClause.Body, found)
			}
		}

		return true
	})

	// Need at least 3 distinct ordinal suffixes.
	return len(found) >= 3
}

// walkOrdinalReturns inspects a case body for returned string literals that
// look like ordinal suffixes. Records any suffix found in `found`.
func walkOrdinalReturns(body []ast.Stmt, found map[string]bool) {
	for _, stmt := range body {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}

		for _, expr := range ret.Results {
			lit, ok := expr.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}

			val, ok := unquoteString(lit)
			if !ok {
				continue
			}

			if ordinalSuffixes[val] {
				found[val] = true
			}
		}
	}
}