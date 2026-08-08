# ADR 0006: Per-Statement Suppression via Line-Range Matching

**Date:** 2026-08-08
**Status:** Accepted

## Context

go-humanize-linter emits exactly one finding per function, positioned at the
function declaration (`fn.Pos()`). Suppression via `//nolint:gohumanize`
originally only worked when the directive appeared on the function declaration
line or its doc comment.

Developers and AI tools naturally try to place `//nolint` on the specific
statement that triggers the finding, inside the function body:

```go
func formatBytes(bytes int64) string {
    // ...
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp]) //nolint:gohumanize
}
```

This did not work — the directive was ignored, and the finding still fired.

## Decision

**Approach B: Line-range matching.** Widen the comment-association check to
accept any `//nolint` directive placed between the function's opening brace
(`fn.Body.Lbrace`) and closing brace (`fn.Body.Rbrace`), inclusive.

The shared helper `commentAssociatedWithFunc` already determines whether a
comment is associated with a function. We extend it with one additional
condition: if the comment's line falls within the `[Lbrace, Rbrace]` range,
it counts as associated.

This approach requires changing exactly one function and zero detector files.
All three suppression paths (CLI `checkFuncDecls`, plugin `HumanizeDetector.Run`,
verification `extractFunctionSuppressions`) already call
`commentAssociatedWithFunc` and inherit the new behavior automatically.

### Rejected Alternative: Approach A (per-statement token.Pos)

Approach A would emit findings at the specific statement position
(e.g., the `return` statement's `token.Pos`) instead of the function
declaration. This would allow `//nolint` to be matched precisely to the
statement that triggers it, matching golangci-lint's native behavior.

However, this requires:

1. Updating all 9 detectors to return per-statement positions instead of
   `posOf(fset, fn.Pos())`.
2. Updating the suppression matching logic to compare comment lines against
   finding lines (not function declaration lines).
3. Updating `directiveMatchesFinding` to key on finding line, not
   `FunctionLine`.

The risk/effort is disproportionate to the benefit. All detectors currently
emit one finding per function; per-statement precision provides no additional
user value until a detector emits multiple findings within a single function.

### Why not wait for Approach A?

Approach B is not a stopgap — it is the correct design for the current
detection model. Because every finding is function-scoped (one per function),
function-scoped suppression is the right granularity. The in-body directive
is an ergonomic improvement (place the comment where the code is), not a
precision improvement.

If a future detector emits per-statement findings, Approach A becomes
necessary and can be layered on top without breaking Approach B users.

## Consequences

- **Positive:** Developers can place `//nolint:gohumanize` on the specific
  line that triggers the finding, improving readability and intent clarity.
- **Positive:** Minimal code change — one helper function, inherited by all
  three suppression paths.
- **Positive:** No detector changes required.
- **Negative:** A directive anywhere in the function body suppresses the
  entire function's finding. This is acceptable because there is only one
  finding per function today.
- **Negative:** `directiveMatchesFinding` in `suppression.go` still keys on
  `FunctionLine` (the function declaration line), not the in-body comment
  line. This works correctly because findings are emitted at `fn.Pos()`.
  If detectors migrate to per-statement positions, this lookup must be
  updated.
