# ADR 0007: Plugin-path H0SUP findings are re-anchored to file position

- Status: Accepted
- Date: 2026-09-02 (decision); recorded as ADR 2026-09-26 during the docs-health pass
- Deciders: go-humanize-linter maintainers

## Context

`--verify-suppressions` produces diagnostics with pseudo-rule ID `H0SUP` for
stale or misspelled `//nolint` directives (ADR 0002). A stale-directive H0SUP
finding is anchored at the `//nolint:gohumanize` directive itself.

golangci-lint's nolint processor drops any issue reported by `gohumanize`
whose position falls under a `//nolint:gohumanize` directive. A stale-directive
H0SUP sits exactly under such a directive — the diagnostic suppressed itself
(self-referential suppression), so plugin users never saw stale-directive
reports. The CLI path has no nolint processor and was unaffected.

## Decision

In the plugin path (`plugin/plugin.go`), H0SUP findings are re-anchored to the
file position (1:1) and the directive coordinates are embedded in the message
text, e.g. `... (directive at main.go:3)`. The CLI path keeps the precise
directive position.

## Consequences

- Plugin users now see stale-directive diagnostics. Verified end-to-end via
  `TestCustomGCLIntegration/verify_suppressions` against a real `custom-gcl`
  binary.
- CLI and plugin report H0SUP at DIFFERENT positions by design. Do not
  "unify" the positions without re-checking that integration test — the
  divergence is what defeats golangci-lint's nolint filter.
- H0SUP also bypasses `minConfidence` filtering in the plugin path, so
  `minConfidence: "full"` cannot hide it (these findings are
  `ConfidenceHigh`).

## Alternatives considered

- Reporting H0SUP at the function declaration in both paths: rejected — the
  function may be far from the directive, and the directive-level anchor is
  the most actionable position for the CLI.
- Renaming the pseudo-rule so the nolint filter does not match: rejected —
  the filter matches on the analyzer name (`gohumanize`), not the rule ID, so
  no rule name change avoids it.
