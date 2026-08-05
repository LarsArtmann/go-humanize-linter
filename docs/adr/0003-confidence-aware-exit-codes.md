# ADR 0003: Confidence-Aware Exit Codes

**Date:** 2026-08-05
**Status:** Accepted

## Context

The `go-linter-sdk` provides `ExitCodeFromReport(report)` which returns a
binary exit code: 0 (no findings) or 1 (any findings). This is the standard
Unix convention for linters.

However, the linter now assigns confidence levels to findings (`Full`,
`High`, `Medium`, `Low`). A finding with `ConfidenceMedium` might be a
false positive, while a finding with `ConfidenceFull` is almost certainly
real. CI pipelines that fail on any finding (exit 1) cannot distinguish
"must fix" from "review when you have time".

## Decision

Implement a **ternary exit code scheme** in the CLI, not in `go-linter-sdk`:

- **Exit 0** — No findings after confidence filtering.
- **Exit 1** — At least one finding with `ConfidenceHigh` or `ConfidenceFull`
  remains. These are "must fix" — almost certainly real reimplementations.
- **Exit 2** — Only `ConfidenceMedium` or `ConfidenceLow` findings remain.
  These are "triage" — likely false positives or borderline cases.

The `exitCodeFromReport()` function in `cmd/go-humanize-linter/main.go`
implements this logic. It replaces the SDK's `ExitCodeFromReport` in the
CLI path.

### Why ternary and not continuous?

Three states map cleanly to CI semantics:

- 0 = green (ship it)
- 1 = red (block the PR)
- 2 = yellow (warn but don't block)

More states would add complexity without actionable benefit. CI pipelines
can branch on `exit code == 2` to post a comment instead of failing the build.

### Why in the CLI, not the SDK?

1. **SDK stability.** `go-linter-sdk` is at v0.1.0. Changing
   `ExitCodeFromReport` would break all SDK consumers. Adding a new function
   (`ExitCodeFromReportConfidence`) is possible but requires an SDK release.

2. **Different linters have different needs.** Not every SDK-based linter
   assigns confidence levels. A ternary exit code is only useful when the
   linter produces meaningful confidence values.

3. **The CLI is the only consumer today.** The plugin path uses
   golangci-lint's exit code, which is always binary. Only the standalone
   CLI needs the ternary scheme.

### Filtering vs. exit codes

`--min-confidence` controls **what appears in the report**. The exit code
is computed from the **filtered** report. So `--min-confidence high` with
one high-confidence finding produces exit 1. `--min-confidence high` with
only medium findings produces exit 0 (they were filtered out).

## Consequences

- **Positive:** CI pipelines can distinguish "must fix" (exit 1) from
  "triage" (exit 2), enabling graduated enforcement.
- **Positive:** The SDK remains unchanged; no breaking change for other
  consumers.
- **Negative:** The CLI's `exitCodeFromReport` diverges from the SDK's
  `ExitCodeFromReport`. This creates a maintenance burden if the SDK later
  adds its own confidence-aware exit code (see TODO T21).
- **Negative:** Exit code 2 is unusual. Most CI systems expect 0 or 1.
  Pipelines that use `set -e` will treat exit 2 as failure. This is
  documented in `--help` output.
