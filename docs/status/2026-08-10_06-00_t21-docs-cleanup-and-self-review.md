# Status Report: T21 Docs Cleanup & Self-Review

**Date:** 2026-08-10 06:00 UTC
**Session scope:** Review TODO T21 (propose `ExitCodeFromReportConfidence` upstream), close documentation drift, self-review

---

## What This Session Did

### Task: Review T21

T21 asked to:

1. Open PR to `go-linter-sdk` with `ExitCodeFromReportConfidence(report, minConfidence)`
2. Replace CLI's local `exitCodeFromReport` with the upstream version once merged

**Finding:** Both items were **already done in code** but all documentation was stale.

- The SDK ships `ExitCodeByConfidence(report, threshold)` at `go-linter-sdk/registry.go:358` (tested at `registry_test.go:976-1047`, documented in SDK CHANGELOG, README, FEATURES.md, example_test.go).
- The CLI calls it at `main.go:232`: `os.Exit(linter.ExitCodeByConfidence(filteredReport, finding.ConfidenceHigh))`.
- No local `exitCodeFromReport()` function exists in any `.go` file anywhere in the project.
- The function was originally proposed as `ExitCodeFromReportConfidence` but shipped under the shorter name `ExitCodeByConfidence`.

### Documentation Fixed (4 files)

| File                                           | What was wrong                                                                                                                                                                                                                                                       | What I did                                                                                                                                                                                        |
| ---------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `TODO_LIST.md`                                 | T21 listed as `_planned_` with unchecked boxes; summary table row present                                                                                                                                                                                            | Deleted T21 entirely (summary row + detail section + orphaned "Upstream contributions" header). Per TODO_LIST.md policy: "When a task is finished, delete it here and record it in CHANGELOG.md." |
| `AGENTS.md`                                    | 3 stale references: (1) "implemented in CLI, not SDK" + "CLI uses its own exitCodeFromReport()" in confidence section, (2) "CLI uses its own exitCodeFromReport()" in CLI-ternary-exit-codes gotcha, (3) `exitCodeFromReport` at end of runScan pipeline description | Updated all 3 to reference `linter.ExitCodeByConfidence` and the SDK                                                                                                                              |
| `docs/adr/0003-confidence-aware-exit-codes.md` | Status "Accepted", body says "Why in the CLI, not the SDK?", consequences say CLI diverges from SDK                                                                                                                                                                  | Status changed to "Updated (2026-08-10)"; added Resolution section documenting SDK adoption and supersession of the CLI-local decision                                                            |
| `CHANGELOG.md`                                 | (1) Listed `exitCodeFromReport()` as a CLI feature, (2) Referenced `TestExitCodeFromReport` (test no longer exists), (3) ADR 0003 description said "decision to implement it in the CLI rather than in go-linter-sdk"                                                | Updated all 3 entries to reflect SDK-provided `ExitCodeByConfidence`; corrected test name to `TestExitCodeByConfidence`                                                                           |

### Verification

- `nix run .#test` results: 3 packages pass, 1 fails (`TestCustomGCLIntegration`)
- Confirmed `TestCustomGCLIntegration` failure is **pre-existing** (fails identically on clean HEAD `76b13ab`) — caused by published SDK version lacking `ToolName` field in `RuleMeta` that the local `replace` directive provides
- Remaining `exitCodeFromReport` / `ExitCodeFromReportConfidence` references are all in historical status reports and planning docs (`docs/status/`, `docs/planning/`) — these are point-in-time snapshots that should not be retroactively edited

---

## a) FULLY DONE

1. **T21 code** — SDK `ExitCodeByConfidence` shipped, CLI uses it (was done before this session)
2. **T21 documentation cleanup** — TODO_LIST.md, AGENTS.md, ADR 0003, CHANGELOG.md all updated (done this session)
3. **Test verification** — all packages pass except pre-existing `TestCustomGCLIntegration` failure
4. **Historical doc preservation** — correctly left `docs/status/` and `docs/planning/` snapshots untouched

---

## b) PARTIALLY DONE

Nothing this session.

---

## c) NOT STARTED

1. ~~**Fix `TestCustomGCLIntegration`** — pre-existing failure: `golangci-lint custom` build fails because published SDK version lacks `ToolName` field in `RuleMeta`. Needs a new `go-linter-sdk` tag published and `vendorHash` updated in flake.nix.~~ done — SDK v0.2.0 ToolName bump unblocked it (verified end-to-end)
2. ~~**T1 — Tag v0.2.0** — all code shipped to main but no git tag exists~~ done at `19bdd44`
3. ~~**T2 — Real-world validation sweep** — not started~~ done at `99b93cf`
4. ~~**T18 — Publish to golangci-lint plugin index** — blocked on T1~~ **Won't implement — declined 2026-09-19 (T18).**

---

## d) TOTALLY FUCKED UP

Nothing this session. The work was clean — documentation-only changes, verified against clean HEAD.

**Pre-existing issue worth flagging:** The `TestCustomGCLIntegration` test has been broken since the local SDK diverged from the published version (the `ToolName` field was added to `RuleMeta` locally via `replace` directive but never published). This means the golangci-lint module plugin integration path is **untested in CI**. This is a significant gap.

---

## e) WHAT WE SHOULD IMPROVE

1. **Doc drift detection** — T21's code was done but docs stayed stale across 4 files for multiple sessions. Consider a CI check or pre-commit hook that greps for known-removed symbols (`exitCodeFromReport`, `ExitCodeFromReportConfidence`) in non-historical docs.

2. **Historical docs vs living docs boundary** — The project has 15+ status reports in `docs/status/` that reference removed functions. This is correct (they're snapshots), but a newcomer might read them as current. Consider a header banner on historical docs: `> **Historical snapshot — not current.** See [AGENTS.md] / [TODO_LIST.md] for current state.`

3. **TODO_LIST.md discipline** — T21 should have been deleted when the code was done. The TODO_LIST header says "When a task is finished, delete it here" but this wasn't followed. Multiple prior sessions marked T21 as "I should not make external PRs without user direction" without realizing the PR was already merged.

4. **SDK version skew** — The local `replace` directive in go.mod points to a newer SDK than what's published. The `TestCustomGCLIntegration` test downloads the published version, causing a build failure. Either publish the SDK or make the integration test use the local replace.

5. **ADR update protocol** — ADR 0003 said "Accepted" with "Why in the CLI, not the SDK?" reasoning that was later invalidated. ADRs should get a "Resolution" or "Superseded" status update immediately when the decision changes, not weeks later during a docs audit.

6. **CHANGELOG test-name accuracy** — The CHANGELOG referenced `TestExitCodeFromReport` which doesn't exist (the test is `TestExitCodeByConfidence` in `main_test.go`). CHANGELOG entries referencing test names should be verified at write time.

7. **CHANGELOG entry for this session's doc fixes** — I did not add a CHANGELOG entry for the documentation corrections themselves. These are unreleased v0.2.0 changes, so the corrections are folded into the existing Unreleased entries. This is probably fine but worth noting.

---

## f) Up to 50 Things to Get Done Next

### Release (Critical Path)

1. ~~**Tag v0.2.0** (T1) — code is on main, just needs `git tag v0.2.0 && git push --tags`~~ done at `19bdd44`
2. ~~**Publish new go-linter-sdk version** — the local `replace` has `ToolName` in `RuleMeta` that published version lacks~~ done — SDK published (v0.2.0+; repo pins v0.3.1)
3. ~~**Update vendorHash in flake.nix** after new SDK tag~~ done — vendorHash maintained through dep bumps
4. ~~**Fix TestCustomGCLIntegration** — blocked on #2~~ done at `c084fd9`
5. ~~**Verify `golangci-lint custom` works end-to-end** after SDK publish~~ done — TestCustomGCLIntegration verified end-to-end in v0.2.0

### SDK Upstream (go-linter-sdk repo)

6. ~~**Verify go-linter-sdk has a tag** matching the version go-humanize-linter depends on~~ done — go-linter-sdk tags verified (v0.2.0 → v0.3.1 pinned)
7. ~~**Check if `ToolName` field in `RuleMeta` is published** or still local-only~~ done — ToolName shipped in SDK v0.2.0
8. ~~**Review go-linter-sdk CHANGELOG** for any other shipped features not yet consumed~~ done — SDK v0.2.0/v0.3.1 features consumed
9. ~~**Consider adding `ExitCode` typed constants** (0/1/2) to the SDK instead of bare `int` returns~~ **Won't implement — not adopted — SDK API frozen.**
10. ~~**SDK example for `ExitCodeByConfidence`** — verify the example_test.go ExampleExitCodeByConfidence is up to date~~ done — example_test.go examples render on pkg.go.dev (verified 2026-09-26)

### Testing

11. ~~**Add CI guard for doc drift** — grep non-historical docs for removed symbols~~ **Won't implement — not adopted — audits serve this role.**
12. ~~**Add integration test for CLI ternary exit codes** (exit 0/1/2 via subprocess)~~ done — covered — CLI smoke + TestCLI_MinConfidence
13. ~~**Verify `--min-confidence` flag actually affects exit code correctly** (filtering + exit code interaction)~~ done — TestCLI_MinConfidence covers the interaction
14. ~~**Add test for `ExitCodeByConfidence` with `ConfidenceFull`** threshold (currently tests use `ConfidenceHigh`)~~ done — covered — plugin min_confidence_full subtest
15. ~~**Plugin-path confidence exit code test** — verify plugin doesn't break with confidence filtering~~ **Won't implement — n/a — plugin reports diagnostics, not exit codes.**

### Documentation

16. ~~**Add historical-doc banner** to all `docs/status/` and `docs/planning/` files~~ **Won't implement — rejected — banners are an anti-pattern; annotation-in-place used instead.**
17. ~~**Audit all CHANGELOG entries** for accuracy against current code~~ done — multiple CHANGELOG audits (08-05 06-26, 2026-09-26)
18. ~~**Review ADR 0001** — is it still accurate?~~ done — ADR 0001 verified accurate in audits
19. ~~**Review ADR 0002** — is it still accurate?~~ done — ADR 0002 verified accurate
20. ~~**Review ADR 0004** (behavior-delta) — is it still accurate?~~ done — ADR 0004 verified accurate
21. ~~**Review ADR 0005** (H009/H002 overlap) — is it still accurate?~~ done — ADR 0005 verified accurate
22. ~~**Review ADR 0006** (per-statement suppression) — is it still accurate?~~ done — ADR 0006 verified accurate
23. ~~**README.md audit** — does it reference the SDK's `ExitCodeByConfidence` or the old local function?~~ done — README uses ExitCodeByConfidence
24. ~~**FEATURES.md audit** — is the exit-code feature accurately described?~~ done — FEATURES describes the confidence/exit system
25. ~~**Verify doc.go example** uses current API (`ExitCodeFromReport` is fine, but check context)~~ done — doc.go example updated (6a3267f)

### Code Quality

26. ~~**`filterReportByConfidence` uses `finding.ByConfidenceAtLeast`** — verify this filter and `ExitCodeByConfidence` threshold are semantically aligned (both use `Compare >= 0`)~~ done — verified — both use Compare >= threshold semantics
27. ~~**Exit code threshold is hardcoded to `ConfidenceHigh`** in main.go:232 — should this be configurable via a flag separate from `--min-confidence`?~~ **Won't implement — deliberate — fixed ConfidenceHigh threshold per ADR 0003.**
28. ~~**Review whether `--min-confidence full` + all-full findings produces exit 1** correctly (filtering removes nothing, all at threshold)~~ done — verified in audits (full+full → exit 1)
29. ~~**Review whether `--min-confidence full` + only-medium findings produces exit 0** correctly (all filtered out, empty report)~~ done — verified in audits (filtered-empty → exit 0)
30. ~~**Doc the exit-code/filtering interaction** more prominently in `--help` output~~ done — README documents the scheme

### Validation (T2)

31. ~~**Run linter against dustin/go-humanize itself** — should find zero findings~~ **Won't implement — fun idea, never adopted.**
32. **Run linter against 5+ real Go projects** from GitHub
33. **Collect false positive rate** per rule
34. **Document projects tested and results**
35. **File issues for any real bugs found** in detection

### Distribution

36. **Publish to golangci-lint plugin index** (T18, blocked on T1)
37. **Verify GitHub Action works** with published v0.2.0
38. **Update action.yml** version references if needed
39. **Create GitHub Release** with release notes from CHANGELOG

### Architecture / Future

40. **Consider severity-based exit codes** (not just confidence) — SDK mentions this as possible future
41. **Consider a `--fail-on` flag** for flexible CI exit-code configuration
42. **Explore SARIF report confidence field** — does SARIF support confidence? Are we populating it?
43. **Review baseline/behavior-delta feature** for confidence-awareness
44. **Consider machine-readable exit-code explanation** in JSON/SARIF output

### Maintenance

45. **Run `nix flake check`** — verify flake is healthy
46. **Run `nix run .#lint`** — verify no new lint issues
47. **Run `nix run .#vet`** — verify no vet issues
48. **Clean up any stale branches** in the repo
49. **Review `.custom-gcl.yml`** for accuracy against current module structure
50. **Review `go.mod` replace directives** — should the local SDK replace be removed now that the function is published?

---

## g) Questions (cannot figure out myself)

### Q1: Should the CLI exit-code threshold (`ConfidenceHigh`) be configurable via a separate flag?

Currently `main.go:232` hardcodes `finding.ConfidenceHigh` as the exit-code threshold regardless of `--min-confidence`. This means `--min-confidence low` (shows all findings) still uses exit 1 only for high+ findings, exit 2 for medium/low. Should there be a separate `--exit-threshold` flag, or is the current behavior (report shows everything, exit code uses fixed threshold) intentional and desirable?

### Q2: Is the `TestCustomGCLIntegration` failure known and accepted, or should I fix it?

The test fails because the published `go-linter-sdk` version lacks the `ToolName` field in `RuleMeta` that the local `replace` directive provides. Fixing it requires publishing a new SDK tag and updating `vendorHash`. Is this already tracked somewhere I haven't seen, or should I proceed to publish the SDK and fix the test?

### Q3: Should I publish the new go-linter-sdk version (with `ToolName`, `ExitCodeByConfidence`, `FilterRules`, etc.) now?

The local `replace` directive gives access to features the published version doesn't have. Multiple things depend on this (the `ToolName` field, `ExitCodeByConfidence`, etc.). Publishing would unblock `TestCustomGCLIntegration`, enable proper versioning, and let other consumers use these features. Should I proceed with tagging and publishing, or do you want to review the SDK changes first?
