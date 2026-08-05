package humanizelint

import (
	"go/ast"
	"strings"
)

// ---------------------------------------------------------------------------
// Relative-time pattern helpers (H003)
// ---------------------------------------------------------------------------

// hasTimeSinceOrSub reports whether fn calls time.Since or a .Sub method on a
// time value.
func hasTimeSinceOrSub(fn *ast.FuncDecl, aliases map[string]string) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// time.Since(x)
		if isPackageCall(call, "time", "Since", aliases) {
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
// or a multiplication of one (e.g. 24*time.Hour). Supports named aliases
// (tm "time" → tm.Hour) and dot imports (. "time" → bare Hour).
func hasTimeThresholdComparison(fn *ast.FuncDecl, aliases map[string]string) bool {
	timeConsts := map[string]bool{
		"Second": true, "Minute": true, "Hour": true,
		"Day": true, "Week": true, "Month": true, "Year": true,
	}

	dotTime := false
	if aliases != nil {
		if resolved, ok := aliases["."]; ok {
			if resolved == "time" || strings.HasSuffix(resolved, "/time") {
				dotTime = true
			}
		}
	}

	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		// time.Hour, tm.Hour (SelectorExpr form)
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if !timeConsts[sel.Sel.Name] {
				return true
			}

			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "time" {
					hit = true
					return false
				}

				if aliases != nil {
					if resolved, ok := aliases[ident.Name]; ok {
						if resolved == "time" || strings.HasSuffix(resolved, "/time") {
							hit = true
							return false
						}
					}
				}
			}

			return true
		}

		// Dot import: bare Hour (Ident form)
		if dotTime {
			if ident, ok := n.(*ast.Ident); ok {
				if timeConsts[ident.Name] {
					hit = true
					return false
				}
			}
		}

		return true
	})

	return hit
}
