# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                                      | Tier   | Effort | Status   |
| --- | ------------------------------------------------------------------------- | ------ | ------ | -------- |
| T1  | Wire H008 + H009 analysistest fixtures into the test run                  | High   | XS     | planned  |
| T2  | Negative testdata for H008 and H009                                       | High   | XS     | planned  |
| T3  | `TestRuleCountConsistency` anti-ghost-rule guard                          | High   | XS     | planned  |
| T4  | `TestHasOrdinalSwitch` table-driven test                                  | High   | XS     | planned  |
| T5  | Rule doc pages: `docs/rules/H008.md` + `H009.md`                          | Medium | XS     | planned  |
| T6  | Tag `v0.2.0` (code shipped; tag missing)                                  | High   | XS     | planned  |
| T7  | Re-run real-world validation sweep with H008 + H009                       | High   | M      | planned  |
| T8  | CONTRIBUTING.md rule-addition checklist                                   | Medium | XS     | planned  |
| T9  | White-box unit tests for `pattern_commaf.go` helpers                      | Medium | S      | planned  |
| T10 | Replace `nix run .#lint` `grep -v` filter with `.golangci.yml` plugin reg | High   | M      | planned  |
| T11 | Configurable rules in plugin mode via `Analyzer.Flags`                    | Medium | M      | planned  |
| T12 | GitHub Action composite `action.yml`                                      | Medium | S      | planned  |
| T13 | Package-level `var` detection for H007                                    | Low    | M      | planned  |
| T14 | go/types type-aware detection                                             | Low    | L      | planned  |
| T15 | `--config` flag for YAML/TOML rule configuration                          | Low    | M      | planned  |
| T16 | Publish to golangci-lint plugin index                                     | Low    | S      | blocked  |

---

## Test gaps

### T1 — Wire H008 + H009 analysistest fixtures into the test run · High · _planned_

The fixtures `testdata/analysistest/h008positive/` and `h009positive/` exist, but
`TestAnalyzerAnalysistest` (`plugin/plugin_test.go:177`) only runs H001–H007 + `clean`.
The two new rules are never exercised through the full `analysis.Analyzer.Run` path.

- [ ] Add `"./h008positive"` and `"./h009positive"` to the `analysistest.Run` call
- [ ] Verify each fixture's `// want` diagnostic matches

### T2 — Negative testdata for H008 and H009 · High · _planned_

H001–H005 each have a `*_negative/` fixture (clean code that must NOT flag).
H006, H007, H008, and H009 have none. A clean ordinal `switch n { case 1: "1st" }`
(no `%` operator) and a `%.Nf` Sprintf without a comma loop must return no findings.

- [ ] `testdata/h008_negative/main.go` — `switch n%10` with non-ordinal returns
- [ ] `testdata/h009_negative/main.go` — `%.Nf` alone (no separator loop)

### T3 — `TestRuleCountConsistency` anti-ghost-rule guard · High · _planned_

H009 shipped as a "ghost rule" (detector + test existed but it was never registered
in `AllRules()`). A single test asserting
`len(rule_*.go factories) == len(AllRules()) == len(allRuleDetectors())` blocks a
repeat of that class of bug.

- [ ] Add a test that lists `rule_*.go` factories and asserts each has an entry in `AllRules()` and `allRuleDetectors()`
- [ ] Assert the three sources agree on count (currently 9)

### T4 — `TestHasOrdinalSwitch` table-driven test · High · _planned_

Every other rule has a `TestHas*` table-driven helper test (`TestHasConst1024`,
`TestHasStepBy3`, `TestHasEqualsOneBranch`, `TestHasCommaOrSeparator`). H008's
`hasOrdinalSwitch` (`pattern_ordinal.go`) has none — it is only integration-tested
via `TestRuleOrdinal_Positive`.

- [ ] Table-driven test mirroring `TestHasEqualsOneBranch` (6+ cases), reusing `boolSrcCase`

### T9 — White-box unit tests for `pattern_commaf.go` helpers · Medium · _planned_

The four pure helpers extracted to stay under cyclop/gocognit thresholds are
trivially testable but currently only exercised through `TestRuleCommaf_Positive`:

- [ ] `TestWalkFormatFloatVerbs` — `%%` escapes, mixed-width verbs, trailing junk
- [ ] `TestScanDottedPercentFloat` / `TestScanBarePercentFloat`
- [ ] `TestHasFormatFloatPrecision`

Risk without them: a future refactor silently regresses H009 to 0 detections.
Evidence: `pattern_commaf.go:102-167`.

---

## Documentation

### T5 — Rule doc pages: `docs/rules/H008.md` + `H009.md` · Medium · _planned_

`docs/rules/` has `H001.md`–`H007.md`. H008 (manual-ordinal) and H009
(manual-commaf) shipped in v0.2.0 but have no doc pages. The 0.2.0 CHANGELOG
advertised "Markdown docs per rule: H001.md through H007.md" — the new rules
were left out.

- [ ] `docs/rules/H008.md` matching the existing format (what it detects, example, fix, suppression)
- [ ] `docs/rules/H009.md`

### T8 — CONTRIBUTING.md rule-addition checklist · Medium · _planned_

`CONTRIBUTING.md` is 14 lines of generic fork/PR guidance. It does not mention
that adding a rule requires updating `rules.go`, `pattern_helpers.go` rule-ID
constants, `example_test.go`, `linter_test.go` count assertion, AGENTS.md, and
CHANGELOG.md — the exact drift that caused the H008/H009 doc rot. Also: the
documented dev commands (`golangci-lint run ./...`) omit the required
`GOEXPERIMENT=jsonv2` env.

- [ ] Add a "Adding a rule" checklist section
- [ ] Fix dev-setup commands to match `flake.nix` / AGENTS.md (env vars or `nix run`)

---

## Release & validation

### T6 — Tag `v0.2.0` (code shipped; tag missing) · High · _planned_

H008, H009, scoped suppression, `HumanizeDetector`, `--explain`, and `--list-files`
are all merged to `main` and documented in `CHANGELOG.md` under `[0.2.0]`. The only
git tag is `v0.1.0`. The release workflow (`.github/workflows/release.yml`) fires
on tags, so the 0.2.0 release artefacts have never been built.

- [ ] Tag `v0.2.0` on `main` after T1–T5 land
- [ ] Publish GitHub release notes

### T7 — Re-run real-world validation sweep with H008 + H009 · High · _planned_

H001–H007 were swept against 190+ Go projects (`docs/validation/2026-07-30_real-world-sweep.md`).
H008 and H009 have never been swept against the corpus, so the "~0% FP" claim does
not yet extend to them.

- [ ] Run the linter over the 190+ project corpus with all 9 rules enabled
- [ ] Record H008/H009 finding counts and false-positive rate
- [ ] Save to `docs/validation/`

---

## Tooling & CI

### T10 — Replace `nix run .#lint` `grep -v` filter with `.golangci.yml` plugin registration · High · _planned_

`flake.nix` silences the `nolint_filter` "unknown linter: gohumanize" warning with a
`grep -v`. This hides a real signal from the dev shell while a raw
`golangci-lint run` (as the CI lint job does) still emits it. The proper fix is to
register the project's own analyzer via `.golangci.yml`'s `plugins:` map.

- [ ] Register `gohumanize` as a golangci-lint plugin in `.golangci.yml`
- [ ] Remove the `grep -v` band-aid from `flake.nix`

### T12 — GitHub Action composite `action.yml` · Medium · _planned_

- [ ] Reusable composite `action.yml` (inputs: `path`, `enable`, `disable`, `format`)
- [ ] README "Use in GitHub Actions" section

---

## Larger / future

### T11 — Configurable rules in plugin mode via `Analyzer.Flags` · Medium · _planned_

The CLI supports `--enable`/`--disable`; the golangci-lint plugin does not. Users
cannot disable a rule from their `.golangci.yml` without `//nolint` (scoped
directives cover the common case, but config-level control is still missing).

- [ ] Add `enable` / `disable` string flags to `plugin.Analyzer.Flags`
- [ ] Filter `DetectFuncDecl` results by configured rule set in `analyzeHumanize`

### T13 — Package-level `var` detection for H007 · Low · _planned_

Only `*ast.FuncDecl` scope is scanned. A `var multiplier = map[string]int64{"KB": 1024}`
at package scope is invisible to H007.

- [ ] File/package-scope scan for `map[string]int64` with byte-unit keys
- [ ] Hook into `checkFuncDecls` as a separate decl-kind pass

### T14 — go/types type-aware detection · Low · _planned_

Detection is purely syntactic — import aliases (`s "strings"`) and typed values
are not resolved. Type info would cut false negatives on generic / aliased code.

- [ ] Resolve import aliases before `isPackageCall`
- [ ] Decide: full `go/types` or lightweight `pass.TypesInfo` in plugin path only
- [ ] Benchmark impact on scan speed

### T15 — `--config` flag for YAML/TOML rule configuration · Low · _planned_

- [ ] Define config schema (enabled/disabled, thresholds)
- [ ] Loader + `--config` flag + precedence over CLI flags

### T16 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on `go-linter-sdk` first tag (so `go.mod` `replace` directives can be removed)
- [ ] Submit to the plugin index once a tagged version is `go install`-able
