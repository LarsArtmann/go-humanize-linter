package humanizelint

import (
	"go/ast"
	"go/token"
	"regexp"
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
// Byte-unit pattern helpers (H001)
// ---------------------------------------------------------------------------

// byteUnitSet contains known byte-size unit suffixes (case-insensitive). The
// detection logic counts distinct units found in string literals within a
// function.
var byteUnitSet = map[string]bool{
	"KiB": true, "MiB": true, "GiB": true, "TiB": true, "PiB": true, "EiB": true,
	"KB": true, "MB": true, "GB": true, "TB": true, "PB": true, "EB": true,
	"kB": true,
}

// siPrefixStrings are the exact string literals used as byte-prefix arrays.
// "KMGTPE" is the canonical 6-character prefix string; "kMGTPE" is the
// lowercase-k variant.
var siPrefixStrings = map[string]bool{
	"KMGTPE": true, "kMGTPE": true,
	"KMGTPEZY": true, "kMGTPEZY": true,
}

// byteUnitRegex matches known byte-size unit suffixes as whole words.
// The \b boundaries prevent false positives like "MEDIUMBLOB" → "MB".
var byteUnitRegex = regexp.MustCompile(`(?i)\b(KiB|MiB|GiB|TiB|PiB|EiB|KB|MB|GB|TB|PB|EB|kB)\b`)

// countByteUnits scans all string literals in fn and returns the set of
// distinct byte-unit suffixes found (using word-boundary matching to avoid
// false positives like "MEDIUMBLOB" containing "MB").
func countByteUnits(fn *ast.FuncDecl) map[string]bool {
	found := map[string]bool{}

	for _, s := range allStringLiterals(fn) {
		for _, m := range byteUnitRegex.FindAllString(s, -1) {
			found[strings.ToUpper(m)] = true
		}
	}

	return found
}

// hasKMGTPEIndex reports whether fn contains an index expression on a string
// literal like "KMGTPE" — the classic byte-prefix trick.
func hasKMGTPEIndex(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		idx, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		if val, ok := unquoteString(getBasicLit(idx.X)); ok {
			if siPrefixStrings[val] {
				hit = true
			}
		}

		return true
	})

	return hit
}

// hasByteUnitSlice reports whether fn contains a []string composite literal
// whose elements include 2+ known byte-unit strings.
func hasByteUnitSlice(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		count := 0

		for _, el := range cl.Elts {
			if val, ok := unquoteString(getBasicLit(el)); ok {
				if byteUnitSet[val] {
					count++
				}
			}
		}

		if count >= 2 {
			hit = true
		}

		return true
	})

	return hit
}

// hasDivisionByPowerOf1024 reports whether fn contains a binary division whose
// divisor is the literal 1024, 1048576 (1024²), 1073741824 (1024³), or any
// expression multiplying by 1024.
func hasDivisionByPowerOf1024(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		be, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if be.Op == token.QUO && isLiteral1024ish(be.Y) {
			hit = true
		}

		return true
	})

	return hit
}

// isLiteral1024ish reports whether expr is the literal 1024, 1048576,
// 1073741824, or a multiplication chain involving 1024 (e.g. 1024*1024).
func isLiteral1024ish(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch normLit(e.Value) {
		case "1024", "1048576", "1073741824", "1099511627776":
			return true
		}

	case *ast.BinaryExpr:
		if e.Op == token.MUL {
			return isLiteral1024ish(e.X) || isLiteral1024ish(e.Y)
		}
	}

	return false
}

// hasConst1024 reports whether fn or its enclosing file defines a constant
// whose value is 1024 or a power of 1024. This catches both in-function
// constants (`const unit = 1024`) and file-level constants (`const bytesKB = 1024`).
func hasConst1024(file *ast.File, fn *ast.FuncDecl) bool {
	hit := false

	// Check constants inside the function body.
	ast.Inspect(fn, func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			return true
		}

		for _, spec := range gd.Specs {
			if specHas1024Value(spec) {
				hit = true
			}
		}

		return true
	})

	// Check file-level constant declarations.
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}

		for _, spec := range gd.Specs {
			if specHas1024Value(spec) {
				hit = true
			}
		}
	}

	return hit
}

// specHas1024Value reports whether a ValueSpec contains a literal 1024-ish value.
func specHas1024Value(spec ast.Spec) bool {
	vs, ok := spec.(*ast.ValueSpec)
	if !ok {
		return false
	}

	for _, val := range vs.Values {
		if isLiteral1024ish(val) {
			return true
		}
	}

	return false
}

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

// ---------------------------------------------------------------------------
// Pluralization pattern helpers (H004)
// ---------------------------------------------------------------------------

// hasEqualsOneBranch reports whether fn contains an `if x == 1` or `if x != 1`
// conditional where x is a simple identifier (not a function call like len(x)),
// AND the branch body contains at least one string literal or string-producing
// statement (return/assignment with a string expression). The string-in-branch
// requirement eliminates false positives like `if n == 1 { return sorted[0] }`
// (percentile calculation), exit-code extraction, and code generation — none of
// which produce different strings based on the count.
func hasEqualsOneBranch(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		if hit {
			return false
		}

		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		be, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if be.Op != token.EQL && be.Op != token.NEQ {
			return true
		}

		var otherSide ast.Expr

		if isLiteralInt(be.Y, 1) {
			otherSide = be.X
		} else if isLiteralInt(be.X, 1) {
			otherSide = be.Y
		} else {
			return true
		}

		if _, ok := otherSide.(*ast.Ident); !ok {
			return true
		}

		if branchContainsString(ifStmt.Body) {
			hit = true
		}

		return true
	})

	return hit
}

// branchContainsString reports whether a block statement contains any string
// literal or string concatenation (BinaryExpr with +). This confirms that an
// `if x == 1` branch is actually producing different text output.
func branchContainsString(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}

	found := false

	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}

		switch e := n.(type) {
		case *ast.BasicLit:
			if _, ok := unquoteString(e); ok {
				found = true
			}

		case *ast.BinaryExpr:
			if e.Op == token.ADD {
				if exprIsStringy(e.X) || exprIsStringy(e.Y) {
					found = true
				}
			}
		}

		return true
	})

	return found
}

// exprIsStringy is a lightweight heuristic: reports whether expr is a basic
// string literal or an identifier (which might be a string variable).
func exprIsStringy(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		_, ok := unquoteString(e)

		return ok
	case *ast.Ident:
		return true
	}

	return false
}

// hasPluralNamedParams reports whether fn has parameters whose names include
// "singular" and "plural".
func hasPluralNamedParams(fn *ast.FuncDecl) bool {
	var hasSingular, hasPlural bool

	if fn.Type != nil && fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			for _, name := range field.Names {
				lower := strings.ToLower(name.Name)
				if strings.Contains(lower, "singular") {
					hasSingular = true
				}

				if strings.Contains(lower, "plural") {
					hasPlural = true
				}
			}
		}
	}

	return hasSingular && hasPlural
}

// ---------------------------------------------------------------------------
// SI / K-M pattern helpers (H005)
// ---------------------------------------------------------------------------

// hasDivisionBy1000 reports whether fn contains a division by the literal
// 1000, 1000000, 1000000000, or the scientific-notation equivalents
// (1e3, 1e6, 1e9, 1e12).
func hasDivisionBy1000(fn *ast.FuncDecl) bool {
	powersOf1000 := map[string]bool{
		"1000": true, "1000000": true, "1000000000": true, "1000000000000": true,
		"1e3": true, "1e6": true, "1e9": true, "1e12": true,
	}

	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		be, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if be.Op == token.QUO {
			if lit, ok := be.Y.(*ast.BasicLit); ok && powersOf1000[normLit(lit.Value)] {
				hit = true
			}
		}

		return true
	})

	return hit
}

// hasKMSuffix reports whether fn contains SI-prefix suffixes — either as
// standalone string literals ("K", "M") or embedded in format strings
// ("%.2fK", "%.1fM"). Excludes byte-unit strings to avoid overlap with H001.
func hasKMSuffix(fn *ast.FuncDecl) bool {
	for _, s := range allStringLiterals(fn) {
		// Standalone single-letter SI prefixes.
		switch s {
		case "K", "M", "G", "T", "k", "m":
			return true
		}

		// Embedded suffixes like "%.2fK" or "%.1fM" — the last char is an SI
		// prefix and the string is NOT a byte-unit pattern (e.g., "%.1f MB").
		if len(s) >= 2 && !byteUnitRegex.MatchString(s) {
			switch s[len(s)-1] {
			case 'K', 'M', 'G', 'T', 'k':
				return true
			}
		}
	}

	return false
}

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

// ---------------------------------------------------------------------------
// ParseBytes pattern helpers (H007)
// ---------------------------------------------------------------------------

// hasByteUnitSuffixChecks reports whether fn contains 2+ calls to
// strings.HasSuffix, strings.CutSuffix, or strings.TrimSuffix where the suffix
// argument is a known byte-unit string (KB, MB, GB, KiB, MiB, etc.).
func hasByteUnitSuffixChecks(fn *ast.FuncDecl) int {
	count := 0

	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isPackageCall(call, "strings", "HasSuffix") &&
			!isPackageCall(call, "strings", "CutSuffix") &&
			!isPackageCall(call, "strings", "TrimSuffix") {
			return true
		}

		if len(call.Args) < 2 {
			return true
		}

		if val, ok := unquoteString(getBasicLit(call.Args[1])); ok {
			if byteUnitSet[val] {
				count++
			}
		}

		return true
	})

	return count
}

// hasByteUnitMultiplierMap reports whether fn contains a map literal with
// string keys that are byte-unit names (KB, MB, GiB, etc.) and integer values.
func hasByteUnitMultiplierMap(fn *ast.FuncDecl) bool {
	hit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		count := 0

		for _, el := range cl.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			if val, ok := unquoteString(getBasicLit(kv.Key)); ok {
				if byteUnitSet[val] {
					count++
				}
			}
		}

		if count >= 2 {
			hit = true
		}

		return true
	})

	return hit
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
