# ADR 0004: Behavior Delta Baseline Comparison

**Date:** 2026-08-05
**Status:** Accepted

## Context

Detector refactors (broadening or narrowing pattern matches, changing
corroborating-signal thresholds, adjusting confidence assignment) can silently
shift the set of findings the linter produces. A reviewer reading a PR diff sees
code changes but has no easy way to see their blast radius on real-world code:

- Did the refactor introduce new false positives?
- Did tightening a threshold drop a real finding that users depend on?
- Did an alias-resolution fix suddenly light up previously-missed detections?

Without a way to compare before/after finding sets, these regressions are
discovered only when a user files a bug — sometimes months later.

`go-linter-sdk` reports findings as `finding.Finding` structs containing a
rule ID, message, suggestion, confidence, and position. Any of these fields
could theoretically be used as a comparison key.

## Decision

Add `--save-baseline <file>` and `--behavior-delta <file>` CLI flags backed by
`cmd/go-humanize-linter/behavior_delta.go`. The comparison key is the
**(rule, file, line)** tuple only — message text, suggestion text, and
confidence are deliberately ignored.

### Why (rule, file, line) and not the full finding?

1. **Message text changes are not behavioral.** Detector wording is refined
   constantly ("manual byte-size formatting" → "manual byte formatting with
   KMGTPE trick"). These are presentation improvements, not detection shifts.
   Including message text in the key would produce noise on every wording PR.

2. **Confidence changes are not behavioral.** Promoting a finding from
   `ConfidenceMedium` to `ConfidenceHigh` is a triage decision, not a
   detection change. The finding is still the same finding.

3. **Column changes are not behavioral.** A refactor that improves position
   precision (function-level → statement-level) is a UX improvement. The
   finding still points at the same line of code.

4. **Rule ID + file + line is the minimum stable identity.** If a different
   rule fires on the same line, or the same rule fires on a different line,
   that IS a behavioral change worth surfacing.

### Why a separate baseline file instead of inline diff?

- **Cross-machine comparison.** Baselines can be saved on `main` CI and
  downloaded by PR builds. An inline diff would require both runs in the same
  process.
- **Human-readable.** The JSON is a sorted array of `{rule, file, line}`
  entries — reviewable in a text editor, diffable in git.
- **Deterministic.** `sortEntries` sorts by (file, line, rule) so two runs
  over the same code produce byte-identical files. This makes baseline files
  merge-friendly in source control.

### Exit code semantics

| Exit | Meaning                                                            |
| ---- | ------------------------------------------------------------------ |
| 0    | No delta — current findings structurally match the baseline.       |
| 1    | Findings added or removed since the baseline (review required).    |
| 2    | Operational error — baseline file missing, unreadable, or invalid. |

Exit 2 uses the same code as `--min-confidence`'s triage exit, but the
meaning differs (error vs. triage). The CLI distinguishes them by context:
`--behavior-delta` runs exit 2 only on baseline-read failure, never on
finding confidence.

### Why `0o600` file mode?

Baseline files may be checked into source control alongside the code they
describe. Restrictive permissions (`0o600`) prevent other users on shared CI
runners from tampering with the baseline between the save and compare steps.

## Consequences

- **Positive:** Refactor PRs can be reviewed by their finding delta, not just
  their code diff. "This change drops 3 findings and adds 1" is actionable;
  "this change rewrites `isPackageCall`" is not.
- **Positive:** The baseline format is stable and JSON — external tooling
  (dashboards, bots) can consume it without parsing linter output.
- **Positive:** Deterministic sorting makes baselines merge-friendly in git.
- **Negative:** A detector that changes a finding's line number (e.g.,
  reporting at the call site instead of the function declaration) will
  register as a delta even though the underlying detection is unchanged.
  This is intentional — line-number shifts ARE user-visible and worth
  reviewing.
- **Negative:** Message-only changes require regenerating the baseline to
  silence the "no delta" check. This is the correct tradeoff: the baseline
  should reflect the current state, not the historical state.
- **Negative:** No schema versioning yet. If the baseline format changes
  (e.g., adding a `column` field), old baselines must be regenerated. A
  future `schema_version` field could address this if it becomes a problem.
