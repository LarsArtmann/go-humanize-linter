package humanizelint

import (
	"go/ast"
)

// ---------------------------------------------------------------------------
// Relative-time pattern helpers (H003)
// ---------------------------------------------------------------------------

// hasTimeSinceOrSub reports whether fn calls time.Since or a .Sub method on a
// time value.
func hasTimeSinceOrSub(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// time.Since(x)
		if isPackageCall(call, "time", "Since") {
			hit = true

			return true
		}

		// x.Sub(y)
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "Sub" {
				hit = true
			}
		}

		return true
	})

	return hit
}

// hasTimeThresholdComparison reports whether fn compares an expression against
// a time-package duration constant (time.Minute, time.Hour, time.Second, etc.)
// or a multiplication of one (e.g. 24*time.Hour).
func hasTimeThresholdComparison(fn *ast.FuncDecl) bool {
	timeConsts := map[string]bool{
		"Second": true, "Minute": true, "Hour": true,
		"Day": true, "Week": true, "Month": true, "Year": true,
	}

	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		// Look for time.Hour, time.Minute, etc. in any selector expression.
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "time" {
			if timeConsts[sel.Sel.Name] {
				hit = true
			}
		}

		return true
	})

	return hit
}
