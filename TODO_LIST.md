# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                              | Tier   | Effort | Status  |
| --- | ----------------------------------------------------------------- | ------ | ------ | ------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                          | High   | XS     | blocked |
| T2  | Real-world validation sweep with new detection + features         | High   | M      | planned |
| T15 | Plugin integration test through `custom-gcl` binary               | Medium | S      | planned |
| T16 | Dot-import (`. "strings"`) support for alias resolution           | Low    | S      | planned |
| T17 | H009/H002 overlap disambiguation                                  | Low    | M      | planned |
| T18 | Publish to golangci-lint plugin index                             | Low    | S      | blocked |
| T19 | Per-statement `//nolint` suppression support                      | Medium | M      | planned |
| T20 | `--behavior-delta` flag for regression testing                    | Low    | M      | planned |
| T21 | Propose `ExitCodeFromReportConfidence` upstream                   | Low    | S      | planned |
| T22 | Protect `findingToTokenPos` against out-of-range line numbers     | High   | XS     | planned |
| T23 | Exclude H0SUP findings from confidence filtering                  | Medium | XS     | planned |
| T24 | Remove dead `RunOverPackage` method                               | Low    | XS     | planned |

---

## Release & validation

### T1 — Tag `v0.2.0` (code shipped; tag missing) · High · _blocked_

All v0.2.0 features are merged to `main` and documented in `CHANGELOG.md` under
`[0.2.0] - Unreleased`. The only git tag is `v0.1.0`. The release workflow
(`/.github/workflows/release.yml`) fires on tags.

- [ ] Tag `v0.2.0` on `main` (requires explicit user approval — never tag without it)
- [ ] Verify the release workflow fires and publishes GitHub release notes

### T2 — Real-world validation sweep with new detection + features · High · _planned_

H001–H009 were swept against 327 Go projects
(`docs/validation/2026-07-31_real-world-sweep.md`), but the new features
(`--verify-suppressions`, `--min-confidence`, plugin confidence filtering, the
H001 size-bucket filter, import-alias-aware detection, and the gogenfilter
generated-file integration) have not been validated against the corpus with
these features active.

- [ ] Run the linter over the 327-project corpus with all 9 rules enabled
- [ ] Run `--verify-suppressions` on the corpus and record stale-directive rate
- [ ] Run `--min-confidence high` and compare finding counts
- [ ] Verify the gogenfilter integration does not skip hand-written files erroneously
- [ ] Save results to `docs/validation/`

---

## Plugin robustness

### T22 — Protect `findingToTokenPos` against out-of-range line numbers · High · _planned_

`findingToTokenPos` (`plugin/plugin.go`) guards against `f.Position.Line < 1`
but does NOT guard against line numbers exceeding the file's actual line count.
`token.File.LineStart(line)` panics with "illegal line number (past end of
file)" on out-of-range input. While unlikely in normal operation (findings use
the same `FileSet`), a defensive check is required — a linter crashing is worse
than a linter reporting at position 0.

- [ ] Add `tf.LineCount() >= f.Position.Line` check before calling `LineStart`
- [ ] Add `TestFindingToTokenPos_OutOfRangeLine` panic-recovery test

### T23 — Exclude H0SUP findings from confidence filtering · Medium · _planned_

Suppression-verification findings (H0SUP) carry `ConfidenceHigh` (0.75). If a
user sets `minConfidence: "full"` (1.0), all H0SUP findings are filtered out.
This is wrong — suppression verification is a meta-diagnostic about directive
correctness, not a confidence-rated code-pattern finding. It should bypass
confidence filtering in `runDetector`.

- [ ] In `runDetector`, exclude H0SUP findings from the confidence filter check
- [ ] Add `TestRunDetector_H0SUPBypassesConfidenceFilter`

### T15 — Plugin integration test through `custom-gcl` binary · Medium · _planned_

Plugin registration is verified via `TestPluginRegisteredWithGolangciLint` and
`TestPluginRegisteredWithSettings`. However, no test exercises the full
`runDetector` pipeline through the golangci-lint runtime: building the
`custom-gcl` binary, running it on testdata with `minConfidence` and
`verifySuppressions` enabled, and asserting the diagnostics appear correctly.

- [ ] Write a test that builds `custom-gcl`, runs it on testdata, asserts findings
- [ ] Gate behind `testing.Short()` skip (requires `golangci-lint custom` — network + git clone)

---

## Detection improvements

### T16 — Dot-import (`. "strings"`) support for alias resolution · Low · _planned_

The `buildImportAliases` function resolves named import aliases (e.g.,
`str "strings"`) but does not handle dot imports (`. "strings"`). Dot imports
inject all exported names into the current scope, so `HasSuffix` would be called
without a package qualifier. ADR 0001 documents this as a known gap.

- [ ] Detect dot imports in `buildImportAliases` and handle bare calls to
      `HasSuffix`, `TrimRight`, etc. in pattern helpers

### T17 — H009/H002 overlap disambiguation · Low · _planned_

H009 (manual-commaf) and H002 (manual-comma-format) can both match the same
code when it involves `%.Nf` formatting plus comma-grouping loops. Currently
both rules fire independently. A prioritization or suppression mechanism
would reduce noise.

- [ ] When H009 fires, suppress H002 on the same function (or vice versa)
- [ ] Document the precedence rule

---

## Suppression & confidence features

### T19 — Per-statement `//nolint` suppression support · Medium · _planned_

Currently, `//nolint:gohumanize` directives are matched at the function-declaration
level (the finding position is `fn.Pos()`). AIs and developers often want to
suppress a finding on a specific line or statement inside a function. This
requires either per-statement findings (each detector returns a specific
`token.Pos` instead of `fn.Pos()`) or a line-range-based suppression matcher.

- [ ] Decide approach: per-statement `token.Pos` in detectors vs. line-range matching
- [ ] Implement the chosen approach
- [ ] Add testdata for per-statement suppression

### T20 — `--behavior-delta` flag for regression testing · Low · _planned_

A `--behavior-delta <baseline.json>` flag would compare the current run's
findings against a saved baseline and report additions/removals. This is useful
for detecting false-positive regressions when detector logic changes.

- [ ] Implement baseline loading and comparison
- [ ] Report added findings (potential new false positives) and removed findings (potential missed detections)
- [ ] Exit code: 0 = no delta, 1 = delta found

---

## Code cleanup

### T24 — Remove dead `RunOverPackage` method · Low · _planned_

`HumanizeDetector.RunOverPackage()` in `rules.go` was the old plugin entry
point. The plugin now uses `runDetector` which inlines the same logic.
`RunOverPackage` has no live callers (verified via grep — only a commented-out
reference exists).

- [ ] Remove `RunOverPackage` or deprecate with a `// Deprecated:` comment
- [ ] Update `CHANGELOG.md` Removed section

---

## Upstream contributions

### T21 — Propose `ExitCodeFromReportConfidence` upstream · Low · _planned_

The CLI implements its own `exitCodeFromReport()` with ternary exit codes
(0=clean, 1=must fix, 2=triage). The SDK's `linter.ExitCodeFromReport` is still
binary (0 or 1). Proposing `ExitCodeFromReportConfidence` upstream would let
all SDK-based linters benefit from confidence-aware exit codes without
re-implementing the logic.

- [ ] Open PR to `go-linter-sdk` with `ExitCodeFromReportConfidence(report, minConfidence)`
- [ ] Replace CLI's local `exitCodeFromReport` with the upstream version once merged

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (v0.2.0 not tagged yet — see T1)
- [ ] Submit to the plugin index once installable
