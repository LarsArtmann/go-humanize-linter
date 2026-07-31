# ADR 0001: Import-Alias-Aware Package Call Detection

**Date:** 2026-07-30
**Status:** Accepted

## Context

The linter's `isPackageCall` helper checks whether a call expression matches
`pkg.FuncName` (e.g. `strings.HasSuffix`). The original implementation did a
pure syntactic comparison: `ident.Name == pkg`. This fails when the source
code uses an import alias:

```go
import str "strings"
str.HasSuffix(s, "KB")  // NOT detected — ident.Name is "str", not "strings"
```

Aliased imports are uncommon but not rare. Without alias resolution, the
linter produces false negatives on any codebase that aliases `strings`,
`strconv`, `fmt`, or `time`.

## Decision

Use **syntactic import-alias resolution** via `buildImportAliases(file)`,
not full `go/types` type checking.

`buildImportAliases` scans the file's `ast.File.Imports` and builds a
`map[string]string` from local names to import paths. For non-aliased
imports, the local name is the package's last path segment. For aliased
imports, the local name is the alias.

`isPackageCall` now accepts a variadic `aliases ...map[string]string`
parameter. When provided, it checks whether the selector's package name
resolves to the expected package via the alias map.

### Why not full go/types?

1. **CLI path lacks type info.** The standalone CLI uses `go/parser` to walk
   directories. It does not run the type checker, so `pass.TypesInfo` is
   unavailable. Using go/types would require either type-checking every file
   (slow) or maintaining two code paths (complex).

2. **Sufficient precision.** Import aliases are the only case where the
   package name in the source differs from the expected name. Dot imports
   (`. "strings"`) are excluded by the `*ast.SelectorExpr` check — a dot
   import produces `HasSuffix(s, "KB")`, not `strings.HasSuffix(...)`, so it
   wouldn't match `isPackageCall` regardless. This is an acceptable gap
   since dot imports for `strings`/`strconv`/`fmt` are extremely rare.

3. **No new dependencies.** Syntactic resolution uses only `go/ast`, which
   is already imported. Full go/types would add complexity to the CLI walker.

## Consequences

- **Positive:** All pattern helpers now detect aliased imports with zero
  new dependencies and no CLI walker changes.
- **Positive:** The `aliases` parameter is variadic, so existing callers
  that don't pass it still compile (backward compatible).
- **Negative:** Dot imports (`. "strings"`) are not detected. This is a
  known limitation, acceptable because dot imports of stdlib packages are
  near-nonexistent in real codebases.
- **Negative:** Every pattern helper that calls `isPackageCall` must accept
  and forward the `aliases` parameter. This adds one parameter to 7 helper
  functions.
