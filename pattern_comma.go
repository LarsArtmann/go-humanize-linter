package humanizelint

import (
	"go/ast"
	"go/token"
)

// ---------------------------------------------------------------------------
// Comma-grouping pattern helpers (H002)
// ---------------------------------------------------------------------------

// hasModulo3 reports whether fn contains a `% 3` expression — the classic
// comma-grouping signal.
func hasModulo3(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		be, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if be.Op == token.REM {
			if lit, ok := be.Y.(*ast.BasicLit); ok && lit.Kind == token.INT && lit.Value == "3" {
				hit = true
			}
		}

		return true
	})

	return hit
}

// hasCommaOrSeparator reports whether fn writes a comma or other thousands
// separator character via WriteString, WriteByte, WriteRune, or appears in a
// strings.Join call with ",".
func hasCommaOrSeparator(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// strings.Join(x, ",")
		if isPackageCall(call, "strings", "Join") && len(call.Args) >= 2 {
			if val, ok := unquoteString(getBasicLit(call.Args[1])); ok {
				if val == "," || val == "." || val == " " {
					hit = true
				}
			}
		}

		// builder.WriteString(","), builder.WriteByte(',')
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		switch sel.Sel.Name {
		case "WriteString":
			if len(call.Args) > 0 {
				if val, ok := unquoteString(getBasicLit(call.Args[0])); ok {
					if val == "," || val == "." {
						hit = true
					}
				}
			}

		case "WriteByte", "WriteRune":
			if len(call.Args) > 0 {
				lit := getBasicLit(call.Args[0])
				if lit != nil && (lit.Value == "','" || lit.Value == "'.'") {
					hit = true
				}
			}
		}

		return true
	})

	return hit
}

// hasStepBy3 reports whether fn contains a for-loop incrementing by 3
// (e.g. i += 3 or i = i + 3).
func hasStepBy3(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		forStmt, ok := n.(*ast.ForStmt)
		if !ok {
			return true
		}

		if forStmt.Post == nil {
			return true
		}

		ast.Inspect(forStmt.Post, func(nn ast.Node) bool {
			if be, ok := nn.(*ast.BinaryExpr); ok {
				if be.Op == token.ADD {
					if lit, ok := be.Y.(*ast.BasicLit); ok && lit.Value == "3" {
						hit = true
					}
				}
			}

			if as, ok := nn.(*ast.AssignStmt); ok && as.Tok == token.ADD_ASSIGN {
				for _, v := range as.Rhs {
					if lit, ok := v.(*ast.BasicLit); ok && lit.Value == "3" {
						hit = true
					}
				}
			}

			return true
		})

		return true
	})

	return hit
}

// hasForLoop reports whether fn contains any for-loop.
func hasForLoop(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		if _, ok := n.(*ast.ForStmt); ok {
			hit = true
		}

		return true
	})

	return hit
}

// hasDigitConversion reports whether fn calls strconv.Itoa, strconv.FormatInt,
// strconv.FormatUint, or strconv.FormatFloat — the integer-to-string
// conversions that comma formatters always use.
func hasDigitConversion(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isPackageCall(call, "strconv", "Itoa") ||
			isPackageCall(call, "strconv", "FormatInt") ||
			isPackageCall(call, "strconv", "FormatUint") ||
			isPackageCall(call, "strconv", "FormatFloat") {
			hit = true
		}

		return true
	})

	return hit
}
