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

// ---------------------------------------------------------------------------
// Suppression directives
// ---------------------------------------------------------------------------

// nolintLinterName is the analyzer name users write in //nolint directives. It
// matches plugin.Analyzer.Name ("gohumanize").
const nolintLinterName = "gohumanize"

// nolintAllMarker is the sentinel token meaning "suppress every linter" in a
// //nolint directive (e.g. "//nolint" or "//nolint:all").
const nolintAllMarker = "all"

// lintIgnorePrefix is the Go-style suppression directive prefix recognised in
// addition to //nolint. Format:
//
//	//lint:ignore gohumanize          (suppress every H-rule)
//	//lint:ignore gohumanize:H001     (only H001)
//	//lint:ignore gohumanize reason   (suppress every H-rule, with reason)
const lintIgnorePrefix = "//lint:ignore"

// hasNoLintDirective reports whether fn carries a //nolint directive that
// suppresses this linter. A directive counts if it appears in the function's
// doc comment, as a trailing comment on the func's own line, or in any comment
// group positioned on the line immediately before the declaration.
//
// Recognised forms:
//
//	//nolint                     (suppresses everything)
//	//nolint:all                 (suppresses everything)
//	//nolint:gohumanize          (suppresses only this linter)
//	//nolint:gohumanize,other    (comma-separated list)
//	//nolint:gohumanize:H001     (scoped — suppresses only H001 in this linter)
func hasNoLintDirective(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) bool {
	if file == nil || fn == nil {
		return false
	}

	fnLine := fset.Position(fn.Pos()).Line

	for _, group := range file.Comments {
		isDoc := group == fn.Doc

		for _, comment := range group.List {
			suppressed := suppressedRules(comment.Text)
			if suppressed == nil {
				continue
			}

			if !isSuppressedAll(suppressed) {
				continue
			}

			cmtLine := fset.Position(comment.Pos()).Line
			if isDoc || cmtLine == fnLine || cmtLine == fnLine-1 {
				return true
			}
		}
	}

	return false
}

// suppressedRules parses a suppression comment and returns the linter/rule
// names it lists. Returns nil if the comment is not a suppression directive.
//
// Recognised syntax (two flavours):
//
//	//nolint                            → ["all"]
//	//nolint:all                        → ["all"]
//	//nolint:gohumanize                 → ["gohumanize"]
//	//nolint:gohumanize:H001            → ["gohumanize", "H001"]
//	//nolint:gohumanize:H001,H002       → ["gohumanize", "H001", "H002"]
//	//nolint:gohumanize,other           → ["gohumanize", "other"]
//	//lint:ignore gohumanize            → ["gohumanize"]   (Go-style)
//	//lint:ignore gohumanize:H001       → ["gohumanize", "H001"]
//	//lint:ignore gohumanize reason     → ["gohumanize"]   (reason after linter)
//
// The "all" token means "suppress every linter". When the first token after
// "nolint:" or "lint:ignore" is "gohumanize" (our linter name), any subsequent
// colon-prefixed tokens like ":H001" are sub-scopes within that linter.
func suppressedRules(commentText string) []string {
	body := strings.TrimSpace(strings.TrimPrefix(commentText, "//"))

	// Two supported directive prefixes.
	var rest string

	switch {
	case strings.HasPrefix(body, "nolint"):
		rest = strings.TrimSpace(strings.TrimPrefix(body, "nolint"))
	case strings.HasPrefix(body, "lint:ignore"):
		rest = strings.TrimSpace(strings.TrimPrefix(body, "lint:ignore"))
	default:
		return nil
	}

	// Bare "//nolint" or "//lint:ignore" (with no linter name) suppresses
	// everything only in the nolint flavour; lint:ignore requires a linter.
	if rest == "" {
		return []string{nolintAllMarker}
	}

	if !strings.HasPrefix(rest, ":") {
		// No colon — first whitespace-delimited token is the linter name (or
		// "all"). For lint:ignore, "all" means "suppress every linter".
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			return nil
		}

		body2 := fields[0]
		if body2 == "" {
			return nil
		}

		return splitAndExpand(body2)
	}

	body2 := noLintList(rest)
	if body2 == "" {
		return nil
	}

	return splitAndExpand(body2)
}

// splitAndExpand splits a comma-separated list of linter names and colon-scoped
// rule IDs into a flat []string. Examples:
//
//	"gohumanize"             → ["gohumanize"]
//	"gohumanize:H001"        → ["gohumanize", "H001"]
//	"H001"                   → ["H001"]
//	"gohumanize:H001,H002"   → ["gohumanize", "H001", "H002"]
func splitAndExpand(body2 string) []string {
	out := []string{}

	for name := range strings.SplitSeq(body2, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		// Split colon-scoped form: "gohumanize:H001" → ["gohumanize", "H001"].
		for part := range strings.SplitSeq(name, ":") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}

	return out
}

// isSuppressedAll reports whether the suppression list (from suppressedRules)
// matches every rule. "all" suppresses everything. If "gohumanize" appears as
// a top-level entry without a sub-rule, every H-rule is suppressed.
func isSuppressedAll(suppressions []string) bool {
	for _, s := range suppressions {
		if s == nolintAllMarker {
			return true
		}

		if s == nolintLinterName {
			return true
		}
	}

	return false
}

// isSuppressedRule reports whether suppression list suppresses the specific
// ruleID (e.g. "H001"). Recognises:
//
//	//nolint:all                            (suppresses everything)
//	//nolint:gohumanize                     (suppresses every H-rule)
//	//nolint:gohumanize:H001                (only H001)
//	//nolint:gohumanize:H001,H002           (H001 and H002)
//	//nolint:H001                           (only H001, unscoped)
//
// "gohumanize" alone is treated as "all H-rules" (suppresses every rule ID)
// so existing //nolint:gohumanize behaviour is preserved. When scoped
// sub-rules appear alongside "gohumanize", only those specific rule IDs are
// suppressed.
func isSuppressedRule(suppressions []string, ruleID string) bool {
	hasAll := false
	hasLinter := false
	scopedRules := map[string]bool{}

	for _, s := range suppressions {
		switch s {
		case nolintAllMarker:
			hasAll = true
		case nolintLinterName:
			hasLinter = true
		default:
			// Sub-rule token (e.g. "H001"). Only meaningful when "gohumanize"
			// is also present as a scope marker, but we record it either way.
			scopedRules[s] = true
		}
	}

	if hasAll {
		return true
	}

	if hasLinter {
		// Scoped sub-rules: only suppress the rules explicitly named.
		if len(scopedRules) > 0 {
			return scopedRules[ruleID]
		}

		// Bare "gohumanize" without sub-rules → suppress every H-rule.
		return true
	}

	// Unscoped sub-rule (e.g. //nolint:H001): suppress only that rule.
	return scopedRules[ruleID]
}

// noLintList extracts the comma-separated linter name list from the text that
// follows "nolint:". Trailing explanation text (e.g. the second "// reason"
// after the names) is ignored. Returns the empty string for a bare //nolint.
func noLintList(rest string) string {
	fields := strings.Fields(strings.TrimPrefix(rest, ":"))
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}
