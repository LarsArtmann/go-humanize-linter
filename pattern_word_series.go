package humanizelint

import (
	"go/ast"
	"go/token"
)

// ---------------------------------------------------------------------------
// WordSeries pattern helpers (H012)
// ---------------------------------------------------------------------------

// wordSeriesEvidence holds the H012 signals collected per function.
type wordSeriesEvidence struct {
	// prefixJoin: strings.Join(x[:len(x)-1], ", ") — joining all but the
	// last element, the distinctive WordSeries core
	prefixJoin bool
	// anyCommaJoin: strings.Join(x, ", ") (any form)
	anyCommaJoin bool
	// lastElementIndex: x[len(x)-1] access
	lastElementIndex bool
	// conjunctionLiteral: an "and"/"or" conjunction string literal used in a
	// concatenation or call (" and ", ", and ", " or ", ", or ")
	conjunctionLiteral bool
	// lengthBranches: len(x) == 1 or len(x) == 2 branch handling
	lengthBranches bool
}

// conjunctionLiterals are the separator strings a hand-rolled word series
// uses to introduce the final element.
var conjunctionLiterals = map[string]bool{ //nolint:gochecknoglobals // package-level lookup table
	" and ":     true,
	", and ":    true,
	" and":      true,
	" or ":      true,
	", or ":     true,
	" or":       true,
	"and":       true,
	"or":        true,
	", and":     true,
	", or":      true,
}

// isCommaJoin reports whether call is strings.Join with a ", " (or ",")
// separator and returns the joined expression.
func isCommaJoin(call *ast.CallExpr, aliases map[string]string) bool {
	if !isPackageCall(call, "strings", "Join", aliases) || len(call.Args) < 2 {
		return false
	}

	sep, ok := unquoteString(getBasicLit(call.Args[1]))

	return ok && (sep == ", " || sep == ",")
}

// isLenMinusOneSlice reports whether expr is a slice expression whose high
// bound is len(x)-1 for some x (e.g. parts[:len(parts)-1]).
func isLenMinusOneSlice(expr ast.Expr) bool {
	slice, ok := expr.(*ast.SliceExpr)
	if !ok || slice.High == nil {
		return false
	}

	binary, ok := slice.High.(*ast.BinaryExpr)
	if !ok || binary.Op != token.SUB {
		return false
	}

	call, ok := binary.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}

	if !isBuiltinLen(call) {
		return false
	}

	lit, ok := binary.Y.(*ast.BasicLit)

	return ok && lit.Kind == token.INT && lit.Value == "1"
}

// isBuiltinLen reports whether call is the builtin len(x).
func isBuiltinLen(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)

	return ok && ident.Name == "len"
}

// isLastElementIndex reports whether expr is x[len(x)-1] for some x.
func isLastElementIndex(expr ast.Expr) bool {
	index, ok := expr.(*ast.IndexExpr)
	if !ok {
		return false
	}

	return isLenMinusOneArith(index.Index)
}

// isLenMinusOneArith reports whether expr is len(x)-1.
func isLenMinusOneArith(expr ast.Expr) bool {
	binary, ok := expr.(*ast.BinaryExpr)
	if !ok || binary.Op != token.SUB {
		return false
	}

	call, ok := binary.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}

	if !isBuiltinLen(call) {
		return false
	}

	lit, ok := binary.Y.(*ast.BasicLit)

	return ok && lit.Kind == token.INT && lit.Value == "1"
}

// hasLenBranch reports whether fn contains a `len(x) == 1` or `len(x) == 2`
// comparison (the singular/dual special cases of a word series).
func hasLenBranch(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		binary, ok := n.(*ast.BinaryExpr)
		if !ok || binary.Op != token.EQL {
			return true
		}

		if isLenCallEqInt(binary.X, binary.Y) || isLenCallEqInt(binary.Y, binary.X) {
			hit = true
		}

		return true
	})

	return hit
}

// isLenCallEqInt reports whether one side of an == is len(x) and the other
// is the literal 1 or 2.
func isLenCallEqInt(a, b ast.Expr) bool {
	call, ok := a.(*ast.CallExpr)
	if !ok || !isBuiltinLen(call) {
		return false
	}

	lit, ok := b.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return false
	}

	return lit.Value == "1" || lit.Value == "2"
}

// collectWordSeriesEvidence walks fn and gathers every H012 signal.
func collectWordSeriesEvidence(fn *ast.FuncDecl, aliases map[string]string) wordSeriesEvidence {
	var evidence wordSeriesEvidence

	ast.Inspect(fn, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if isCommaJoin(node, aliases) {
				evidence.anyCommaJoin = true

				if len(node.Args) >= 1 && isLenMinusOneSlice(node.Args[0]) {
					evidence.prefixJoin = true
				}
			}
		case *ast.BasicLit:
			if val, ok := unquoteString(node); ok && conjunctionLiterals[val] {
				evidence.conjunctionLiteral = true
			}
		case *ast.IndexExpr:
			if isLastElementIndex(node) {
				evidence.lastElementIndex = true
			}
		case *ast.BinaryExpr:
			if node.Op == token.ADD {
				if isLastElementIndex(node.X) || isLastElementIndex(node.Y) {
					evidence.lastElementIndex = true
				}
			}
		}

		return true
	})

	evidence.lengthBranches = hasLenBranch(fn)

	return evidence
}
