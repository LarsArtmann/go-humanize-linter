# Real-World Validation Sweep

**Date:** 2026-07-30 (Session 2)
**Scope:** All Go projects under `~/projects/` (190+ repos scanned)
**Linter build:** pre-v0.1.0, rules H001–H007
**Full report:** [`../status/2026-07-30_17-48_real-world-validation-and-plugin-mode.html`](../status/2026-07-30_17-48_real-world-validation-and-plugin-mode.html)

> This is a **point-in-time snapshot** of detection accuracy on real-world code.
> It is preserved as a baseline for regression tracking; numbers will drift as
> rules and targets evolve.

## Headline results

| Metric              | Value                                |
| ------------------- | ------------------------------------ |
| Repos scanned       | 190+                                 |
| Repos with findings | ~30                                  |
| Total findings      | 97                                   |
| Overall FP rate     | Near-zero (after the H004 fix below) |

## Findings by rule

| Rule | Count | Humanize API to use instead               |
| ---- | ----- | ----------------------------------------- |
| H001 | 27    | `humanize.Bytes` / `humanize.IBytes`      |
| H003 | 25    | `humanize.RelTime` / `humanize.Time`      |
| H004 | 24    | `humanize.Plural` / `humanize.PluralWord` |
| H005 | 10    | `humanize.ComputeSI`                      |
| H002 | 7     | `humanize.Comma`                          |
| H006 | 2     | `humanize.FormatFloat`                    |
| H007 | 2     | `humanize.ParseBytes`                     |

## Top-hit projects

| Project                 | Findings |
| ----------------------- | -------- |
| ast-state-analyzer      | 9        |
| SEC                     | 8        |
| KeyCountdown            | 8        |
| BuildFlow               | 8        |
| CreditReformBilanzampel | 7        |

## Accuracy notes

- **H001, H002, H003, H005, H006, H007** — near-zero false-positive rate across the
  sweep. Production-usable as-is.
- **H004** was the outlier at **~60% false-positive rate** in the raw scan. A
  two-stage discriminator was added during this session and brought it to ~0%:
  1. `branchContainsString(body)` — the `if x == 1` branch must contain a string
     literal or concatenation.
  2. `funcReturnsString(fn)` — the enclosing function must return `string` or
     `[]string`.

  Effect: **44 → 24 findings**, all 24 remaining judged genuine pluralization.

- Two borderline H004 findings on wave-label functions
  (`project-dependency-graph`) were judged acceptable noise: count-based text
  selection where the "plural" is a full sentence rather than an `s` suffix.

## Known precision gap (carried forward)

The golangci-lint plugin reports every finding at the **FuncDecl position**
(`fn.Pos()`), not the precise line/column stored in `finding.Position`. The CLI
path is correct; only the plugin path loses precision. Tracked as **per-line
diagnostics** in [`../../ROADMAP.md`](../../ROADMAP.md) ("Precision & ergonomics").

## Reproduction

```bash
go-humanize-linter ~/projects/<repo>
```

Re-running against the same repos today may yield different counts as both the
linter and the target projects have moved on; this document captures the
session-2 baseline.
