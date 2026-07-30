package humanizelint

import (
	"go/ast"
)

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
