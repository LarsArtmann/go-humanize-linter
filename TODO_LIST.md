# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                      | Tier   | Effort | Status  |
| --- | --------------------------------------------------------- | ------ | ------ | ------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                  | High   | XS     | blocked |
| T2  | Real-world validation sweep with new detection + features | High   | M      | planned |
| T15 | Plugin integration test through `custom-gcl` binary       | Medium | S      | planned |
| T18 | Publish to golangci-lint plugin index                     | Low    | S      | blocked |
| T19 | Per-statement `//nolint` suppression support              | Medium | M      | planned |
| T21 | Propose `ExitCodeFromReportConfidence` upstream           | Low    | S      | planned |

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
H001 size-bucket filter, import-alias-aware detection, dot-import support, and
the gogenfilter generated-file integration) have not been validated against the
corpus with these features active.

- [ ] Run the linter over the 327-project corpus with all 9 rules enabled
- [ ] Run `--verify-suppressions` on the corpus and record stale-directive rate
- [ ] Run `--min-confidence high` and compare finding counts
- [ ] Verify the gogenfilter integration does not skip hand-written files erroneously
- [ ] Save results to `docs/validation/`

---

## Plugin robustness

### T15 — Plugin integration test through `custom-gcl` binary · Medium · _planned_

Plugin registration is verified via `TestPluginRegisteredWithGolangciLint` and
`TestPluginRegisteredWithSettings`. Analysistest-based tests
(`TestRunDetector_FiltersByConfidence`, `TestRunDetector_VerifySuppressions`,
`TestRunDetector_H0SUPBypassesConfidenceFilter`) cover the `runDetector` pipeline.
However, no test exercises the full golangci-lint runtime: building the
`custom-gcl` binary, running it on testdata with `minConfidence` and
`verifySuppressions` enabled, and asserting the diagnostics appear correctly.

- [ ] Write a test that builds `custom-gcl`, runs it on testdata, asserts findings
- [ ] Gate behind `testing.Short()` skip (requires `golangci-lint custom` — network + git clone)

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
