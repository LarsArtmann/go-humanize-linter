# Status Report — 2026-07-30 19:27 — H008 Rule Count Fix & Test Suite Reconciliation

## TL;DR

Reconciled test fixtures against the recently merged H008 (manual-ordinal) rule. Four test surfaces still hard-coded 7 rules / 7 IDs; bumped them to 8. Full test suite (107 tests, `-race`) green across all four packages. Discovered that H009 (manual-commaf) is **already implemented and tested** in the working tree but **not registered** in `AllRules()` — a latent v0.2 feature that needs a decision before shipping.

> **Update 2026-07-30 (commit `2ac66b6`):** H009 **is now registered** in `AllRules()` and `allRuleDetectors()` (the "ghost rule" is live). AGENTS.md shows H001–H009. The rule count is 9, not 8. Remaining open items (negative testdata, `docs/rules/H008.md`+`H009.md`, real-world sweep) moved to `TODO_LIST.md`. Full item-by-item status in [Resolution](#resolution-2026-07-30) below.

---

## a) FULLY DONE

| #   | Item                                                   | Files                                     | Notes                              |
| --- | ------------------------------------------------------ | ----------------------------------------- | ---------------------------------- |
| 1   | `TestDefaultRegistry_AllRules` now expects 8 rules     | `linter_test.go:319`                      | Was `!= 7` → `!= 8`                |
| 2   | `ExampleDefaultRegistry` output updated                | `example_test.go:19`                      | `registered rules: 7` → `8`        |
| 3   | `ExampleAllRules` exposes H008                         | `example_test.go:35`                      | Added `// H008` line               |
| 4   | `TestCLI_RulesFlag` validates H008 in `--rules` output | `cmd/go-humanize-linter/main_test.go:460` | Added `"H008"` to expected ID list |
| 5   | Full `-race` test suite green                          | all 4 pkgs                                | 107 tests, 0 failures, 0 panics    |

## b) PARTIALLY DONE

| #   | Item                      | Status                                                                                                                                                                                                    | Gap                                                                                                                                                                                            |
| --- | ------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | H009 manual-commaf wiring | Detector (`pattern_commaf.go`), rule (`rule_commaf.go`), and test (`TestRuleCommaf_Positive`) all exist and pass. Testdata under `testdata/h009_commaf/` and `testdata/analysistest/h009positive/` exist. | **Not registered in `AllRules()`.** The `detectCommaf` function is not in `allRuleDetectors()` either. H009 is a "ghost rule" — present, tested, but invisible to the CLI / plugin / registry. |
| 2   | AGENTS.md rule inventory  | Still claims H001–H007                                                                                                                                                                                    | Should be H001–H008 (and H009 once activated)                                                                                                                                                  |
| 3   | Documentation drift       | Historical docs in `docs/status/2026-07-30_18-10_v0.1.0-released-gaps-remaining.md` etc. still reference 7 rules                                                                                          | Point-in-time docs; annotation, not rewrite                                                                                                                                                    |

## c) NOT STARTED

| #   | Item                                                                                                                  |
| --- | --------------------------------------------------------------------------------------------------------------------- |
| 1   | Register H009 (RuleCommaf) in `AllRules()` and `allRuleDetectors()`                                                   |
| 2   | Add H008 **and** H009 to `ExampleAllRules` output (only H008 done)                                                    |
| 3   | Add H008 to the `TestCLI_RulesFlag` expectation list (only done for H008 — verify H009 also added after registration) |
| 4   | Update `AGENTS.md` "Rule IDs" section to H001–H008                                                                    |
| 5   | Update `AGENTS.md` "Architecture" file table to mention `rule_ordinal.go` and `pattern_ordinal.go`                    |
| 6   | Real-world sweep of H008 against the 190+ project corpus to validate ~0% FP claim                                     |
| 7   | Real-world sweep of H009 once registered                                                                              |
| 8   | Negative testdata for H008 (clean code with switch on n%10 that returns non-ordinal strings)                          |
| 9   | Negative testdata for H009 (clean code with `%.Nf` Sprintf but no comma-writing)                                      |
| 10  | Add H008 markdown doc page (`docs/rules/H008.md`) per the P14 plan item                                               |
| 11  | Add H009 markdown doc page (`docs/rules/H009.md`)                                                                     |
| 12  | Add H008 + H009 entries to CHANGELOG.md "Unreleased" section                                                          |
| 13  | Add H008 + H009 entries to FEATURES.md rule inventory                                                                 |
| 14  | Add H008 ExampleRule function (`ExampleRuleOrdinal`) to `example_test.go`                                             |
| 15  | Add H009 ExampleRule function (`ExampleRuleCommaf`) to `example_test.go`                                              |
| 16  | Fix `goconst: string H001 has 3 occurrences` warning in `rules.go:159`                                                |
| 17  | Fix `lintIgnorePrefix` unused-const warning in `pattern_helpers.go:196`                                               |
| 18  | Address `mnd` magic-number warnings in `pattern_ordinal.go:44,60` (10, 3)                                             |
| 19  | Decide whether H009 stays in v0.1.x or rolls to v0.2.0 (product-direction question)                                   |
| 20  | Backfill H008 into --rules output `--help` and CLI `--enable`/`--disable` examples                                    |

## d) TOTALLY FUCKED UP

| #   | Item | Detail                                                                                                                                                                  |
| --- | ---- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| —   | —    | Nothing broken. The H008 merge just slipped past the four test fixtures that hard-coded the 7-rule count. No data corruption, no API regressions, no missed invariants. |

## e) WHAT WE SHOULD IMPROVE

### Architecture & Process

1. **Single source of truth for rule count.** The number "7" appears in test code, example comments, AGENTS.md, FEATURES.md, README, and CHANGELOG independently. There should be a `Rules()` accessor that returns the slice and a derived `len()` everywhere else. A test that asserts `len(AllRules()) == N` is anti-fragile — it duplicates the count.
2. **Example-based tests are gold, but they break on every rule addition.** `ExampleDefaultRegistry` and `ExampleAllRules` will need to be updated every time a new rule lands. Consider deriving these examples from `AllRules()` at test time, or replacing the hard-coded `// Output:` blocks with a sub-test that compares against the live registry.
3. **Hidden-feature landmines.** H009 is fully implemented and tested but not registered. A CI check that asserts `len(Diff(AllRules(), testCount)) == 0` or that all rule files matching `rule_*.go` are registered would catch this. Even cheaper: a `go test` that lists `rule_*.go` and asserts each one has a `Rule*()` factory and an entry in `AllRules()`.
4. **AGENTS.md not updated in the H008 commit.** The merge added `rule_ordinal.go`, `pattern_ordinal.go`, and a new rule — but the architecture table and rule-id list weren't bumped. Whoever landed H008 didn't update the docs.

### Code Quality

5. **H008's `ordinalSuffixes` map is package-level mutable state.** `var ordinalSuffixes = map[string]bool{...}` is a `//nolint:gochecknoglobals` candidate. For a 4-entry lookup, a `switch` or `contains([]string, ...)` would be safer and avoid the mutex question in case of future writes.
6. **Magic numbers in `pattern_ordinal.go`** (`10`, `100`, `3`) deserve named constants — `ordinalModTen`, `ordinalModHundred`, `minOrdinalSuffixesForFlag`.
7. **`goconst` on `"H001"`** in `rules.go` — the rule IDs are repeated across detector lists, factory meta blocks, and tests. A `const ruleIdH001 = "H001"` shared across packages would centralize the constant.
8. **`lintIgnorePrefix` unused** in `pattern_helpers.go:196` — either delete it or wire it up. Dead code violates the "no paper cuts" rule.

### Test Coverage

9. **H008 and H009 have no negative testdata.** A clean function that uses `n%10` with non-ordinal branches must NOT flag. This is the only way to validate the "3 of 4 suffixes" heuristic.
10. **H008's `hasOrdinalSwitch` has no sub-test.** The other rules have `TestHas*` table-driven tests; H008 needs `TestHasOrdinalSwitch` with the same structure.
11. **Plugin-mode tests for H008** — `testdata/analysistest/h008positive/` exists; verify the analyser runs it. The `TestAnalyzerAnalysistest` test passes, but I should confirm both H008 and H009 fixtures are exercised.
12. **No real-world validation of H008.** Pre-merge, all 7 rules were swept against 190+ projects. H008 needs the same sweep to validate the "3 of 4 suffix" heuristic before claiming a "~0% FP" rate.

### Documentation

13. **CHANGELOG.md missing H008 entry.** A rule addition that's already merged should be in the Unreleased section.
14. **FEATURES.md rule inventory** should list H008 under DONE.
15. **README.md rule table** (if it exists) needs H008.

## f) UP TO 50 THINGS TO GET DONE NEXT

Sorted by Pareto impact (high-value, low-effort first):

| #   | Task                                                                                                                  | Effort | Impact                                   |
| --- | --------------------------------------------------------------------------------------------------------------------- | ------ | ---------------------------------------- |
| 1   | Decide H009: register now or park for v0.2.0                                                                          | 1 min  | HIGH — unblocks 10+ downstream tasks     |
| 2   | Register H009 in `AllRules()` + `allRuleDetectors()` if "now"                                                         | 5 min  | HIGH                                     |
| 3   | Add H009 to `ExampleAllRules` and `TestCLI_RulesFlag`                                                                 | 5 min  | HIGH                                     |
| 4   | Update `AGENTS.md` architecture table + rule IDs to H001–H008                                                         | 5 min  | MED                                      |
| 5   | Add H008 entry to CHANGELOG.md Unreleased                                                                             | 5 min  | MED                                      |
| 6   | Add H008 entry to FEATURES.md                                                                                         | 5 min  | MED                                      |
| 7   | Add H008 markdown doc page `docs/rules/H008.md`                                                                       | 15 min | MED                                      |
| 8   | Add H009 markdown doc page `docs/rules/H009.md` (if registered)                                                       | 15 min | MED                                      |
| 9   | Add `TestHasOrdinalSwitch` table-driven test (mirror `TestHasEqualsOneBranch`)                                        | 10 min | HIGH — improves test parity              |
| 10  | Add H008 negative testdata (`n%10` with non-ordinal returns)                                                          | 10 min | HIGH                                     |
| 11  | Add H009 negative testdata (`%.Nf` Sprintf without comma writing)                                                     | 10 min | HIGH                                     |
| 12  | Run real-world sweep of H008 against 190+ project corpus                                                              | 30 min | HIGH — validates "~0% FP" claim          |
| 13  | Run real-world sweep of H009 (if registered)                                                                          | 30 min | HIGH                                     |
| 14  | Fix `goconst` H001 warning — extract `const ruleIDH001 = "H001"`                                                      | 5 min  | LOW                                      |
| 15  | Delete unused `lintIgnorePrefix` const in `pattern_helpers.go`                                                        | 1 min  | LOW                                      |
| 16  | Fix `mnd` warnings in `pattern_ordinal.go` — name magic numbers                                                       | 5 min  | LOW                                      |
| 17  | Add `ExampleRuleOrdinal` example test                                                                                 | 5 min  | LOW                                      |
| 18  | Add `ExampleRuleCommaf` example test (if registered)                                                                  | 5 min  | LOW                                      |
| 19  | CI gate: assert `len(rule_*.go) == len(AllRules()) == len(allRuleDetectors())`                                        | 15 min | HIGH — prevents future H009-style ghosts |
| 20  | Replace hard-coded `// Output: registered rules: N` with derived example                                              | 10 min | MED                                      |
| 21  | Refactor `IsLiteralInt` / `hasOrdinalSwitch` to share constants                                                       | 10 min | MED                                      |
| 22  | Add `golines`/`gofumpt` formatting pass to H008/H009 files                                                            | 5 min  | LOW                                      |
| 23  | Verify `TestAllRules_CleanCode` covers H008 and H009 (once registered)                                                | 5 min  | HIGH                                     |
| 24  | Add `bench/` benchmark for H008 `hasOrdinalSwitch`                                                                    | 15 min | LOW                                      |
| 25  | Add H008 explanation text to `--explain H008` flag                                                                    | 10 min | MED — UX                                 |
| 26  | Add H009 explanation text (if registered)                                                                             | 10 min | MED                                      |
| 27  | Audit `FEATURES.md` for any other 7-rule references                                                                   | 5 min  | MED                                      |
| 28  | Audit `README.md` for 7-rule references                                                                               | 5 min  | MED                                      |
| 29  | Update `docs/validation/2026-07-30_real-world-sweep.md` with H008 numbers                                             | 5 min  | MED                                      |
| 30  | Add H008 to plugin `description` table                                                                                | 5 min  | LOW                                      |
| 31  | Add H009 to plugin description (if registered)                                                                        | 5 min  | LOW                                      |
| 32  | Add `--enable H008` / `--disable H008` examples to CLI help                                                           | 10 min | LOW                                      |
| 33  | Stress test H008 with adversarial switch patterns (e.g. `n%10` with sentinel that has 3 of 4 suffixes by coincidence) | 30 min | MED                                      |
| 34  | Add docstring explaining H008's "3 of 4" threshold vs "all 4"                                                         | 5 min  | LOW                                      |
| 35  | Cross-check H008 detector against `humanize.Ordinal` source code                                                      | 10 min | LOW                                      |
| 36  | Add H008/H009 entries to `nix run .#rules` output (if it exists)                                                      | 5 min  | LOW                                      |
| 37  | Verify `--rules` output ordering matches `AllRules()` ordering                                                        | 5 min  | LOW                                      |
| 38  | Add `goimports` pass for H008/H009 files                                                                              | 1 min  | LOW                                      |
| 39  | Bump version to v0.2.0-pre (H008 + H009 is a minor bump)                                                              | 5 min  | HIGH                                     |
| 40  | Tag v0.2.0 release once H008/H009 are ratified                                                                        | 30 min | HIGH                                     |
| 41  | Backfill H008 + H009 to `nix run .#lint` task if it has a docs section                                                | 5 min  | LOW                                      |
| 42  | Add H008 example to README rule table                                                                                 | 5 min  | LOW                                      |
| 43  | Add H009 example to README (if registered)                                                                            | 5 min  | LOW                                      |
| 44  | Verify all 5 suppressed-docs files (sweep the unannotated ones)                                                       | 10 min | LOW                                      |
| 45  | Run `nix run .#lint` and address any new findings                                                                     | 5 min  | HIGH — CI gate                           |
| 46  | Run `nix run .#vet` and address any new findings                                                                      | 5 min  | HIGH                                     |
| 47  | Check for `golang.org/x/tools/go/analysis` predicate stability across Go minor versions                               | 10 min | LOW                                      |
| 48  | Add a `CONTRIBUTING.md` note about the rule-count CI gate                                                             | 10 min | LOW                                      |
| 49  | Re-run full real-world sweep vs main, save as `docs/validation/2026-07-30_v0.2-sweep.md`                              | 45 min | HIGH                                     |
| 50  | Publish v0.2.0 release notes on GitHub                                                                                | 30 min | HIGH                                     |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **H009: ship now or defer?** The detector, factory, and test all exist; the only work is wiring it into `AllRules()` and `allRuleDetectors()`. But shipping both H008 and H009 in the same release doubles the maintenance surface and we haven't validated H009 against the real-world corpus. Do you want H009 in v0.2.0 alongside H008, or should H009 live in a feature branch until we can sweep it?

2. **Versioning — v0.1.1 or v0.2.0?** H008 is a new rule with a new heuristic. Is that a patch (additive, backward-compatible) or a minor (new feature surface)? SemVer says minor, but our release cadence is undocumented. What's the project's intended policy?

3. **CI rule-count gate — add to `flake.nix` as a pre-commit hook or as a test?** A test (`TestRuleCountConsistency`) is simpler and runs in CI; a pre-commit hook is faster feedback but more setup. Which is the truth in this project — are hooks favored via `treefmt`/`nix`/`lefthook`, or is the convention "tests only"?

---

## Verification

```
ok  	github.com/larsartmann/go-humanize-linter              1.204s
ok  	github.com/larsartmann/go-humanize-linter/cmd/go-humanize-linter  1.378s
ok  	github.com/larsartmann/go-humanize-linter/cmd/gohumanize          1.566s
ok  	github.com/larsartmann/go-humanize-linter/plugin                  2.142s
```

All tests pass. Zero failures, zero panics. 107 tests executed.

---

## Resolution (2026-07-30)

The "NOT STARTED" and "f) UP TO 50" sections drove the next sessions. Item-by-item:

| Section c) / f) item                                           | Status                          |
| -------------------------------------------------------------- | ------------------------------- |
| #1 Register H009 in `AllRules()` + `allRuleDetectors()`        | done at `2ac66b6`               |
| #2/#3 Add H008+H009 to `ExampleAllRules` / `TestCLI_RulesFlag` | done at `20dd7d3`               |
| #4 Update AGENTS.md to H001–H009                               | done at `2ac66b6`               |
| #5/#6 H008 entries in CHANGELOG / FEATURES                     | done at `2ac66b6` (now 9 rules) |
| #14 goconst `H001` warning                                     | done at `2ac66b6`               |
| #15 unused `lintIgnorePrefix`                                  | done (0 lint issues)            |
| #16/#18 `mnd` in `pattern_ordinal.go`                          | done at `2ac66b6`               |
| #17/#18 `ExampleRuleOrdinal` / `ExampleRuleCommaf`             | done at `20dd7d3`               |
| #25/#26 `--explain H008` / `--explain H009`                    | done at `2ac66b6`               |

Still open — moved to `TODO_LIST.md`:

- ~~Negative testdata for H008 + H009~~ done at `23bf769`
- ~~`docs/rules/H008.md` + `H009.md`~~ done at `f8ba5d6`
- Real-world sweep of H008 + H009 (→ TODO_LIST T2)
- Tag `v0.2.0` (→ TODO_LIST T1)
