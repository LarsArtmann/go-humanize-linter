package humanizelint

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/larsartmann/go-finding"
)

// ---------------------------------------------------------------------------
// String literal helpers
// ---------------------------------------------------------------------------

// unquoteString extracts the Go string value from a *ast.BasicLit, returning
// the unquoted content and true. Returns "", false for non-string literals.
func unquoteString(lit *ast.BasicLit) (string, bool) {
	if lit == nil || lit.Kind != token.STRING {
		return "", false
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return val, true
}

// allStringLiterals collects the values of all string basic literals within a
// function body (including nested call arguments, composite literals, etc.).
func allStringLiterals(fn *ast.FuncDecl) []string {
	var out []string

	ast.Inspect(fn, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok {
			if val, ok := unquoteString(lit); ok {
				out = append(out, val)
			}
		}

		return true
	})

	return out
}

// hasAnyStringLiteral reports whether any string literal in fn contains any of
// the substrings.
func hasAnyStringLiteral(fn *ast.FuncDecl, substrs ...string) bool {
	all := allStringLiterals(fn)

	for _, s := range all {
		for _, sub := range substrs {
			if strings.Contains(s, sub) {
				return true
			}
		}
	}

	return false
}

// ---------------------------------------------------------------------------
// Shared AST utilities
// ---------------------------------------------------------------------------

// normLit strips Go underscore digit separators from a basic literal value.
// e.g. "1_000_000" → "1000000", "1_024" → "1024".
func normLit(val string) string {
	return strings.ReplaceAll(val, "_", "")
}

// funcReturnsString reports whether fn has at least one string-typed return
// value. Used to filter H004 false positives: pluralization functions always
// return strings, while `if x == 1` in error/validation/code-gen paths usually
// returns error, bool, or int.
func funcReturnsString(fn *ast.FuncDecl) bool {
	if fn.Type == nil || fn.Type.Results == nil {
		return false
	}

	for _, field := range fn.Type.Results.List {
		if isStringType(field.Type) {
			return true
		}
	}

	return false
}

// isStringType reports whether an AST type expression denotes string or []string.
func isStringType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name == "string"
	case *ast.ArrayType:
		if ident, ok := t.Elt.(*ast.Ident); ok {
			return ident.Name == "string"
		}
	}

	return false
}

// getBasicLit returns the *ast.BasicLit underlying expr, unwrapping parentheses
// and unary plus/minus operators.
func getBasicLit(expr ast.Expr) *ast.BasicLit {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e
	case *ast.ParenExpr:
		return getBasicLit(e.X)
	case *ast.UnaryExpr:
		return getBasicLit(e.X)
	}

	return nil
}

// isLiteralInt reports whether expr is a basic literal with the given int value.
func isLiteralInt(expr ast.Expr, val int) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return false
	}

	return lit.Value == strconv.Itoa(val)
}

// isPackageCall reports whether call is a function call of the form
// pkg.FuncName (e.g. strings.TrimRight).
func isPackageCall(call *ast.CallExpr, pkg, name string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != name {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == pkg
}

// makeFindingWithConfidence constructs a finding with an explicit confidence
// and suggestion text.
func makeFindingWithConfidence(
	ruleID, message, suggestion string,
	line, col int,
	filePath string,
	conf finding.Confidence,
) finding.Finding {
	b := finding.NewBuilder(
		finding.RuleName(ruleID),
		finding.ToolName("go-humanize-linter"),
		message,
		finding.SeverityWarning,
		finding.Pos(finding.FilePath(filePath), line, col),
	).
		WithCategory(finding.CategoryStyle).
		WithFixStrategy(finding.FixStrategySuggest).
		WithConfidence(conf)

	if suggestion != "" {
		b = b.WithSuggestion(suggestion)
	}

	return b.MustBuild()
}
