package humanizelint

import (
	"go/ast"
)

// ---------------------------------------------------------------------------
// Float trailing-zero helpers (H006)
// ---------------------------------------------------------------------------

// hasNestedTrimRight reports whether fn contains
// strings.TrimRight(strings.TrimRight(x, "0"), ".") — the classic Ftoa pattern.
func hasNestedTrimRight(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		outer, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isPackageCall(outer, "strings", "TrimRight") {
			return true
		}

		if len(outer.Args) < 1 {
			return true
		}

		inner, ok := outer.Args[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		if isPackageCall(inner, "strings", "TrimRight") {
			hit = true
		}

		return true
	})

	return hit
}
