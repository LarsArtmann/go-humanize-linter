// pattern_commaf.go — H009 manual-commaf detection.
//
// Commaf detection: float formatting + manual separator grouping. Triggered
// when a function calls fmt.Sprintf("%.Nf", x) and writes the result to a
// string with manual comma/separator insertion, OR uses strconv.FormatFloat
// inside a comma-grouping loop.

package humanizelint

import (
	"go/ast"
	"go/token"
)

// hasCommafPattern reports whether fn contains the "%.Nf + manual comma
// grouping" anti-pattern. The signature is:
//
//  1. A call to fmt.Sprintf with a "%.Nf" format string (or
//     strconv.FormatFloat with a non-zero precision), AND
//  2. A manual separator group loop in the same function.

func hasCommafPattern(fn *ast.FuncDecl) bool {
	hasFloatFormat := false
	hasSeparatorLoop := false

	ast.Inspect(fn, func(n ast.Node) bool {
		if hasFloatFormat && hasSeparatorLoop {
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isPackageCall(call, "fmt", "Sprintf") {
			if hasPercentNF(call) {
				hasFloatFormat = true
			}
		}

		if isPackageCall(call, "strconv", "FormatFloat") && len(call.Args) >= 3 {
			// FormatFloat(x, 'f', prec, bits) — third arg is precision.
			lit := getBasicLit(call.Args[2])
			if lit != nil && lit.Kind == token.INT && lit.Value != "0" && lit.Value != "-1" {
				hasFloatFormat = true
			}
		}

		if isCommaSeparatorCall(call) {
			hasSeparatorLoop = true
		}

		return true
	})

	return hasFloatFormat && hasSeparatorLoop
}

// hasPercentNF reports whether a fmt.Sprintf call uses a "%.Nf" format
// (e.g. "%.1f", "%.2f") — a non-zero precision float.
func hasPercentNF(call *ast.CallExpr) bool {
	if len(call.Args) < 1 {
		return false
	}

	lit := getBasicLit(call.Args[0])
	if lit == nil || lit.Kind != token.STRING {
		return false
	}

	val, ok := unquoteString(lit)
	if !ok {
		return false
	}

	// Walk the format string looking for %.<digit>f verbs.
	for i := 0; i+2 < len(val); i++ {
		if val[i] != '%' {
			continue
		}

		if val[i+1] == '.' {
			// %.<digit>f
			if i+2 < len(val) && val[i+2] >= '0' && val[i+2] <= '9' {
				for j := i + 3; j < len(val); j++ {
					if val[j] == 'f' {
						return true
					}

					if val[j] != val[i+2] {
						break
					}
				}
			}
		} else if val[i+1] >= '0' && val[i+1] <= '9' {
			// %<digit>f (no dot)
			for j := i + 2; j < len(val); j++ {
				if val[j] == 'f' {
					return true
				}

				if val[j] != val[i+1] {
					break
				}
			}
		}
	}

	return false
}
