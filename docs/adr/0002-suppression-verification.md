# ADR 0002: Suppression Verification as a Post-Detection Pass

**Date:** 2026-08-05
**Status:** Accepted

## Context

The linter's main detection loop (`checkFuncDecls` in `walker.go`) skips
suppressed functions entirely — when a function has a `//nolint:gohumanize`
directive, no detectors run on it and no findings are produced. This means
the linter has no record of whether a suppression directive actually
suppressed a real finding.

Two common problems result:

1. **Stale suppressions** — A `//nolint:gohumanize:H001` directive was added
   months ago. Since then, the code was refactored and the H001 pattern no
   longer exists. The directive is now dead weight, but nobody knows because
   the linter silently skips the function.

2. **Misspelled linter names** — AIs and developers frequently write
   `//nolint:go-humanize-linter/H003` (the module path) instead of
   `//nolint:gohumanize:H003` (the analyzer name). The directive silently
   does nothing because golangci-lint matches on the registered analyzer name
   (`gohumanize`), not the Go module path.

Both problems are invisible: the linter produces no output for these
directives. Users believe their suppressions are working when they are not.

## Decision

Implement suppression verification as a **separate post-detection pass**,
not integrated into the main detection loop.

`VerifySuppressions(dir string, report *finding.Report)` runs after the
main report is built. It:

1. Collects every `//nolint` directive in the directory via
   `collectSuppressions`.
2. Checks each directive for unknown linter names (strings containing
   "humanize" but not exactly "gohumanize").
3. For directives that do target `gohumanize`, checks whether any finding
   in the report was actually produced at the same function position.
4. Produces `H0SUP` findings for directives that are misspelled or stale.

### Why not integrate into the detection loop?

1. **The detection loop skips suppressed functions.** By the time the loop
   reaches a suppressed function, it has already decided to skip it. To
   verify suppression, the loop would need to run detectors on suppressed
   functions anyway, defeating the purpose of suppression (performance).

2. **Position matching is simpler post-hoc.** The report already contains
   all findings with their file/line positions. Matching directives to
   findings is a simple lookup (`findingsByFileLine`). Doing this during the
   detection loop would require threading the directive list through every
   detector call.

3. **Separation of concerns.** Detection ("is this code a humanize
   reimplementation?") and verification ("is this suppression directive
   still needed?") are different questions. Coupling them would make the
   detection loop harder to maintain.

### Why H0SUP as a pseudo-rule?

`H0SUP` is not registered in `AllRules()` or `DefaultRegistry()`. It is a
hardcoded constant in `suppression.go`. This was chosen because:

- H0SUP is a **meta-diagnostic** about directives, not about code patterns.
  It does not fit the RuleFunc model (which takes a directory and returns
  findings by scanning AST).
- It should not be filterable by `--enable`/`--disable` — you either want
  suppression verification or you don't (controlled by `--verify-suppressions`).
- Adding it to `AllRules()` would make it appear in `--rules` output and
  `.golangci.yml` enable/disable lists, which would be confusing.

**Trade-off:** H0SUP cannot be suppressed itself. If the suppression
verification produces a false positive, the user must remove or fix the
offending directive rather than suppressing H0SUP.

## Consequences

- **Positive:** Stale and misspelled suppressions are detected without
  changes to the detection loop or performance impact on normal runs.
- **Positive:** The verification is opt-in (`--verify-suppressions` flag),
  so existing CI pipelines that don't use it see no change.
- **Positive:** `VerifySuppressions` is exported, so the plugin path can
  wire it in later without changes to the core library.
- **Negative:** H0SUP cannot be filtered by `--enable`/`--disable`. If
  needed in the future, it would need to be promoted to a registered rule.
- **Negative:** The post-detection pass walks the directory a second time
  (via `collectSuppressions`). For large codebases, this adds a second AST
  parse of every Go file. This is acceptable because `--verify-suppressions`
  is an opt-in flag, not a default behavior.
