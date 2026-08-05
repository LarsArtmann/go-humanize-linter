# ADR 0005: H009/H002 Disambiguation

**Date:** 2026-08-05
**Status:** Accepted

## Context

H002 (`manual-comma-format`) and H009 (`manual-commaf`) both detect manual
thousands-separator insertion. They differ in specificity:

- **H002** triggers on grouping signals (modulo-3, step-by-3, for-loop +
  digit conversion) combined with comma/separator writing. It suggests
  `humanize.Comma` (integer) or `humanize.Commaf` (float).
- **H009** triggers on the same comma/separator writing **plus** a `%.Nf`
  `fmt.Sprintf` call — a strictly more specific signal indicating
  float-with-comma formatting. It suggests `humanize.Commaf` directly.

A function that formats a float with `%.2f` and then manually inserts commas
matches **both** patterns. Without disambiguation, the user sees two findings
on the same function:

```
H002: manual comma formatting (mod3=true, step3=true) — use humanize.Comma instead
H009: manual float-with-comma formatting (%.Nf + manual separator) — use humanize.Commaf instead
```

This is noise. The two findings point at the same code, suggest overlapping
fixes, and force the user to mentally deduplicate. Worse, the H002 suggestion
(`humanize.Comma(int64(n))`) is wrong for a float — it would truncate.

## Decision

**H002 suppresses itself when H009's pattern matches on the same function.**
The suppression is unidirectional: H009 always fires when its pattern
matches; H002 yields when `hasCommafPattern(fn, aliases)` returns true.

```go
func detectCommaFormat(...) []finding.Finding {
    // Suppress H002 when H009 (commaf) would fire on the same function.
    if hasCommafPattern(fn, aliases) {
        return nil
    }
    // ... H002 detection logic ...
}
```

### Why suppress H002 and not H009?

1. **H009 is more specific.** It requires an additional signal (`%.Nf`) that
   H002 does not. The more specific diagnosis should win.
2. **H009's suggestion is more accurate.** `humanize.Commaf(f)` is the correct
   replacement for float-with-comma code. `humanize.Comma(int64(n))` is wrong
   for floats.
3. **H009 carries more information.** The `%.Nf` signal tells the user (and
   the linter author) that this is definitively a float formatter, not an
   integer formatter. H002's message ("mod3=true, step3=true") does not
   convey this.

### Why unidirectional?

If both rules suppressed each other on overlap, no finding would be produced
at all — the user gets zero guidance on a real reimplementation. Unidirectional
suppression ensures exactly one finding (the better one) fires.

### Why not merge H002 and H009 into one rule?

The two rules have distinct positive test fixtures (`h002_comma_*` for
integer-only formatting, `h009_commaf` for float-with-comma). Merging would
conflate two detection patterns and produce a less precise message for the
integer-only case. Keeping them separate preserves diagnostic precision:

- Integer-only comma formatting → H002 (`humanize.Comma`)
- Float-with-comma formatting → H009 (`humanize.Commaf`)

### Implementation note

The suppression is a **direct pattern check**, not a cross-rule dependency.
`detectCommaFormat` calls `hasCommafPattern(fn, aliases)` directly rather
than inspecting H009's findings. This keeps the rules independently runnable
— `RuleComma()` can be enabled without `RuleCommaf()` and vice versa. The
tradeoff is that `hasCommafPattern` is evaluated twice when both rules are
enabled (once by H009, once by H002's suppression check). This is negligible:
the function is a single AST walk, and rules run sequentially per file.

## Consequences

- **Positive:** Users see exactly one finding per function for comma-related
  formatting, with the more accurate suggestion.
- **Positive:** Rules remain independently enable/disable-able. Disabling
  H009 does not break H002's suppression check — `hasCommafPattern` is a
  pure pattern function, not a dependency on H009's rule execution.
- **Negative:** A small redundancy: `hasCommafPattern` may be evaluated
  twice per function when both rules are enabled. Acceptable given the AST
  walk is cheap.
- **Negative:** If a future rule (e.g., H010 "currency formatting") also
  overlaps with H002, the suppression logic would need extending. The current
  pattern (direct check in `detectCommaFormat`) scales to a second check but
  does not generalize to N-way overlap. A priority-ordered suppression
  mechanism would be needed if overlaps become common.
